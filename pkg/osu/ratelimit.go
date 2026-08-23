package osu

import (
	"container/heap"
	"context"
	"net/http"
	"sync"
	"time"
)

// defaultRequestsPerMinute is the osu! API terms-of-use ceiling: no more than 60
// requests per minute. Every request this package makes passes through the pacer
// so the whole process, across all callers, stays under it.
const defaultRequestsPerMinute = 60

// Priority orders requests at the shared pacer: when several wait at once, a higher
// value is granted first, and ties break by arrival order. This package only orders
// the levels; the caller assigns their meaning by building a Client at each level.
// The zero value is the level a default Client carries.
type Priority int

// waiter is one pending reservation; ready closes when the pacer grants it.
type waiter struct {
	prio  Priority
	seq   uint64
	ready chan struct{}
}

// waiterHeap orders by descending priority, then ascending seq so a level is served
// first-come first-served.
type waiterHeap []*waiter

func (h waiterHeap) Len() int { return len(h) }
func (h waiterHeap) Less(i, j int) bool {
	if h[i].prio != h[j].prio {
		return h[i].prio > h[j].prio
	}
	return h[i].seq < h[j].seq
}
func (h waiterHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *waiterHeap) Push(x any)   { *h = append(*h, x.(*waiter)) }
func (h *waiterHeap) Pop() any {
	old := *h
	n := len(old)
	w := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return w
}

// pacer spaces osu! requests to keep the whole process under the API ceiling,
// granting one slot per interval to the highest-priority waiter. A caller reserves
// at a priority it chooses, so a burst of low-priority traffic never delays a
// higher-priority request by more than a single slot.
type pacer struct {
	mu       sync.Mutex
	cond     *sync.Cond
	queue    waiterHeap
	seq      uint64
	interval time.Duration
	started  bool
}

func newPacer(perMinute int) *pacer {
	p := &pacer{interval: rateInterval(perMinute)}
	p.cond = sync.NewCond(&p.mu)
	return p
}

func rateInterval(perMinute int) time.Duration {
	if perMinute < 1 {
		perMinute = 1
	}
	return time.Minute / time.Duration(perMinute)
}

// reserve blocks until the pacer grants a slot or ctx is cancelled. A higher
// priority is granted ahead of any lower one waiting at the same time.
func (p *pacer) reserve(ctx context.Context, prio Priority) error {
	w := &waiter{prio: prio, ready: make(chan struct{})}

	p.mu.Lock()
	if !p.started {
		p.started = true
		go p.run()
	}
	w.seq = p.seq
	p.seq++
	heap.Push(&p.queue, w)
	p.cond.Signal()
	p.mu.Unlock()

	select {
	case <-w.ready:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// run grants one waiter per interval, always the highest priority waiting.
func (p *pacer) run() {
	for {
		p.mu.Lock()
		for p.queue.Len() == 0 {
			p.cond.Wait()
		}
		w := heap.Pop(&p.queue).(*waiter)
		interval := p.interval
		p.mu.Unlock()

		close(w.ready)
		time.Sleep(interval)
	}
}

func (p *pacer) setRate(perMinute int) {
	p.mu.Lock()
	p.interval = rateInterval(perMinute)
	p.mu.Unlock()
}

// throttledTransport blocks each request until the pacer grants it a slot at prio,
// the level of the Client that owns this transport. The request context is honored
// only for cancellation.
type throttledTransport struct {
	base  http.RoundTripper
	pacer *pacer
	prio  Priority
}

func (t *throttledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.pacer.reserve(req.Context(), t.prio); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

// globalPacer is the process-wide pacer every Client reserves on, so the rate ceiling
// holds across all callers and priority levels.
var globalPacer = newPacer(defaultRequestsPerMinute)

// SetRateLimit caps every osu! request this package makes to perMinute requests
// per minute, shared process-wide. Values below 1 clamp to 1.
func SetRateLimit(perMinute int) {
	globalPacer.setRate(perMinute)
}
