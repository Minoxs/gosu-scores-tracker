package phantom

import (
	"github.com/minoxs/osu-phantom/pkg/osu/player"
	"log/slog"
	"testing"
	"time"
)

var testScores = player.Scores{
	{
		ID:        1,
		PP:        100,
		EndedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 1, Status: "ranked"},
	},
	{
		ID:        2,
		PP:        110,
		EndedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 1, Status: "ranked"},
	},
	{
		ID:        3,
		PP:        90,
		EndedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 1, Status: "approved"},
	},
	{
		ID:        4,
		PP:        200,
		EndedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 2, Status: "ranked"},
	},
	{
		ID:        5,
		PP:        300,
		EndedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 3, Status: "approved"},
	},
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

func TestClient_FoldPageSkipsUnranked(t *testing.T) {
	var mixed = player.Scores{
		{ID: 10, PP: 100, EndedAt: time.Now(), Beatmap: player.Beatmap{ID: 10, Status: "ranked"}},
		{ID: 11, PP: 120, EndedAt: time.Now(), Beatmap: player.Beatmap{ID: 11, Status: "loved"}},
		{ID: 12, PP: 130, EndedAt: time.Now(), Beatmap: player.Beatmap{ID: 12, Status: "graveyard"}},
		{ID: 13, PP: 140, EndedAt: time.Now(), Beatmap: player.Beatmap{ID: 13, Status: "approved"}},
		{ID: 14, PP: 150, EndedAt: time.Now(), Beatmap: player.Beatmap{ID: 14, Status: "qualified"}},
	}

	var counted = make(chan int, 1)
	var test = &Client{
		Logger:      slog.Default(),
		OnNewScores: func(scores []NewScore) { counted <- len(scores) },
	}

	test.foldPage(mixed, time.Time{})
	if got := <-counted; got != 2 {
		t.Fatalf("counted %d ranked scores, want 2", got)
	}
	if got := test.Ranking().Count(); got != 2 {
		t.Fatalf("ranking holds %d scores, want 2", got)
	}
}

func TestClient_RestoreSkipsUnranked(t *testing.T) {
	var mixed = player.Scores{
		{ID: 20, PP: 100, EndedAt: time.Now(), Beatmap: player.Beatmap{ID: 20, Status: "ranked"}},
		{ID: 21, PP: 120, EndedAt: time.Now(), Beatmap: player.Beatmap{ID: 21, Status: "loved"}},
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
