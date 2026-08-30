package tracker

import (
	"context"
	"sync"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

// ScoreTracker narrows a feed to a chosen set of users, each from a start time.
// Track and Untrack adjust that set at any time; Scores yields the scores of
// tracked users set after their start.
type ScoreTracker interface {
	Track(userID int, since time.Time)
	Untrack(userID int)
	Scores() <-chan gosu.Score
}

// FilterTracker implements ScoreTracker over a ScoreProvider: it keeps a tracked
// set and forwards only scores whose user is tracked and which were set after
// that user's start time.
type FilterTracker struct {
	mu      sync.Mutex
	tracked map[int]time.Time
	out     chan gosu.Score
}

// NewFilterTracker subscribes to provider and starts forwarding matching scores.
// It stops when ctx is cancelled or the provider closes the subscription, at
// which point the Scores channel is closed.
func NewFilterTracker(ctx context.Context, provider ScoreProvider) *FilterTracker {
	t := &FilterTracker{
		tracked: make(map[int]time.Time),
		out:     make(chan gosu.Score, DefaultSubBuffer),
	}
	go t.run(ctx, provider.Subscribe())
	return t
}

func (t *FilterTracker) run(ctx context.Context, in <-chan gosu.Score) {
	defer close(t.out)
	for {
		select {
		case <-ctx.Done():
			return
		case s, ok := <-in:
			if !ok {
				return
			}
			if !t.wants(s) {
				continue
			}
			select {
			case t.out <- s:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (t *FilterTracker) wants(s gosu.Score) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	since, ok := t.tracked[s.UserID]
	return ok && s.EndedAt.After(since)
}

// Track starts forwarding userID's scores set after since. Calling it again with
// a new start moves the boundary.
func (t *FilterTracker) Track(userID int, since time.Time) {
	t.mu.Lock()
	t.tracked[userID] = since
	t.mu.Unlock()
}

// Untrack stops forwarding userID's scores.
func (t *FilterTracker) Untrack(userID int) {
	t.mu.Lock()
	delete(t.tracked, userID)
	t.mu.Unlock()
}

// Scores is the stream of tracked users' scores.
func (t *FilterTracker) Scores() <-chan gosu.Score { return t.out }
