package phantom

import (
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
