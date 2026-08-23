package phantom

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/minoxs/osu-phantom/pkg/osu"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// ScoresFetcher fetches one page of osu!'s global recent-scores feed newer than
// cursor, returning the page and the cursor for the next. It is the one osu-facing
// dependency of RealtimePoller, injected so the tail loop tests without the network.
type ScoresFetcher func(cursor string) (player.Scores, string, error)

// RealtimeConfig sets the tail cadence, which ruleset to follow, and the pacer
// priority the feed poll reserves at.
type RealtimeConfig struct {
	// Interval is the wait between polls once caught up. A poll returns every score
	// set since the last one, so the interval bounds latency, not completeness, as
	// long as it stays short enough that a page does not overflow.
	Interval time.Duration
	// Ruleset is the osu! ruleset name to follow, e.g. "osu". Empty means "osu".
	Ruleset string
	// Priority is the pacer level the feed poll reserves at. The caller sets it so
	// the feed never starves behind lower-priority traffic.
	Priority osu.Priority
}

const (
	// DefaultInterval keeps the tail well inside one page. osu!standard sets on the
	// order of twenty scores a second globally, and a page holds about a thousand,
	// so this leaves wide margin.
	DefaultInterval = 15 * time.Second

	// drainThreshold is the page size above which a poll likely left more behind, so
	// the next poll fires at once instead of waiting the interval.
	drainThreshold = 500
)

// RealtimePoller tails osu!'s global scores feed, the passing scores every player
// sets, fanning each out through an embedded Broadcaster until its Run context is
// done.
//
// The scores it emits are lean: osu!'s global feed omits the beatmap, sending only
// a flat beatmap_id, so Beatmap and BeatmapSet are zero on emitted scores and only
// Score.BeatmapID is set. Filling them is a consumer's concern; RealtimeTracker does
// it for tracked users.
type RealtimePoller struct {
	*Broadcaster
	fetch    ScoresFetcher
	interval time.Duration
	logger   *slog.Logger

	mu     sync.Mutex
	cursor string
}

// NewRealtimePoller builds a poller driven by the given fetcher.
func NewRealtimePoller(fetch ScoresFetcher, interval time.Duration) *RealtimePoller {
	if interval <= 0 {
		interval = DefaultInterval
	}
	return &RealtimePoller{
		Broadcaster: NewBroadcaster(),
		fetch:       fetch,
		interval:    interval,
		logger:      slog.Default().With("component", "realtime"),
	}
}

// NewOsuRealtimePoller builds a RealtimePoller that tails osu!'s feed through
// osu-phantom.
func NewOsuRealtimePoller(provider AuthProvider, cfg RealtimeConfig) *RealtimePoller {
	ruleset := cfg.Ruleset
	if ruleset == "" {
		ruleset = "osu"
	}
	client := osu.NewClient(cfg.Priority)
	fetch := func(cursor string) (player.Scores, string, error) {
		return client.GetScores(provider.GetToken(), ruleset, cursor)
	}
	return NewRealtimePoller(fetch, cfg.Interval)
}

// Run streams the feed forward, emitting every score set after it started, until
// ctx is cancelled. It blocks, so callers run it in a goroutine. With no resume
// cursor it first seeds at the newest score, whose page is a historical snapshot
// and so is used only to seed and is not emitted; resumed from a persisted cursor
// it instead streams the scores set since that cursor, so a restart does not skip
// scores set during downtime.
func (p *RealtimePoller) Run(ctx context.Context) {
	cursor := p.Cursor()
	if cursor == "" {
		seeded, ok := p.seed(ctx)
		if !ok {
			return
		}
		cursor = seeded
		p.setCursor(cursor)
	}

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			scores, next, err := p.fetch(cursor)
			if err != nil {
				p.logger.Warn("realtime poll failed", "err", err)
				timer.Reset(p.interval)
				continue
			}
			for i := range scores {
				p.Emit(scores[i])
			}
			if next != "" {
				cursor = next
				p.setCursor(next)
			}
			if len(scores) >= drainThreshold {
				timer.Reset(0)
			} else {
				timer.Reset(p.interval)
			}
		}
	}
}

// Cursor returns the feed position the poller has reached, empty before it seeds.
// Persist it and pass it back through Resume to continue the tail across a restart,
// so scores set during downtime are streamed rather than skipped.
func (p *RealtimePoller) Cursor() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cursor
}

// Resume sets the position the next Run starts from, so the tail streams the scores
// set since it instead of seeding at the newest and skipping the gap. Call before
// Run; an empty cursor leaves Run to seed as usual.
func (p *RealtimePoller) Resume(cursor string) { p.setCursor(cursor) }

func (p *RealtimePoller) setCursor(cursor string) {
	p.mu.Lock()
	p.cursor = cursor
	p.mu.Unlock()
}

// seed fetches the newest page to adopt its cursor without emitting it, retrying on
// error. It returns false only when ctx is cancelled.
func (p *RealtimePoller) seed(ctx context.Context) (string, bool) {
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", false
		case <-timer.C:
			_, cursor, err := p.fetch("")
			if err == nil && cursor != "" {
				return cursor, true
			}
			if err != nil {
				p.logger.Warn("realtime seed failed", "err", err)
			}
			timer.Reset(p.interval)
		}
	}
}
