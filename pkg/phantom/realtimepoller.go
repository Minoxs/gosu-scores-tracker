package phantom

import (
	"context"
	"log/slog"
	"time"

	"github.com/minoxs/osu-phantom/pkg/osu"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// ScoresFetcher fetches one page of osu!'s global recent-scores feed newer than
// cursor, returning the page and the cursor for the next. It is the one osu-facing
// dependency of RealtimePoller, injected so the tail loop tests without the network.
type ScoresFetcher func(cursor string) (player.Scores, string, error)

// RealtimeConfig sets the tail cadence and which ruleset to follow.
type RealtimeConfig struct {
	// Interval is the wait between polls once caught up. A poll returns every score
	// set since the last one, so the interval bounds latency, not completeness, as
	// long as it stays short enough that a page does not overflow.
	Interval time.Duration
	// Ruleset is the osu! ruleset name to follow, e.g. "osu". Empty means "osu".
	Ruleset string
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
	fetch := func(cursor string) (player.Scores, string, error) {
		return osu.GetScores(provider.GetToken(), ruleset, cursor)
	}
	return NewRealtimePoller(fetch, cfg.Interval)
}

// Run seeds the cursor at the newest score, then streams forward, emitting every
// score set after it started, until ctx is cancelled. It blocks, so callers run it
// in a goroutine. The newest page at start is a historical snapshot, not new
// events, so it is used only to seed the cursor and is not emitted.
func (p *RealtimePoller) Run(ctx context.Context) {
	cursor, ok := p.seed(ctx)
	if !ok {
		return
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
			}
			if len(scores) >= drainThreshold {
				timer.Reset(0)
			} else {
				timer.Reset(p.interval)
			}
		}
	}
}

// seed fetches the newest page to adopt its cursor without emitting it, retrying
// on error so streaming never falls back to replaying the whole newest page. It
// returns false only when ctx is cancelled.
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
