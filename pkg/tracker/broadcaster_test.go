package tracker

import (
	"context"
	"testing"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

func score(id int, user int64, at time.Time) gosu.Score {
	return gosu.Score{ID: int64(id), UserID: user, EndedAt: at}
}

func recv(t *testing.T, ch <-chan gosu.Score) gosu.Score {
	t.Helper()
	select {
	case s := <-ch:
		return s
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for score")
		return gosu.Score{}
	}
}

// assertNoScore fails if any score arrives within a short grace period.
func assertNoScore(t *testing.T, ch <-chan gosu.Score) {
	t.Helper()
	select {
	case s := <-ch:
		t.Fatalf("unexpected score %d", s.ID)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestBroadcaster_FanOut(t *testing.T) {
	b := newBroadcaster[gosu.Score]()
	a := b.subscribe()
	c := b.subscribe()

	b.emit(context.Background(), score(1, 7, time.Now()))

	if got := recv(t, a); got.ID != 1 {
		t.Fatalf("subscriber a got ID %d, want 1", got.ID)
	}
	if got := recv(t, c); got.ID != 1 {
		t.Fatalf("subscriber c got ID %d, want 1", got.ID)
	}
}

// TestBroadcaster_EmitBlocksUntilConsumed proves emit does not drop past the buffer:
// it blocks once the buffer fills and completes only after the consumer drains.
func TestBroadcaster_EmitBlocksUntilConsumed(t *testing.T) {
	b := newBroadcaster[gosu.Score]()
	ch := b.subscribe()

	done := make(chan struct{})
	go func() {
		for i := 0; i < subBuffer+2; i++ {
			b.emit(context.Background(), score(i, 1, time.Now()))
		}
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("emit did not block once the buffer filled")
	case <-time.After(50 * time.Millisecond):
	}

	go func() {
		for i := 0; i < subBuffer+2; i++ {
			<-ch
		}
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("emit did not finish after the consumer drained")
	}
}

// TestBroadcaster_EmitReleasesOnContextCancel proves a blocked emit lets go when its
// context is cancelled, so a stalled consumer never wedges shutdown.
func TestBroadcaster_EmitReleasesOnContextCancel(t *testing.T) {
	b := newBroadcaster[gosu.Score]()
	_ = b.subscribe() // never drained

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		for i := 0; i < subBuffer+5; i++ {
			b.emit(ctx, score(i, 1, time.Now()))
		}
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("cancel did not release a blocked emit")
	}
}

func TestBroadcaster_CloseEndsSubscriptions(t *testing.T) {
	b := newBroadcaster[gosu.Score]()
	ch := b.subscribe()
	b.close()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected subscription channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for close")
	}
}

// TestBroadcaster_CloseToUnsubscribe proves a consumer closing its channel drops it
// without panicking the emit and without disturbing the remaining subscribers.
func TestBroadcaster_CloseToUnsubscribe(t *testing.T) {
	b := newBroadcaster[gosu.Score]()
	stay := b.subscribe()
	leave := b.subscribe()

	close(leave)

	b.emit(context.Background(), score(1, 7, time.Now()))

	if got := recv(t, stay); got.ID != 1 {
		t.Fatalf("remaining subscriber got %d, want 1", got.ID)
	}

	b.mu.Lock()
	_, present := b.subs[leave]
	b.mu.Unlock()
	if present {
		t.Fatal("closed subscriber was not dropped")
	}
}
