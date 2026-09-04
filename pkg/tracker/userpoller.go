package tracker

import (
	"context"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

// fullScoreFetcher returns a user's scores set after since, with the map embedded,
// plus the newest score's time so the caller can advance its watermark. It is the
// one osu!-facing dependency of UserPoller, injected so the loop tests without the network
type fullScoreFetcher func(userID int64, since time.Time) (scores gosu.FullScores, newest time.Time, err error)

// PollConfig sets the per-user polling cadence
type PollConfig struct {
	// BaseInterval floors the cadence, at least a minute per the osu! terms
	BaseInterval time.Duration
	// MaxInterval caps the backoff applied while a user sets nothing new
	MaxInterval time.Duration
	// JitterFrac is the fraction of BaseInterval added at random, never subtracted
	JitterFrac float64
}

// UserPoller polls tracked users for new scores, fanning each to its subscribers.
// Each tracked user runs its own poll loop, backing off while idle. Its scores are
// full, with the beatmap and beatmapset embedded
type UserPoller struct {
	fetch  fullScoreFetcher
	cfg    PollConfig
	bc     *broadcaster[gosu.FullScore]
	logger *slog.Logger

	mu      sync.Mutex
	ctx     context.Context
	running bool
	tracked map[int64]time.Time
	cancels map[int64]context.CancelFunc
	wg      sync.WaitGroup
}

// NewUserPoller builds a poller that fetches through app, sharing its rate ceiling
// so your own app.GuestClient calls are never crowded out
func NewUserPoller(app *gosu.App, cfg PollConfig) *UserPoller {
	return newUserPoller(osuScoreFetcher(app), cfg)
}

func newUserPoller(fetch fullScoreFetcher, cfg PollConfig) *UserPoller {
	if cfg.BaseInterval <= 0 {
		cfg.BaseInterval = time.Minute
	}
	if cfg.MaxInterval < cfg.BaseInterval {
		cfg.MaxInterval = 30 * time.Minute
	}
	return &UserPoller{
		fetch:   fetch,
		cfg:     cfg,
		bc:      newBroadcaster[gosu.FullScore](),
		logger:  slog.Default().With("component", "userpoller"),
		tracked: make(map[int64]time.Time),
		cancels: make(map[int64]context.CancelFunc),
	}
}

// Subscribe returns an independent stream of tracked users' scores. Drain it or
// delivery blocks and slows polling. Close it to unsubscribe; Run closes the rest
func (t *UserPoller) Subscribe() chan gosu.FullScore { return t.bc.subscribe() }

// Track starts polling userID for scores set after since. Before Run it is
// remembered and started when Run begins; while running it starts at once. Idempotent
func (t *UserPoller) Track(userID int64, since time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.tracked[userID]; ok {
		return
	}
	t.tracked[userID] = since
	if t.running {
		t.startLoop(userID, since)
	}
}

// Untrack stops polling userID
func (t *UserPoller) Untrack(userID int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tracked, userID)
	if cancel, ok := t.cancels[userID]; ok {
		cancel()
		delete(t.cancels, userID)
	}
}

// Run polls every tracked user until ctx is cancelled. It blocks, so run it in a
// goroutine
func (t *UserPoller) Run(ctx context.Context) {
	t.mu.Lock()
	if t.running {
		t.mu.Unlock()
		return
	}
	t.running = true
	t.ctx = ctx
	for userID, since := range t.tracked {
		t.startLoop(userID, since)
	}
	t.mu.Unlock()

	<-ctx.Done()

	t.mu.Lock()
	t.running = false
	for userID, cancel := range t.cancels {
		cancel()
		delete(t.cancels, userID)
	}
	t.mu.Unlock()

	t.wg.Wait()
	t.bc.close()
}

// startLoop launches userID's poll loop. The caller holds t.mu
func (t *UserPoller) startLoop(userID int64, since time.Time) {
	ctx, cancel := context.WithCancel(t.ctx)
	t.cancels[userID] = cancel
	t.wg.Add(1)
	go t.loop(ctx, userID, since)
}

func (t *UserPoller) loop(ctx context.Context, userID int64, since time.Time) {
	defer t.wg.Done()

	watermark := since
	interval := t.cfg.BaseInterval
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			scores, newest, err := t.fetch(userID, watermark)
			if err != nil {
				t.logger.Warn("poll fetch failed", "user", userID, "err", err)
			}

			fresh := false
			for _, s := range scores {
				if !s.EndedAt.After(watermark) {
					continue
				}
				t.bc.emit(ctx, s)
				fresh = true
			}
			if newest.After(watermark) {
				watermark = newest
			}

			if fresh {
				interval = t.cfg.BaseInterval
			} else {
				interval = clampInterval(interval*2, t.cfg.BaseInterval, t.cfg.MaxInterval)
			}
			timer.Reset(jitter(interval, t.cfg.JitterFrac))
		}
	}
}

func clampInterval(d, min, max time.Duration) time.Duration {
	if d < min {
		return min
	}
	if d > max {
		return max
	}
	return d
}

// jitter adds up to frac of d, never subtracting, so the per-user minimum holds
func jitter(d time.Duration, frac float64) time.Duration {
	if frac <= 0 {
		return d
	}
	return d + time.Duration(rand.Float64()*frac*float64(d))
}

// recentScorePageSize is the per-request score count. The osu! API caps it at 100
const recentScorePageSize = 100

// osuScoreFetcher fetches a user's recent scores through gosu-api, paging while a
// whole page is newer than since
func osuScoreFetcher(app *gosu.App) fullScoreFetcher {
	client := app.GuestClient(0)
	return func(userID int64, since time.Time) (gosu.FullScores, time.Time, error) {
		var out gosu.FullScores
		var newest time.Time

		for offset := 0; ; offset += recentScorePageSize {
			page, err := client.GetRecentScores(userID, recentScorePageSize, offset)
			if err != nil {
				return out, newest, err
			}
			if len(page) == 0 {
				break
			}
			if offset == 0 {
				newest = page[0].EndedAt
			}

			allNew := true
			for i := range page {
				s := page[i]
				if !s.EndedAt.After(since) {
					allNew = false
					break
				}
				out = append(out, s)
			}

			if !allNew || len(page) < recentScorePageSize {
				break
			}
		}
		return out, newest, nil
	}
}
