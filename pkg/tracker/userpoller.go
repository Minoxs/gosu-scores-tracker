package tracker

import (
	"context"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

// FullScoreProvider is a source of scores with their maps embedded. UserPoller is
// one: unlike the lean realtime ScoreProvider, the user-scores endpoint returns
// FullScore, so its stream carries the beatmap and beatmapset.
type FullScoreProvider interface {
	Scores() <-chan gosu.FullScore
}

// FullScoreFetcher returns a user's scores set after since, as the API returns them
// with the map embedded, plus the newest score's time seen so the caller can advance
// its watermark. It is the one osu!-facing dependency of UserPoller, injected so the
// polling loop can be tested without the network.
type FullScoreFetcher func(userID int64, since time.Time) (scores gosu.FullScores, newest time.Time, err error)

// PollConfig sets the per-user polling cadence. The osu! terms of use forbid
// polling a user more than once a minute, so BaseInterval floors the cadence and
// MaxInterval caps the backoff applied while a user sets nothing new.
type PollConfig struct {
	BaseInterval time.Duration
	MaxInterval  time.Duration
	JitterFrac   float64
}

// UserPoller is a ScoreTracker that surfaces tracked users' new ranked scores by
// polling the osu! API. Each tracked user runs its own poll loop, backing off
// while idle. Accumulation is the consumer's job.
type UserPoller struct {
	fetch FullScoreFetcher
	cfg   PollConfig
	out   chan gosu.FullScore

	// OnCheck, when set before the first Track, is called after every poll with
	// the time of the poll and when the next is due, so a consumer can report how
	// fresh a user's data is even across polls that found nothing new. It must not
	// be mutated once Track has been called.
	OnCheck func(userID int64, checked, next time.Time)

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu      sync.Mutex
	tracked map[int64]context.CancelFunc
	closed  bool
}

// NewUserPoller builds a tracker driven by the given fetcher.
func NewUserPoller(fetch FullScoreFetcher, cfg PollConfig) *UserPoller {
	if cfg.BaseInterval <= 0 {
		cfg.BaseInterval = time.Minute
	}
	if cfg.MaxInterval < cfg.BaseInterval {
		cfg.MaxInterval = 30 * time.Minute
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &UserPoller{
		fetch:   fetch,
		cfg:     cfg,
		out:     make(chan gosu.FullScore, DefaultSubBuffer),
		ctx:     ctx,
		cancel:  cancel,
		tracked: make(map[int64]context.CancelFunc),
	}
}

// NewOsuUserPoller builds a UserPoller that fetches through gosu-api.
func NewOsuUserPoller(provider AuthProvider, cfg PollConfig) *UserPoller {
	return NewUserPoller(osuScoreFetcher(provider), cfg)
}

// Track starts polling userID, surfacing scores set after since. It is idempotent
// and a no-op once the tracker is closed.
func (t *UserPoller) Track(userID int64, since time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return
	}
	if _, ok := t.tracked[userID]; ok {
		return
	}
	ctx, cancel := context.WithCancel(t.ctx)
	t.tracked[userID] = cancel
	t.wg.Add(1)
	go t.loop(ctx, userID, since)
}

// Untrack stops polling userID.
func (t *UserPoller) Untrack(userID int64) {
	t.mu.Lock()
	cancel, ok := t.tracked[userID]
	delete(t.tracked, userID)
	t.mu.Unlock()
	if ok {
		cancel()
	}
}

// Scores is the stream of tracked users' new scores, each with its map embedded.
func (t *UserPoller) Scores() <-chan gosu.FullScore { return t.out }

// Close stops every poll loop and closes the Scores channel once they have all
// exited. The tracker cannot be reused afterwards.
func (t *UserPoller) Close() {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return
	}
	t.closed = true
	t.mu.Unlock()

	t.cancel()
	t.wg.Wait()
	close(t.out)
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
				slog.Error("poll fetch failed", "user", userID, "err", err)
			}

			fresh := false
			for _, s := range scores {
				if !s.EndedAt.After(watermark) {
					continue
				}
				if !t.emit(ctx, s) {
					return
				}
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

			delay := jitter(interval, t.cfg.JitterFrac)
			if t.OnCheck != nil {
				now := time.Now()
				t.OnCheck(userID, now, now.Add(delay))
			}
			timer.Reset(delay)
		}
	}
}

// emit forwards a score, blocking until the consumer takes it so no score is
// dropped. It bails out if the loop is cancelled while waiting.
func (t *UserPoller) emit(ctx context.Context, s gosu.FullScore) bool {
	select {
	case t.out <- s:
		return true
	case <-ctx.Done():
		return false
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

// jitter adds up to frac of d, never subtracting, so the per-user minimum holds.
func jitter(d time.Duration, frac float64) time.Duration {
	if frac <= 0 {
		return d
	}
	return d + time.Duration(rand.Float64()*frac*float64(d))
}

// osuScoreFetcher fetches a user's recent scores through gosu-api, paging while
// a whole page is newer than since. It forwards the crude scores the API returns.
func osuScoreFetcher(provider AuthProvider) FullScoreFetcher {
	client := gosu.NewClient(0)
	return func(userID int64, since time.Time) (gosu.FullScores, time.Time, error) {
		var out gosu.FullScores
		var newest time.Time

		for offset := 0; ; offset += recentScorePageSize {
			page := client.GetRecentScores(provider.GetToken(), userID, recentScorePageSize, offset)
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
