package phantom

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// fakeBeatmaps scripts the beatmap fetch and counts it, so a test can assert that a
// gate dropped a score before any fetch.
type fakeBeatmaps struct {
	beatmap func(id int64) (player.Beatmap, player.BeatmapSet, error)
	fetches int
}

func (f *fakeBeatmaps) Beatmap(id int64) (player.Beatmap, player.BeatmapSet, error) {
	f.fetches++
	return f.beatmap(id)
}

// rankedBeatmap fills a ranked map so enrichment passes the status gate.
func rankedBeatmap(id int64) (player.Beatmap, player.BeatmapSet, error) {
	return player.Beatmap{ID: id, Status: player.StatusRanked}, player.BeatmapSet{Title: "T"}, nil
}

func newTracked(ctx context.Context, b *Broadcaster, beatmaps BeatmapProvider, user int) *RealtimeTracker {
	tr := NewRealtimeTracker(ctx, b, beatmaps)
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
	tr := newTracked(ctx, b, &fakeBeatmaps{beatmap: rankedBeatmap}, 7)

	b.Emit(rankedScore(7))

	got := recv(t, tr.Scores())
	if got.Beatmap.Status != player.StatusRanked || got.BeatmapSet.Title != "T" {
		t.Fatalf("score not enriched: %+v", got.Beatmap)
	}
	if got.PP != 100 {
		t.Fatalf("feed pp not trusted: got %v, want 100", got.PP)
	}
}

func TestRealtimeTracker_DropsUntrackedUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	beatmaps := &fakeBeatmaps{beatmap: rankedBeatmap}
	tr := newTracked(ctx, b, beatmaps, 7)

	b.Emit(rankedScore(9)) // different user

	assertNoScore(t, tr.Scores())
	if beatmaps.fetches != 0 {
		t.Fatalf("untracked user triggered %d fetches, want 0", beatmaps.fetches)
	}
}

func TestRealtimeTracker_DropsNonStandardBeforeFetch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	beatmaps := &fakeBeatmaps{beatmap: rankedBeatmap}
	tr := newTracked(ctx, b, beatmaps, 7)

	s := rankedScore(7)
	s.RulesetID = 3 // mania
	b.Emit(s)

	assertNoScore(t, tr.Scores())
	if beatmaps.fetches != 0 {
		t.Fatalf("non-standard score triggered %d fetches, want 0", beatmaps.fetches)
	}
}

func TestRealtimeTracker_DropsUnrankedFlagBeforeFetch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	beatmaps := &fakeBeatmaps{beatmap: rankedBeatmap}
	tr := newTracked(ctx, b, beatmaps, 7)

	s := rankedScore(7)
	s.Ranked = false
	b.Emit(s)

	assertNoScore(t, tr.Scores())
	if beatmaps.fetches != 0 {
		t.Fatalf("unranked-flag score triggered %d fetches, want 0", beatmaps.fetches)
	}
}

func TestRealtimeTracker_DropsWhenStatusAwardsNoPP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	loved := &fakeBeatmaps{beatmap: func(id int64) (player.Beatmap, player.BeatmapSet, error) {
		return player.Beatmap{ID: id, Status: player.StatusLoved}, player.BeatmapSet{}, nil
	}}
	tr := newTracked(ctx, b, loved, 7)

	b.Emit(rankedScore(7))

	assertNoScore(t, tr.Scores())
}

func TestRealtimeTracker_DropsOnFetchError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	failing := &fakeBeatmaps{beatmap: func(int64) (player.Beatmap, player.BeatmapSet, error) {
		return player.Beatmap{}, player.BeatmapSet{}, errors.New("boom")
	}}
	tr := newTracked(ctx, b, failing, 7)

	b.Emit(rankedScore(7))

	assertNoScore(t, tr.Scores())
}

func TestRealtimeTracker_DropsScoreBeforeSince(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := NewBroadcaster()
	tr := NewRealtimeTracker(ctx, b, &fakeBeatmaps{beatmap: rankedBeatmap})
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
