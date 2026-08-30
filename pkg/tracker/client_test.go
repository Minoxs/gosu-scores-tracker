package tracker

import (
	"log/slog"
	"testing"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
)

func full(id int64, beatmapID int64, pp float64) gosu.FullScore {
	return gosu.FullScore{
		Score:   gosu.Score{ID: id, BeatmapID: beatmapID, PP: pp, EndedAt: time.Now()},
		Beatmap: gosu.Beatmap{ID: beatmapID},
	}
}

var testScores = gosu.FullScores{
	full(1, 1, 100),
	full(2, 1, 110),
	full(3, 1, 90),
	full(4, 2, 200),
	full(5, 3, 300),
}

func TestClient_FoldPage(t *testing.T) {
	var counted = make(chan int, 1)

	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			counted <- len(scores)
		},
	}

	test.foldPage(testScores, time.Time{})
	if scoreCount := <-counted; scoreCount != 4 {
		t.Fatal("Score count mismatch")
	}
}

func TestClient_Ranking(t *testing.T) {
	var counted = make(chan int, 1)

	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			counted <- len(scores)
		},
	}

	test.foldPage(testScores, time.Time{})
	if scoreCount := <-counted; scoreCount != test.Ranking().Count()+1 {
		t.Fatal("Score count mismatch")
	}
}

// A score with no pp is not ranked.
func TestClient_FoldPageSkipsNoPP(t *testing.T) {
	var mixed = gosu.FullScores{
		full(10, 10, 100),
		full(11, 11, 0),
		full(12, 12, 0),
		full(13, 13, 140),
		full(14, 14, 0),
	}

	var counted = make(chan int, 1)
	var test = &Client{
		Logger:      slog.Default(),
		OnNewScores: func(scores []NewScore) { counted <- len(scores) },
	}

	test.foldPage(mixed, time.Time{})
	if got := <-counted; got != 2 {
		t.Fatalf("counted %d scores with pp, want 2", got)
	}
	if got := test.Ranking().Count(); got != 2 {
		t.Fatalf("ranking holds %d scores, want 2", got)
	}
}

func TestClient_RestoreSkipsNoPP(t *testing.T) {
	var mixed = gosu.Scores{
		{ID: 20, BeatmapID: 20, PP: 100, EndedAt: time.Now()},
		{ID: 21, BeatmapID: 21, PP: 0, EndedAt: time.Now()},
	}

	var test = &Client{Logger: slog.Default()}
	test.Restore(mixed)
	if got := test.Ranking().Count(); got != 1 {
		t.Fatalf("ranking holds %d scores after restore, want 1", got)
	}
}

func BenchmarkClient_FoldPage(b *testing.B) {
	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			_ = len(scores)
		},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		test.foldPage(testScores, time.Time{})
	}
}

func BenchmarkClient_Ranking(b *testing.B) {
	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			_ = len(scores)
		},
	}
	test.foldPage(testScores, time.Time{})

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = test.Ranking().Count()
	}
}
