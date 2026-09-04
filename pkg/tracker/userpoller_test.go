package tracker

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

// fakeFetcher scripts the osu! fetch: it records the since it is called with and,
// when emit is set, returns one score per call at a strictly increasing time.
type fakeFetcher struct {
	mu     sync.Mutex
	base   time.Time
	nextID int64
	sinces []time.Time
	emit   bool
}

func (f *fakeFetcher) fetch(userID int64, since time.Time) (gosu.FullScores, time.Time, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sinces = append(f.sinces, since)
	if !f.emit {
		return nil, since, nil
	}
	f.nextID++
	at := f.base.Add(time.Duration(f.nextID) * time.Second)
	return gosu.FullScores{{Score: gosu.Score{ID: f.nextID, UserID: userID, EndedAt: at}}}, at, nil
}

func (f *fakeFetcher) seenSinces() []time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.sinces)
}

// recvFull reads one full score off the poller stream, failing on timeout.
func recvFull(t *testing.T, ch <-chan gosu.FullScore) gosu.FullScore {
	t.Helper()
	select {
	case s := <-ch:
		return s
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for score")
		return gosu.FullScore{}
	}
}

const fastPoll = time.Millisecond

func fastConfig() PollConfig {
	return PollConfig{BaseInterval: fastPoll, MaxInterval: fastPoll}
}

// assertQuiet drains anything in flight for a grace period, then fails if a
// further score arrives, proving polling has stopped.
func assertQuiet(t *testing.T, ch <-chan gosu.FullScore) {
	t.Helper()
	grace := time.After(50 * time.Millisecond)
	for draining := true; draining; {
		select {
		case <-ch:
		case <-grace:
			draining = false
		}
	}
	select {
	case s, ok := <-ch:
		if ok {
			t.Fatalf("still emitting after stop: score %d", s.ID)
		}
	case <-time.After(100 * time.Millisecond):
	}
}

func TestUserPoller_EmitsScores(t *testing.T) {
	f := &fakeFetcher{base: time.Now(), emit: true}
	tr := newUserPoller(f.fetch, fastConfig())

	sub := tr.Subscribe()
	tr.Track(7, time.Now().Add(-time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go tr.Run(ctx)

	if got := recvFull(t, sub); got.UserID != 7 {
		t.Fatalf("emitted score for user %d, want 7", got.UserID)
	}
}

func TestUserPoller_AdvancesWatermark(t *testing.T) {
	f := &fakeFetcher{base: time.Now(), emit: true}
	tr := newUserPoller(f.fetch, fastConfig())

	sub := tr.Subscribe()
	tr.Track(7, time.Now().Add(-time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	go tr.Run(ctx)

	recvFull(t, sub)
	recvFull(t, sub)
	cancel()

	sinces := f.seenSinces()
	if len(sinces) < 2 {
		t.Fatalf("fetch called %d times, want at least 2", len(sinces))
	}
	if !sinces[1].After(sinces[0]) {
		t.Fatalf("watermark did not advance: %v then %v", sinces[0], sinces[1])
	}
}

func TestUserPoller_Untrack(t *testing.T) {
	f := &fakeFetcher{base: time.Now(), emit: true}
	tr := newUserPoller(f.fetch, fastConfig())

	sub := tr.Subscribe()
	tr.Track(7, time.Now().Add(-time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go tr.Run(ctx)

	recvFull(t, sub)
	tr.Untrack(7)

	assertQuiet(t, sub)
}

func TestUserPoller_TrackIsIdempotent(t *testing.T) {
	f := &fakeFetcher{base: time.Now(), emit: true}
	tr := newUserPoller(f.fetch, fastConfig())

	sub := tr.Subscribe()
	since := time.Now().Add(-time.Hour)
	tr.Track(7, since)
	tr.Track(7, since) // duplicate must not start a second loop

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go tr.Run(ctx)

	recvFull(t, sub)
	tr.Untrack(7) // a single untrack stops all polling of user 7

	assertQuiet(t, sub)
}

func TestUserPoller_ShutdownClosesStreams(t *testing.T) {
	f := &fakeFetcher{base: time.Now()}
	tr := newUserPoller(f.fetch, fastConfig())

	sub := tr.Subscribe()
	tr.Track(7, time.Now().Add(-time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	go tr.Run(ctx)
	cancel()

	select {
	case _, ok := <-sub:
		if ok {
			t.Fatal("expected Subscribe channel to close on shutdown")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for stream to close")
	}
}
