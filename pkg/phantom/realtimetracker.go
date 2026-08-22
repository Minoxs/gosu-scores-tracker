package phantom

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// BeatmapFetcher fills the beatmap a lean feed score omits. osu.GetBeatmap
// satisfies it once bound to a token.
type BeatmapFetcher func(id int64) (player.Beatmap, player.BeatmapSet, error)

// PPResolver sets a score's pp when the feed did not report it. osu.GetPP
// satisfies it. Injected so this package stays free of the cgo pp path.
type PPResolver func(*player.Score)

// RealtimeTracker is a ScoreTracker over a lean score provider such as
// RealtimePoller. It follows a chosen set of users and, for each of their new
// scores, applies cheap gates before spending a beatmap fetch: tracked user and
// since, osu!standard ruleset, then the feed's ranked flag. Only a score past all
// three is enriched with its beatmap, checked against the authoritative map status,
// given pp, and emitted.
type RealtimeTracker struct {
	mu      sync.Mutex
	tracked map[int]time.Time
	out     chan player.Score

	fetch     BeatmapFetcher
	resolvePP PPResolver
	logger    *slog.Logger
}

var _ ScoreTracker = (*RealtimeTracker)(nil)

// NewRealtimeTracker subscribes to provider and starts forwarding tracked users'
// ranked osu!standard scores, each enriched and with pp resolved. It stops and
// closes Scores when ctx is cancelled or the provider closes the subscription.
func NewRealtimeTracker(ctx context.Context, provider ScoreProvider, fetch BeatmapFetcher, resolvePP PPResolver) *RealtimeTracker {
	t := &RealtimeTracker{
		tracked:   make(map[int]time.Time),
		out:       make(chan player.Score, DefaultSubBuffer),
		fetch:     fetch,
		resolvePP: resolvePP,
		logger:    slog.Default().With("component", "realtime"),
	}
	go t.run(ctx, provider.Subscribe())
	return t
}

func (t *RealtimeTracker) run(ctx context.Context, in <-chan player.Score) {
	defer close(t.out)
	for {
		select {
		case <-ctx.Done():
			return
		case s, ok := <-in:
			if !ok {
				return
			}
			if !t.admit(s) {
				continue
			}
			if !t.enrich(&s) {
				continue
			}
			select {
			case t.out <- s:
			case <-ctx.Done():
				return
			}
		}
	}
}

// admit runs the gates that need no fetch: tracked user set after their start, an
// osu!standard ruleset, and the feed's ranked flag as a cheap pre-check.
func (t *RealtimeTracker) admit(s player.Score) bool {
	t.mu.Lock()
	since, tracked := t.tracked[s.UserID]
	t.mu.Unlock()

	if !tracked || !s.EndedAt.After(since) {
		return false
	}
	if s.RulesetID != 0 {
		return false
	}
	return s.Ranked
}

// enrich fills the score's beatmap and resolves pp, returning false when the fetch
// fails or the map's authoritative status awards no pp. It reports whether the
// score should be emitted.
func (t *RealtimeTracker) enrich(s *player.Score) bool {
	bm, bs, err := t.fetch(s.BeatmapID)
	if err != nil {
		t.logger.Warn("realtime beatmap fetch failed", "beatmap", s.BeatmapID, "err", err)
		return false
	}
	s.Beatmap = bm
	s.BeatmapSet = bs
	if !s.IsRanked() {
		return false
	}
	t.resolvePP(s)
	return true
}

// Track starts forwarding userID's scores set after since. Calling it again with a
// new start moves the boundary.
func (t *RealtimeTracker) Track(userID int, since time.Time) {
	t.mu.Lock()
	t.tracked[userID] = since
	t.mu.Unlock()
}

// Untrack stops forwarding userID's scores.
func (t *RealtimeTracker) Untrack(userID int) {
	t.mu.Lock()
	delete(t.tracked, userID)
	t.mu.Unlock()
}

// Scores is the stream of tracked users' enriched, pp-resolved ranked scores.
func (t *RealtimeTracker) Scores() <-chan player.Score { return t.out }
