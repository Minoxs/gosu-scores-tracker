package tracker

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

type fetchResult struct {
	scores gosu.Scores
	cursor string
	err    error
}

// scriptedFetcher returns scripted pages in order, then reports caught-up (no
// scores, holding the last cursor) once the script is spent. It records the cursor
// each call was made with.
type scriptedFetcher struct {
	mu      sync.Mutex
	steps   []fetchResult
	n       int
	cursors []string
}

func (s *scriptedFetcher) fetch(cursor string) (gosu.Scores, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cursors = append(s.cursors, cursor)
	if s.n < len(s.steps) {
		r := s.steps[s.n]
		s.n++
		return r.scores, r.cursor, r.err
	}
	last := ""
	if len(s.steps) > 0 {
		last = s.steps[len(s.steps)-1].cursor
	}
	return nil, last, nil
}

func (s *scriptedFetcher) seenCursors() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.cursors)
}

func TestRealtimePoller_SeedsCursorWithoutEmittingNewestPage(t *testing.T) {
	sf := &scriptedFetcher{steps: []fetchResult{
		{scores: gosu.Scores{{ID: 999}}, cursor: "c0"}, // seed page, must not emit
		{scores: gosu.Scores{{ID: 1}}, cursor: "c1"},   // first streamed page
	}}
	p := NewRealtimePoller(sf.fetch, time.Millisecond)
	sub := p.Subscribe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	if got := recv(t, sub); got.ID != 1 {
		t.Fatalf("first emitted score id %d, want 1 (seed page must be skipped)", got.ID)
	}
	assertNoScore(t, sub)

	cursors := sf.seenCursors()
	if len(cursors) < 2 || cursors[0] != "" || cursors[1] != "c0" {
		t.Fatalf("cursors %v, want seed \"\" then stream \"c0\"", cursors[:min(len(cursors), 2)])
	}
}

func TestRealtimePoller_ResumeStreamsFromCursorAndAdvances(t *testing.T) {
	sf := &scriptedFetcher{steps: []fetchResult{
		{scores: gosu.Scores{{ID: 7}}, cursor: "c1"}, // first streamed page from the resumed cursor
	}}
	p := NewRealtimePoller(sf.fetch, time.Millisecond)
	p.Resume("saved")
	sub := p.Subscribe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	if got := recv(t, sub); got.ID != 7 {
		t.Fatalf("first emitted score id %d, want 7 (resume must not seed-skip a page)", got.ID)
	}
	if first := sf.seenCursors()[0]; first != "saved" {
		t.Fatalf("first fetch cursor %q, want resumed \"saved\"", first)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && p.Cursor() != "c1" {
		time.Sleep(time.Millisecond)
	}
	if p.Cursor() != "c1" {
		t.Fatalf("Cursor() = %q, want advanced to c1", p.Cursor())
	}
}

func TestRealtimePoller_RetriesSeedOnError(t *testing.T) {
	sf := &scriptedFetcher{steps: []fetchResult{
		{err: errors.New("boom")},                    // seed fails
		{cursor: "c0"},                               // seed retry succeeds
		{scores: gosu.Scores{{ID: 5}}, cursor: "c1"}, // streams
	}}
	p := NewRealtimePoller(sf.fetch, time.Millisecond)
	sub := p.Subscribe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Run(ctx)

	if got := recv(t, sub); got.ID != 5 {
		t.Fatalf("emitted score id %d, want 5 after seed retry", got.ID)
	}
}
