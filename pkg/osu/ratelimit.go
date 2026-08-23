package osu

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// defaultRequestsPerMinute is the osu! API terms-of-use ceiling: no more than 60
// requests per minute. Every request this package makes passes through the pacer
// so the whole process, across all callers, stays under it.
const defaultRequestsPerMinute = 60

// priorityKey marks a request context as interactive so it reserves ahead of
// background traffic at the shared rate limit.
type priorityKey struct{}

// prioritize tags req as interactive, so the pacer serves it before background
// requests such as the global feed tail and its beatmap enrichment.
func prioritize(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), priorityKey{}, true))
}

// pacer spaces osu! requests to keep the whole process under the API ceiling,
// granting one slot per interval. Interactive reservations (a user lookup behind a
// page load) are served before background ones (the global feed tail and its
// enrichment), so a burst of background polling never delays an interactive
// request by more than a single slot.
type pacer struct {
	high  chan chan struct{}
	low   chan chan struct{}
	start sync.Once

	mu       sync.Mutex
	interval time.Duration
}

func newPacer(perMinute int) *pacer {
	return &pacer{
		high:     make(chan chan struct{}),
		low:      make(chan chan struct{}),
		interval: rateInterval(perMinute),
	}
}

func rateInterval(perMinute int) time.Duration {
	if perMinute < 1 {
		perMinute = 1
	}
	return time.Minute / time.Duration(perMinute)
}

// reserve blocks until the pacer grants a slot or ctx is cancelled. A high-priority
// reservation jumps ahead of any waiting low-priority ones.
func (p *pacer) reserve(ctx context.Context, high bool) error {
	p.start.Do(func() { go p.run() })

	slot := make(chan struct{})
	queue := p.low
	if high {
		queue = p.high
	}
	select {
	case queue <- slot:
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case <-slot:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// run grants one slot per interval, always preferring a waiting high-priority
// reservation, and idles at no cost when nobody is waiting.
func (p *pacer) run() {
	for {
		close(p.nextWaiter())
		p.mu.Lock()
		interval := p.interval
		p.mu.Unlock()
		time.Sleep(interval)
	}
}

// nextWaiter blocks until a reservation is waiting and returns the one to grant,
// draining a high-priority waiter first.
func (p *pacer) nextWaiter() chan struct{} {
	select {
	case s := <-p.high:
		return s
	default:
	}
	select {
	case s := <-p.high:
		return s
	case s := <-p.low:
		return s
	}
}

func (p *pacer) setRate(perMinute int) {
	p.mu.Lock()
	p.interval = rateInterval(perMinute)
	p.mu.Unlock()
}

// throttledTransport blocks each request until the pacer grants it a slot, at the
// priority carried on the request context.
type throttledTransport struct {
	base  http.RoundTripper
	pacer *pacer
}

func (t *throttledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	high, _ := req.Context().Value(priorityKey{}).(bool)
	if err := t.pacer.reserve(req.Context(), high); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

// globalPacer throttles apiClient. It is installed in config.go.
var globalPacer = newPacer(defaultRequestsPerMinute)

// SetRateLimit caps every osu! request this package makes to perMinute requests
// per minute, shared process-wide. Values below 1 clamp to 1.
func SetRateLimit(perMinute int) {
	globalPacer.setRate(perMinute)
}
