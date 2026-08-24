package phantom

import (
	"context"
	"testing"
	"time"
)

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
	b.Emit(score(22, 1, now))                     // tracked user, in window

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
