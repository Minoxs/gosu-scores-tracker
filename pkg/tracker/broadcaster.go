package tracker

import (
	"context"
	"sync"
)

// subBuffer is the per-subscriber channel capacity. Past it emit blocks, so a slow
// consumer slows polling rather than losing scores
const subBuffer = 256

// broadcaster fans a value out to every subscriber, blocking on each until it is
// taken so nothing is dropped. A subscriber unsubscribes by closing its channel:
// the send then panics and the recover drops it
type broadcaster[T any] struct {
	mu   sync.Mutex
	subs map[chan T]struct{}
}

func newBroadcaster[T any]() *broadcaster[T] {
	return &broadcaster[T]{subs: make(map[chan T]struct{})}
}

func (b *broadcaster[T]) subscribe() chan T {
	ch := make(chan T, subBuffer)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

// emit blocks until every subscriber has taken v, so a full buffer slows the caller
// instead of dropping the value. It abandons a send only when ctx is cancelled
func (b *broadcaster[T]) emit(ctx context.Context, v T) {
	for _, ch := range b.snapshot() {
		b.send(ctx, ch, v)
	}
}

// snapshot copies the current subscribers so emit can block on a send without
// holding the lock, which would otherwise stall subscribe, close, and unsubscribe
func (b *broadcaster[T]) snapshot() []chan T {
	b.mu.Lock()
	defer b.mu.Unlock()
	subs := make([]chan T, 0, len(b.subs))
	for ch := range b.subs {
		subs = append(subs, ch)
	}
	return subs
}

// send blocks until ch takes v or ctx is cancelled. A closed channel means the
// consumer unsubscribed, so the send panics and the recover drops the subscriber
func (b *broadcaster[T]) send(ctx context.Context, ch chan T, v T) {
	defer func() {
		if recover() != nil {
			b.remove(ch)
		}
	}()
	select {
	case ch <- v:
	case <-ctx.Done():
	}
}

func (b *broadcaster[T]) remove(ch chan T) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
}

// close closes every remaining subscriber so a ranging consumer ends. It recovers
// per channel in case a consumer already closed its own to unsubscribe
func (b *broadcaster[T]) close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.subs {
		safeClose(ch)
		delete(b.subs, ch)
	}
}

func safeClose[T any](ch chan T) {
	defer func() { _ = recover() }()
	close(ch)
}
