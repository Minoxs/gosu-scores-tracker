package osu

import (
	"net/http"
	"sync"
	"time"
)

// defaultRequestsPerMinute is the osu! API terms-of-use ceiling: no more than 60
// requests per minute. Every request this package makes passes through the pacer
// so the whole process, across all callers, stays under it.
const defaultRequestsPerMinute = 60

// pacer spaces requests at least interval apart. Reserve hands each caller the
// earliest instant it may proceed, advancing the shared cursor so concurrent
// callers queue instead of bursting.
type pacer struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

func (p *pacer) reserve() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.next.Before(now) {
		p.next = now
	}
	wait := p.next.Sub(now)
	p.next = p.next.Add(p.interval)
	return wait
}

func (p *pacer) setRate(perMinute int) {
	if perMinute < 1 {
		perMinute = 1
	}
	p.mu.Lock()
	p.interval = time.Minute / time.Duration(perMinute)
	p.mu.Unlock()
}

// throttledTransport blocks each request until the pacer releases a slot.
type throttledTransport struct {
	base  http.RoundTripper
	pacer *pacer
}

func (t *throttledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if wait := t.pacer.reserve(); wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}
	return t.base.RoundTrip(req)
}

// globalPacer throttles apiClient. It is installed in config.go's init.
var globalPacer = &pacer{interval: time.Minute / defaultRequestsPerMinute}

// SetRateLimit caps every osu! request this package makes to perMinute requests
// per minute, shared process-wide. Values below 1 clamp to 1.
func SetRateLimit(perMinute int) {
	globalPacer.setRate(perMinute)
}
