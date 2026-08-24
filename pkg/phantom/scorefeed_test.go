package phantom

import (
	"context"
	"testing"
	"time"

	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

func score(id, user int, at time.Time) player.Score {
	return player.Score{ID: int64(id), UserID: user, EndedAt: at}
}

func recv(t *testing.T, ch <-chan player.Score) player.Score {
	t.Helper()
	select {
	case s := <-ch:
		return s
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for score")
		return player.Score{}
	}
}

// assertNoScore fails if any score arrives within a short grace period.
func assertNoScore(t *testing.T, ch <-chan player.Score) {
	t.Helper()
	select {
	case s := <-ch:
		t.Fatalf("unexpected score %d", s.ID)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestBroadcaster_FanOut(t *testing.T) {
	b := NewBroadcaster()
	a := b.Subscribe()
	c := b.Subscribe()

	b.Emit(score(1, 7, time.Now()))

	if got := recv(t, a); got.ID != 1 {
		t.Fatalf("subscriber a got ID %d, want 1", got.ID)
	}
	if got := recv(t, c); got.ID != 1 {
		t.Fatalf("subscriber c got ID %d, want 1", got.ID)
	}
}

func TestBroadcaster_EmitDoesNotBlockOnFullSubscriber(t *testing.T) {
	b := NewBroadcaster()
	_ = b.Subscribe() // never drained

	done := make(chan struct{})
	go func() {
		for i := 0; i < DefaultSubBuffer+10; i++ {
			b.Emit(score(i, 1, time.Now()))
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Emit blocked on a subscriber that never drains")
	}
}

func TestBroadcaster_CloseEndsSubscriptions(t *testing.T) {
	b := NewBroadcaster()
	ch := b.Subscribe()
	b.Close()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected subscription channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for close")
	}
}

func TestFilterTracker_ForwardsTracked(t *testing.T) {
	b := NewBroadcaster()
	tr := NewFilterTracker(context.Background(), b)
	tr.Track(1, time.Now().Add(-time.Hour))

	b.Emit(score(11, 1, time.Now()))

	if got := recv(t, tr.Scores()); got.ID != 11 || got.UserID != 1 {
		t.Fatalf("got score %d for user %d, want 11 for user 1", got.ID, got.UserID)
	}
}

func TestFilterTracker_SkipsUntrackedAndPreWindow(t *testing.T) {
	b := NewBroadcaster()
	tr := NewFilterTracker(context.Background(), b)
	since := time.Now().Add(-time.Hour)
	tr.Track(1, since)

	now := time.Now()
	b.Emit(score(20, 2, now))                     // untracked user
	b.Emit(score(21, 1, since.Add(-time.Minute))) // tracked user, before window
	b.Emit(score(22, 1, now))                      // tracked user, in window

	if got := recv(t, tr.Scores()); got.ID != 22 {
		t.Fatalf("first forwarded score was ID %d, want 22", got.ID)
	}
}

func TestFilterTracker_Untrack(t *testing.T) {
	b := NewBroadcaster()
	tr := NewFilterTracker(context.Background(), b)
	since := time.Now().Add(-time.Hour)
	tr.Track(1, since)
	tr.Untrack(1)
	tr.Track(2, since)

	now := time.Now()
	b.Emit(score(30, 1, now)) // untracked again, dropped
	b.Emit(score(31, 2, now)) // tracked, forwarded

	if got := recv(t, tr.Scores()); got.ID != 31 {
		t.Fatalf("first forwarded score was ID %d, want 31", got.ID)
	}
}

func TestFilterTracker_ClosesOnContextCancel(t *testing.T) {
	b := NewBroadcaster()
	ctx, cancel := context.WithCancel(context.Background())
	tr := NewFilterTracker(ctx, b)

	cancel()

	select {
	case _, ok := <-tr.Scores():
		if ok {
			t.Fatal("expected Scores channel to close after cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for Scores to close")
	}
}
