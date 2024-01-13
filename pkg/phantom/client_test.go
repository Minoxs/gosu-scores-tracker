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

func TestClient_ProcessNewScores(t *testing.T) {
	var scoreCount int

	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			scoreCount = len(scores)
		},
	}

	test.processNewScores(testScores)
	if scoreCount != 4 {
		t.Fatal("Score count mismatch")
	}
}

func TestClient_Ranking(t *testing.T) {
	var scoreCount int

	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			scoreCount = len(scores)
		},
	}

	test.processNewScores(testScores)
	if scoreCount != test.Ranking().Count()+1 {
		t.Fatal("Score count mismatch")
	}
}

func BenchmarkClient_ProcessNewScores(b *testing.B) {
	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			_ = len(scores)
		},
	}
	var ref = test.LastUpdate

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		test.LastUpdate = ref
		test.processNewScores(testScores)
	}
}

func BenchmarkClient_Ranking(b *testing.B) {
	var test = &Client{
		Logger: slog.Default(),
		OnNewScores: func(scores []NewScore) {
			_ = len(scores)
		},
	}
	test.processNewScores(testScores)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = test.Ranking().Count()
	}
}
