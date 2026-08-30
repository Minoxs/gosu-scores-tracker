package tracker

import (
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

func (f *fakeFetcher) fetch(userID int, since time.Time) (gosu.FullScores, time.Time, error) {
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

func (f *fakeFetcher) seenSinces() []time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.sinces)
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
	tr := NewUserPoller(f.fetch, fastConfig())
	defer tr.Close()

	tr.Track(7, time.Now().Add(-time.Hour))

	if got := recvFull(t, tr.Scores()); got.UserID != 7 {
		t.Fatalf("emitted score for user %d, want 7", got.UserID)
	}
}

func TestUserPoller_AdvancesWatermark(t *testing.T) {
	f := &fakeFetcher{base: time.Now(), emit: true}
	tr := NewUserPoller(f.fetch, fastConfig())

	tr.Track(7, time.Now().Add(-time.Hour))
	recvFull(t, tr.Scores())
	recvFull(t, tr.Scores())
	tr.Close()

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
	tr := NewUserPoller(f.fetch, fastConfig())
	defer tr.Close()

	tr.Track(7, time.Now().Add(-time.Hour))
	recvFull(t, tr.Scores())
	tr.Untrack(7)

	assertQuiet(t, tr.Scores())
}

func TestUserPoller_TrackIsIdempotent(t *testing.T) {
	f := &fakeFetcher{base: time.Now(), emit: true}
	tr := NewUserPoller(f.fetch, fastConfig())
	defer tr.Close()

	since := time.Now().Add(-time.Hour)
	tr.Track(7, since)
	tr.Track(7, since) // duplicate must not start a second loop
	recvFull(t, tr.Scores())
	tr.Untrack(7) // a single untrack stops all polling of user 7

	assertQuiet(t, tr.Scores())
}

func TestUserPoller_ReportsChecks(t *testing.T) {
	type check struct {
		user    int
		checked time.Time
		next    time.Time
	}
	f := &fakeFetcher{base: time.Now()} // empty polls still report a check
	tr := NewUserPoller(f.fetch, fastConfig())
	defer tr.Close()

	checks := make(chan check, 8)
	tr.OnCheck = func(userID int, checked, next time.Time) {
		select {
		case checks <- check{userID, checked, next}:
		default:
		}
	}

	tr.Track(7, time.Now().Add(-time.Hour))

	select {
	case c := <-checks:
		if c.user != 7 {
			t.Fatalf("check for user %d, want 7", c.user)
		}
		if !c.next.After(c.checked) {
			t.Fatalf("next check %v not after checked %v", c.next, c.checked)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for a check report")
	}
}

func TestUserPoller_CloseClosesScores(t *testing.T) {
	f := &fakeFetcher{base: time.Now()}
	tr := NewUserPoller(f.fetch, fastConfig())

	tr.Track(7, time.Now().Add(-time.Hour))
	tr.Close()

	select {
	case _, ok := <-tr.Scores():
		if ok {
			t.Fatal("expected Scores channel to close")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for Scores to close")
	}
}
