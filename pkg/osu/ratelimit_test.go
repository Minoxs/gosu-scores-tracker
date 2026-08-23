package osu

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestPacerPrefersHighPriority verifies that a high-priority reservation enqueued
// after a low-priority one is still granted first.
func TestPacerPrefersHighPriority(t *testing.T) {
	p := newPacer(600) // 100ms between slots
	ctx := context.Background()

	// The first reservation takes the immediate slot and puts the dispatcher into
	// its inter-slot sleep, so the next two reservations queue behind one grant.
	if err := p.reserve(ctx, false); err != nil {
		t.Fatalf("first reserve: %v", err)
	}

	order := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_ = p.reserve(ctx, false)
		order <- "low"
	}()
	// Let the low reservation reach the queue before the high one, so priority,
	// not arrival order, decides.
	time.Sleep(20 * time.Millisecond)
	go func() {
		defer wg.Done()
		_ = p.reserve(ctx, true)
		order <- "high"
	}()

	wg.Wait()
	if first := <-order; first != "high" {
		t.Errorf("first granted = %q, want high", first)
	}
}

// TestPacerReserveHonorsContext verifies a cancelled context releases a waiter
// instead of blocking on a slot that never comes.
func TestPacerReserveHonorsContext(t *testing.T) {
	p := newPacer(1) // one slot per minute, so the second reservation waits

	if err := p.reserve(context.Background(), false); err != nil {
		t.Fatalf("first reserve: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := p.reserve(ctx, false); err == nil {
		t.Error("reserve with cancelled context returned nil, want context error")
	}
}
