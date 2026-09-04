package tracker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

// scoresFetcher fetches one page of the global feed newer than cursor, returning
// the page and the cursor for the next
type scoresFetcher func(cursor string) (gosu.Scores, string, error)

// RealtimeConfig sets the tail cadence, ruleset, and pacer priority
type RealtimeConfig struct {
	// Interval is the wait between polls once caught up
	Interval time.Duration
	// Ruleset is the osu! ruleset to follow, e.g. "osu". Empty means "osu"
	Ruleset string
	// Priority is the pacer level the feed poll reserves at
	Priority gosu.Priority
}

const (
	// DefaultInterval keeps the tail well inside one page. The feed returns up to
	// 1000 scores a page and osu!standard sets on the order of twenty a second
	DefaultInterval = 15 * time.Second

	// drainThreshold is the page size above which a poll likely left more behind,
	// so the next poll fires at once
	drainThreshold = 500
)

// RealtimePoller tails osu!'s global scores feed, fanning each score to its
// subscribers until Run stops. The scores are lean: the feed omits the beatmap and
// sends only a beatmap_id
type RealtimePoller struct {
	bc       *broadcaster[gosu.Score]
	fetch    scoresFetcher
	interval time.Duration
	logger   *slog.Logger

	mu     sync.Mutex
	cursor string
}

// NewRealtimePoller builds a poller that tails the feed through app, sharing its
// rate ceiling so your own app.GuestClient calls are never crowded out
func NewRealtimePoller(app *gosu.App, cfg RealtimeConfig) *RealtimePoller {
	ruleset := cfg.Ruleset
	if ruleset == "" {
		ruleset = "osu"
	}
	client := app.GuestClient(cfg.Priority)
	fetch := func(cursor string) (gosu.Scores, string, error) {
		return client.GetScores(ruleset, cursor)
	}
	return newRealtimePoller(fetch, cfg.Interval)
}

func newRealtimePoller(fetch scoresFetcher, interval time.Duration) *RealtimePoller {
	if interval <= 0 {
		interval = DefaultInterval
	}
	return &RealtimePoller{
		bc:       newBroadcaster[gosu.Score](),
		fetch:    fetch,
		interval: interval,
		logger:   slog.Default().With("component", "realtime"),
	}
}

// Subscribe returns an independent stream of the feed's scores. Drain it or delivery
// blocks and slows the feed. Close it to unsubscribe; Run closes the rest when it stops
func (p *RealtimePoller) Subscribe() chan gosu.Score { return p.bc.subscribe() }

// Run streams the feed forward until ctx is cancelled. It blocks, so run it in a
// goroutine. Seeded fresh it skips the newest page as a historical snapshot and
// streams what follows; resumed from a cursor it streams the scores set since, so a
// restart does not skip scores set during downtime
func (p *RealtimePoller) Run(ctx context.Context) {
	defer p.bc.close()

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
				p.bc.emit(ctx, scores[i])
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
// Persist it and pass it to Resume to continue the tail across a restart
func (p *RealtimePoller) Cursor() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cursor
}

// Resume sets the position the next Run starts from. Call before Run; an empty
// cursor leaves Run to seed as usual
func (p *RealtimePoller) Resume(cursor string) { p.setCursor(cursor) }

func (p *RealtimePoller) setCursor(cursor string) {
	p.mu.Lock()
	p.cursor = cursor
	p.mu.Unlock()
}

// seed fetches the newest page to adopt its cursor without emitting it, retrying on
// error. It returns false only when ctx is cancelled
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
