package player

import (
	"testing"
)

func assertEqual[T comparable](t *testing.T, val1 T, val2 T) {
	if val1 != val2 {
		t.Fatalf("Values are not equal\nExpected=%v\nActual=%v\n", val1, val2)
	}
}

func TestRanking_AddScore(t *testing.T) {
	panic("TODO")
}

func TestRanking_GetTotalPPSingle(t *testing.T) {
	const ExpectedTotalPP = float64(100)
	var rank = Ranking{
		count: 1,
		scores: [RankSize]Score{
			{PP: ExpectedTotalPP},
		},
	}

	assertEqual(t, ExpectedTotalPP, rank.GetTotalPP())
}

func TestRanking_GetTotalPPMultiple(t *testing.T) {
	const ExpectedTotalPP = float64(336 + 316 + 298 + 282 + 267)
	var rank = Ranking{
		count: 5,
		scores: [RankSize]Score{
			{PP: 336.242},
			{PP: 332.834},
			{PP: 330.403},
			{PP: 328.735},
			{PP: 328.239},
		},
	}

	assertEqual(t, ExpectedTotalPP, rank.GetTotalPP())
}
