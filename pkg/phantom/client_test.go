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
		CreatedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 1},
	},
	{
		ID:        2,
		PP:        110,
		CreatedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 1},
	},
	{
		ID:        3,
		PP:        90,
		CreatedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 1},
	},
	{
		ID:        4,
		PP:        200,
		CreatedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 2},
	},
	{
		ID:        5,
		PP:        300,
		CreatedAt: time.Now(),
		Beatmap:   player.Beatmap{ID: 3},
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
