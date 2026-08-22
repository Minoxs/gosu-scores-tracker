package phantom

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// rankedFetch fills a ranked beatmap so enrichment passes the status gate.
func rankedFetch(id int64) (player.Beatmap, player.BeatmapSet, error) {
	return player.Beatmap{ID: id, Status: player.StatusRanked}, player.BeatmapSet{Title: "T"}, nil
}

// countingResolve marks that pp resolution ran and only fills when unset.
func countingResolve(hits *int) PPResolver {
	return func(s *player.Score) {
		*hits++
		if s.PP == 0 {
			s.PP = 1
		}
	}
}

func newTracked(ctx context.Context, b *Broadcaster, fetch BeatmapFetcher, resolve PPResolver, user int) *RealtimeTracker {
	tr := NewRealtimeTracker(ctx, b, fetch, resolve)
	tr.Track(user, time.Unix(0, 0))
	return tr
}

func rankedScore(user int) player.Score {
	return player.Score{
		ID: 1, UserID: user, BeatmapID: 42, RulesetID: 0, Ranked: true,
		PP: 100, EndedAt: time.Unix(100, 0),
	}
}

func TestRealtimeTracker_EmitsEnrichedRankedScore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	hits := 0
	tr := newTracked(ctx, b, rankedFetch, countingResolve(&hits), 7)

	b.Emit(rankedScore(7))

	got := recv(t, tr.Scores())
	if got.Beatmap.Status != player.StatusRanked || got.BeatmapSet.Title != "T" {
		t.Fatalf("score not enriched: %+v", got.Beatmap)
	}
	if hits != 1 {
		t.Fatalf("pp resolver ran %d times, want 1", hits)
	}
}

func TestRealtimeTracker_DropsUntrackedUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	fetched := 0
	tr := newTracked(ctx, b, func(id int64) (player.Beatmap, player.BeatmapSet, error) {
		fetched++
		return rankedFetch(id)
	}, countingResolve(new(int)), 7)

	b.Emit(rankedScore(9)) // different user

	assertNoScore(t, tr.Scores())
	if fetched != 0 {
		t.Fatalf("untracked user triggered %d fetches, want 0", fetched)
	}
}

func TestRealtimeTracker_DropsNonStandardBeforeFetch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	fetched := 0
	tr := newTracked(ctx, b, func(id int64) (player.Beatmap, player.BeatmapSet, error) {
		fetched++
		return rankedFetch(id)
	}, countingResolve(new(int)), 7)

	s := rankedScore(7)
	s.RulesetID = 3 // mania
	b.Emit(s)

	assertNoScore(t, tr.Scores())
	if fetched != 0 {
		t.Fatalf("non-standard score triggered %d fetches, want 0", fetched)
	}
}

func TestRealtimeTracker_DropsUnrankedFlagBeforeFetch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	fetched := 0
	tr := newTracked(ctx, b, func(id int64) (player.Beatmap, player.BeatmapSet, error) {
		fetched++
		return rankedFetch(id)
	}, countingResolve(new(int)), 7)

	s := rankedScore(7)
	s.Ranked = false
	b.Emit(s)

	assertNoScore(t, tr.Scores())
	if fetched != 0 {
		t.Fatalf("unranked-flag score triggered %d fetches, want 0", fetched)
	}
}

func TestRealtimeTracker_DropsWhenStatusAwardsNoPP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	resolved := 0
	loved := func(id int64) (player.Beatmap, player.BeatmapSet, error) {
		return player.Beatmap{ID: id, Status: player.StatusLoved}, player.BeatmapSet{}, nil
	}
	tr := newTracked(ctx, b, loved, countingResolve(&resolved), 7)

	b.Emit(rankedScore(7))

	assertNoScore(t, tr.Scores())
	if resolved != 0 {
		t.Fatalf("pp resolved for a no-pp map %d times, want 0", resolved)
	}
}

func TestRealtimeTracker_DropsOnFetchError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	failing := func(int64) (player.Beatmap, player.BeatmapSet, error) {
		return player.Beatmap{}, player.BeatmapSet{}, errors.New("boom")
	}
	tr := newTracked(ctx, b, failing, countingResolve(new(int)), 7)

	b.Emit(rankedScore(7))

	assertNoScore(t, tr.Scores())
}

func TestRealtimeTracker_DropsScoreBeforeSince(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	tr := NewRealtimeTracker(ctx, b, rankedFetch, countingResolve(new(int)))
	tr.Track(7, time.Unix(500, 0))

	b.Emit(rankedScore(7)) // ended at t=100, before since=500

	assertNoScore(t, tr.Scores())
}

func assertNoScore(t *testing.T, ch <-chan player.Score) {
	t.Helper()
	select {
	case s := <-ch:
		t.Fatalf("unexpected score emitted: id %d", s.ID)
	case <-time.After(150 * time.Millisecond):
	}
}
