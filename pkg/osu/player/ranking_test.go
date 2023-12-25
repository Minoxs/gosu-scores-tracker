package player

import (
	"log"
	"testing"
)

type TestCaseInfo struct {
	Name          string
	Scores        Scores
	ExpectedCount int
}

func assertEqual[T comparable](t *testing.T, val1 T, val2 T) {
	if val1 != val2 {
		t.Fatalf("Values are not equal\nExpected=%v\nActual=%v\n", val1, val2)
	}
}

func getTestCases() []TestCaseInfo {
	return []TestCaseInfo{
		{
			Name: "Single",
			Scores: Scores{
				{Beatmap: Beatmap{ID: 1}, PP: 250.0},
			},
			ExpectedCount: 1,
		},
		{
			Name: "Triple",
			Scores: Scores{
				{Beatmap: Beatmap{ID: 1}, PP: 369.15}, {Beatmap: Beatmap{ID: 2}, PP: 405.17}, {Beatmap: Beatmap{ID: 3}, PP: 5.89},
			},
			ExpectedCount: 3,
		},
		{
			Name: "Six scores",
			Scores: Scores{
				{Beatmap: Beatmap{ID: 1}, PP: 690.64}, {Beatmap: Beatmap{ID: 2}, PP: 216.06}, {Beatmap: Beatmap{ID: 3}, PP: 141.08},
				{Beatmap: Beatmap{ID: 4}, PP: 6.03}, {Beatmap: Beatmap{ID: 5}, PP: 14.58}, {Beatmap: Beatmap{ID: 6}, PP: 390.45},
			},
			ExpectedCount: 6,
		},
		{
			Name: "Big ranking",
			Scores: Scores{
				{Beatmap: Beatmap{ID: 1}, PP: 132.05}, {Beatmap: Beatmap{ID: 2}, PP: 87.38}, {Beatmap: Beatmap{ID: 3}, PP: 72.42}, {Beatmap: Beatmap{ID: 4}, PP: 1.12},
				{Beatmap: Beatmap{ID: 5}, PP: 291.75}, {Beatmap: Beatmap{ID: 6}, PP: 34.44}, {Beatmap: Beatmap{ID: 7}, PP: 251.46}, {Beatmap: Beatmap{ID: 8}, PP: 38.55},
				{Beatmap: Beatmap{ID: 9}, PP: 536.43}, {Beatmap: Beatmap{ID: 10}, PP: 323.23}, {Beatmap: Beatmap{ID: 11}, PP: 45.86}, {Beatmap: Beatmap{ID: 12}, PP: 414.44},
				{Beatmap: Beatmap{ID: 13}, PP: 587.47}, {Beatmap: Beatmap{ID: 14}, PP: 202.37}, {Beatmap: Beatmap{ID: 15}, PP: 427.90}, {Beatmap: Beatmap{ID: 16}, PP: 243.74},
				{Beatmap: Beatmap{ID: 17}, PP: 34.29}, {Beatmap: Beatmap{ID: 18}, PP: 97.35}, {Beatmap: Beatmap{ID: 19}, PP: 171.57}, {Beatmap: Beatmap{ID: 20}, PP: 295.12},
				{Beatmap: Beatmap{ID: 21}, PP: 235.46}, {Beatmap: Beatmap{ID: 22}, PP: 26.26}, {Beatmap: Beatmap{ID: 23}, PP: 685.67}, {Beatmap: Beatmap{ID: 24}, PP: 302.57},
				{Beatmap: Beatmap{ID: 25}, PP: 201.04},
			},
			ExpectedCount: 25,
		},
		{
			Name: "Exactly 100",
			Scores: Scores{
				{Beatmap: Beatmap{ID: 1}, PP: 155.34}, {Beatmap: Beatmap{ID: 2}, PP: 253.52}, {Beatmap: Beatmap{ID: 3}, PP: 247.91}, {Beatmap: Beatmap{ID: 4}, PP: 0.00},
				{Beatmap: Beatmap{ID: 5}, PP: 42.61}, {Beatmap: Beatmap{ID: 6}, PP: 205.86}, {Beatmap: Beatmap{ID: 7}, PP: 170.99}, {Beatmap: Beatmap{ID: 8}, PP: 137.15},
				{Beatmap: Beatmap{ID: 9}, PP: 376.94}, {Beatmap: Beatmap{ID: 10}, PP: 607.78}, {Beatmap: Beatmap{ID: 11}, PP: 0.54}, {Beatmap: Beatmap{ID: 12}, PP: 68.71},
				{Beatmap: Beatmap{ID: 13}, PP: 96.11}, {Beatmap: Beatmap{ID: 14}, PP: 172.66}, {Beatmap: Beatmap{ID: 15}, PP: 113.40}, {Beatmap: Beatmap{ID: 16}, PP: 250.38},
				{Beatmap: Beatmap{ID: 17}, PP: 299.31}, {Beatmap: Beatmap{ID: 18}, PP: 467.02}, {Beatmap: Beatmap{ID: 19}, PP: 87.90}, {Beatmap: Beatmap{ID: 20}, PP: 555.12},
				{Beatmap: Beatmap{ID: 21}, PP: 347.32}, {Beatmap: Beatmap{ID: 22}, PP: 3.51}, {Beatmap: Beatmap{ID: 23}, PP: 230.51}, {Beatmap: Beatmap{ID: 24}, PP: 419.33},
				{Beatmap: Beatmap{ID: 25}, PP: 322.57}, {Beatmap: Beatmap{ID: 26}, PP: 32.23}, {Beatmap: Beatmap{ID: 27}, PP: 882.48}, {Beatmap: Beatmap{ID: 28}, PP: 117.55},
				{Beatmap: Beatmap{ID: 29}, PP: 86.10}, {Beatmap: Beatmap{ID: 30}, PP: 15.07}, {Beatmap: Beatmap{ID: 31}, PP: 428.13}, {Beatmap: Beatmap{ID: 32}, PP: 189.02},
				{Beatmap: Beatmap{ID: 33}, PP: 515.16}, {Beatmap: Beatmap{ID: 34}, PP: 122.10}, {Beatmap: Beatmap{ID: 35}, PP: 220.60}, {Beatmap: Beatmap{ID: 36}, PP: 68.68},
				{Beatmap: Beatmap{ID: 37}, PP: 15.34}, {Beatmap: Beatmap{ID: 38}, PP: 589.47}, {Beatmap: Beatmap{ID: 39}, PP: 225.46}, {Beatmap: Beatmap{ID: 40}, PP: 184.98},
				{Beatmap: Beatmap{ID: 41}, PP: 114.67}, {Beatmap: Beatmap{ID: 42}, PP: 451.38}, {Beatmap: Beatmap{ID: 43}, PP: 122.08}, {Beatmap: Beatmap{ID: 44}, PP: 15.76},
				{Beatmap: Beatmap{ID: 45}, PP: 538.66}, {Beatmap: Beatmap{ID: 46}, PP: 401.96}, {Beatmap: Beatmap{ID: 47}, PP: 105.42}, {Beatmap: Beatmap{ID: 48}, PP: 386.58},
				{Beatmap: Beatmap{ID: 49}, PP: 378.19}, {Beatmap: Beatmap{ID: 50}, PP: 540.71}, {Beatmap: Beatmap{ID: 51}, PP: 33.16}, {Beatmap: Beatmap{ID: 52}, PP: 289.03},
				{Beatmap: Beatmap{ID: 53}, PP: 226.25}, {Beatmap: Beatmap{ID: 54}, PP: 498.74}, {Beatmap: Beatmap{ID: 55}, PP: 74.20}, {Beatmap: Beatmap{ID: 56}, PP: 70.60},
				{Beatmap: Beatmap{ID: 57}, PP: 71.17}, {Beatmap: Beatmap{ID: 58}, PP: 237.44}, {Beatmap: Beatmap{ID: 59}, PP: 248.87}, {Beatmap: Beatmap{ID: 60}, PP: 749.55},
				{Beatmap: Beatmap{ID: 61}, PP: 18.85}, {Beatmap: Beatmap{ID: 62}, PP: 57.17}, {Beatmap: Beatmap{ID: 63}, PP: 115.54}, {Beatmap: Beatmap{ID: 64}, PP: 578.59},
				{Beatmap: Beatmap{ID: 65}, PP: 660.24}, {Beatmap: Beatmap{ID: 66}, PP: 55.10}, {Beatmap: Beatmap{ID: 67}, PP: 110.93}, {Beatmap: Beatmap{ID: 68}, PP: 50.04},
				{Beatmap: Beatmap{ID: 69}, PP: 236.64}, {Beatmap: Beatmap{ID: 70}, PP: 328.36}, {Beatmap: Beatmap{ID: 71}, PP: 442.37}, {Beatmap: Beatmap{ID: 72}, PP: 42.97},
				{Beatmap: Beatmap{ID: 73}, PP: 120.55}, {Beatmap: Beatmap{ID: 74}, PP: 254.91}, {Beatmap: Beatmap{ID: 75}, PP: 142.38}, {Beatmap: Beatmap{ID: 76}, PP: 70.01},
				{Beatmap: Beatmap{ID: 77}, PP: 3.38}, {Beatmap: Beatmap{ID: 78}, PP: 119.79}, {Beatmap: Beatmap{ID: 79}, PP: 411.25}, {Beatmap: Beatmap{ID: 80}, PP: 511.07},
				{Beatmap: Beatmap{ID: 81}, PP: 83.76}, {Beatmap: Beatmap{ID: 82}, PP: 206.13}, {Beatmap: Beatmap{ID: 83}, PP: 103.37}, {Beatmap: Beatmap{ID: 84}, PP: 256.97},
				{Beatmap: Beatmap{ID: 85}, PP: 28.96}, {Beatmap: Beatmap{ID: 86}, PP: 403.87}, {Beatmap: Beatmap{ID: 87}, PP: 205.81}, {Beatmap: Beatmap{ID: 88}, PP: 308.41},
				{Beatmap: Beatmap{ID: 89}, PP: 8.43}, {Beatmap: Beatmap{ID: 90}, PP: 48.51}, {Beatmap: Beatmap{ID: 91}, PP: 351.35}, {Beatmap: Beatmap{ID: 92}, PP: 74.64},
				{Beatmap: Beatmap{ID: 93}, PP: 44.60}, {Beatmap: Beatmap{ID: 94}, PP: 39.09}, {Beatmap: Beatmap{ID: 95}, PP: 12.92}, {Beatmap: Beatmap{ID: 96}, PP: 33.84},
				{Beatmap: Beatmap{ID: 97}, PP: 518.59}, {Beatmap: Beatmap{ID: 98}, PP: 129.40}, {Beatmap: Beatmap{ID: 99}, PP: 766.85}, {Beatmap: Beatmap{ID: 100}, PP: 26.06},
			},
			ExpectedCount: 100,
		},
		{
			Name: "More than a full ranking",
			Scores: Scores{
				{Beatmap: Beatmap{ID: 1}, PP: 166.14}, {Beatmap: Beatmap{ID: 2}, PP: 638.01}, {Beatmap: Beatmap{ID: 3}, PP: 732.72}, {Beatmap: Beatmap{ID: 4}, PP: 212.14},
				{Beatmap: Beatmap{ID: 5}, PP: 262.47}, {Beatmap: Beatmap{ID: 6}, PP: 394.15}, {Beatmap: Beatmap{ID: 7}, PP: 25.37}, {Beatmap: Beatmap{ID: 8}, PP: 395.60},
				{Beatmap: Beatmap{ID: 9}, PP: 74.06}, {Beatmap: Beatmap{ID: 10}, PP: 733.72}, {Beatmap: Beatmap{ID: 11}, PP: 114.69}, {Beatmap: Beatmap{ID: 12}, PP: 88.51},
				{Beatmap: Beatmap{ID: 13}, PP: 305.42}, {Beatmap: Beatmap{ID: 14}, PP: 215.11}, {Beatmap: Beatmap{ID: 15}, PP: 477.38}, {Beatmap: Beatmap{ID: 16}, PP: 134.25},
				{Beatmap: Beatmap{ID: 17}, PP: 211.30}, {Beatmap: Beatmap{ID: 18}, PP: 436.71}, {Beatmap: Beatmap{ID: 19}, PP: 380.79}, {Beatmap: Beatmap{ID: 20}, PP: 184.20},
				{Beatmap: Beatmap{ID: 21}, PP: 226.27}, {Beatmap: Beatmap{ID: 22}, PP: 260.12}, {Beatmap: Beatmap{ID: 23}, PP: 92.39}, {Beatmap: Beatmap{ID: 24}, PP: 89.83},
				{Beatmap: Beatmap{ID: 25}, PP: 84.89}, {Beatmap: Beatmap{ID: 26}, PP: 653.27}, {Beatmap: Beatmap{ID: 27}, PP: 263.21}, {Beatmap: Beatmap{ID: 28}, PP: 92.73},
				{Beatmap: Beatmap{ID: 29}, PP: 175.19}, {Beatmap: Beatmap{ID: 30}, PP: 144.43}, {Beatmap: Beatmap{ID: 31}, PP: 362.05}, {Beatmap: Beatmap{ID: 32}, PP: 19.02},
				{Beatmap: Beatmap{ID: 33}, PP: 10.85}, {Beatmap: Beatmap{ID: 34}, PP: 360.89}, {Beatmap: Beatmap{ID: 35}, PP: 99.82}, {Beatmap: Beatmap{ID: 36}, PP: 146.40},
				{Beatmap: Beatmap{ID: 37}, PP: 18.29}, {Beatmap: Beatmap{ID: 38}, PP: 96.41}, {Beatmap: Beatmap{ID: 39}, PP: 216.26}, {Beatmap: Beatmap{ID: 40}, PP: 65.15},
				{Beatmap: Beatmap{ID: 41}, PP: 119.52}, {Beatmap: Beatmap{ID: 42}, PP: 599.61}, {Beatmap: Beatmap{ID: 43}, PP: 433.27}, {Beatmap: Beatmap{ID: 44}, PP: 585.17},
				{Beatmap: Beatmap{ID: 45}, PP: 207.78}, {Beatmap: Beatmap{ID: 46}, PP: 69.52}, {Beatmap: Beatmap{ID: 47}, PP: 399.58}, {Beatmap: Beatmap{ID: 48}, PP: 303.71},
				{Beatmap: Beatmap{ID: 49}, PP: 198.79}, {Beatmap: Beatmap{ID: 50}, PP: 558.97}, {Beatmap: Beatmap{ID: 51}, PP: 26.11}, {Beatmap: Beatmap{ID: 52}, PP: 14.84},
				{Beatmap: Beatmap{ID: 53}, PP: 335.68}, {Beatmap: Beatmap{ID: 54}, PP: 223.33}, {Beatmap: Beatmap{ID: 55}, PP: 24.43}, {Beatmap: Beatmap{ID: 56}, PP: 620.81},
				{Beatmap: Beatmap{ID: 57}, PP: 896.72}, {Beatmap: Beatmap{ID: 58}, PP: 160.86}, {Beatmap: Beatmap{ID: 59}, PP: 26.66}, {Beatmap: Beatmap{ID: 60}, PP: 486.88},
				{Beatmap: Beatmap{ID: 61}, PP: 521.61}, {Beatmap: Beatmap{ID: 62}, PP: 6.96}, {Beatmap: Beatmap{ID: 63}, PP: 64.94}, {Beatmap: Beatmap{ID: 64}, PP: 670.44},
				{Beatmap: Beatmap{ID: 65}, PP: 235.29}, {Beatmap: Beatmap{ID: 66}, PP: 61.38}, {Beatmap: Beatmap{ID: 67}, PP: 156.77}, {Beatmap: Beatmap{ID: 68}, PP: 90.61},
				{Beatmap: Beatmap{ID: 69}, PP: 282.06}, {Beatmap: Beatmap{ID: 70}, PP: 123.00}, {Beatmap: Beatmap{ID: 71}, PP: 194.18}, {Beatmap: Beatmap{ID: 72}, PP: 75.08},
				{Beatmap: Beatmap{ID: 73}, PP: 290.83}, {Beatmap: Beatmap{ID: 74}, PP: 119.96}, {Beatmap: Beatmap{ID: 75}, PP: 354.32}, {Beatmap: Beatmap{ID: 76}, PP: 254.35},
				{Beatmap: Beatmap{ID: 77}, PP: 77.67}, {Beatmap: Beatmap{ID: 78}, PP: 663.65}, {Beatmap: Beatmap{ID: 79}, PP: 643.69}, {Beatmap: Beatmap{ID: 80}, PP: 154.35},
				{Beatmap: Beatmap{ID: 81}, PP: 18.37}, {Beatmap: Beatmap{ID: 82}, PP: 720.36}, {Beatmap: Beatmap{ID: 83}, PP: 1.22}, {Beatmap: Beatmap{ID: 84}, PP: 296.96},
				{Beatmap: Beatmap{ID: 85}, PP: 31.98}, {Beatmap: Beatmap{ID: 86}, PP: 784.22}, {Beatmap: Beatmap{ID: 87}, PP: 413.07}, {Beatmap: Beatmap{ID: 88}, PP: 325.01},
				{Beatmap: Beatmap{ID: 89}, PP: 104.54}, {Beatmap: Beatmap{ID: 90}, PP: 128.57}, {Beatmap: Beatmap{ID: 91}, PP: 268.60}, {Beatmap: Beatmap{ID: 92}, PP: 63.30},
				{Beatmap: Beatmap{ID: 93}, PP: 87.16}, {Beatmap: Beatmap{ID: 94}, PP: 164.85}, {Beatmap: Beatmap{ID: 95}, PP: 126.42}, {Beatmap: Beatmap{ID: 96}, PP: 192.25},
				{Beatmap: Beatmap{ID: 97}, PP: 225.16}, {Beatmap: Beatmap{ID: 98}, PP: 95.03}, {Beatmap: Beatmap{ID: 99}, PP: 82.38}, {Beatmap: Beatmap{ID: 100}, PP: 745.63},
				{Beatmap: Beatmap{ID: 101}, PP: 308.36}, {Beatmap: Beatmap{ID: 102}, PP: 177.12}, {Beatmap: Beatmap{ID: 103}, PP: 7.31}, {Beatmap: Beatmap{ID: 104}, PP: 246.99},
				{Beatmap: Beatmap{ID: 105}, PP: 712.97}, {Beatmap: Beatmap{ID: 106}, PP: 32.05}, {Beatmap: Beatmap{ID: 107}, PP: 17.60}, {Beatmap: Beatmap{ID: 108}, PP: 949.36},
				{Beatmap: Beatmap{ID: 109}, PP: 286.49}, {Beatmap: Beatmap{ID: 110}, PP: 937.03}, {Beatmap: Beatmap{ID: 111}, PP: 381.00}, {Beatmap: Beatmap{ID: 112}, PP: 100.73},
				{Beatmap: Beatmap{ID: 113}, PP: 433.25}, {Beatmap: Beatmap{ID: 114}, PP: 363.07}, {Beatmap: Beatmap{ID: 115}, PP: 628.86}, {Beatmap: Beatmap{ID: 116}, PP: 573.03},
				{Beatmap: Beatmap{ID: 117}, PP: 46.61}, {Beatmap: Beatmap{ID: 118}, PP: 486.21}, {Beatmap: Beatmap{ID: 119}, PP: 23.12}, {Beatmap: Beatmap{ID: 120}, PP: 39.01},
				{Beatmap: Beatmap{ID: 121}, PP: 99.12}, {Beatmap: Beatmap{ID: 122}, PP: 234.70}, {Beatmap: Beatmap{ID: 123}, PP: 112.00}, {Beatmap: Beatmap{ID: 124}, PP: 80.54},
				{Beatmap: Beatmap{ID: 125}, PP: 411.21}, {Beatmap: Beatmap{ID: 126}, PP: 240.20}, {Beatmap: Beatmap{ID: 127}, PP: 317.47}, {Beatmap: Beatmap{ID: 128}, PP: 87.49},
				{Beatmap: Beatmap{ID: 129}, PP: 80.78}, {Beatmap: Beatmap{ID: 130}, PP: 491.77}, {Beatmap: Beatmap{ID: 131}, PP: 722.88}, {Beatmap: Beatmap{ID: 132}, PP: 175.83},
				{Beatmap: Beatmap{ID: 133}, PP: 82.34}, {Beatmap: Beatmap{ID: 134}, PP: 68.62}, {Beatmap: Beatmap{ID: 135}, PP: 265.61}, {Beatmap: Beatmap{ID: 136}, PP: 42.63},
				{Beatmap: Beatmap{ID: 137}, PP: 39.10}, {Beatmap: Beatmap{ID: 138}, PP: 31.04}, {Beatmap: Beatmap{ID: 139}, PP: 258.76}, {Beatmap: Beatmap{ID: 140}, PP: 589.09},
				{Beatmap: Beatmap{ID: 141}, PP: 92.92}, {Beatmap: Beatmap{ID: 142}, PP: 52.16}, {Beatmap: Beatmap{ID: 143}, PP: 66.14}, {Beatmap: Beatmap{ID: 144}, PP: 140.49},
				{Beatmap: Beatmap{ID: 145}, PP: 634.26}, {Beatmap: Beatmap{ID: 146}, PP: 36.78}, {Beatmap: Beatmap{ID: 147}, PP: 166.03}, {Beatmap: Beatmap{ID: 148}, PP: 384.21},
				{Beatmap: Beatmap{ID: 149}, PP: 313.35}, {Beatmap: Beatmap{ID: 150}, PP: 31.25},
			},
			ExpectedCount: 100,
		},
		{
			Name: "Some repeated scores",
			Scores: Scores{
				{Beatmap: Beatmap{ID: 1}, PP: 230.98}, {Beatmap: Beatmap{ID: 1}, PP: 307.03}, {Beatmap: Beatmap{ID: 1}, PP: 400.48},
				{Beatmap: Beatmap{ID: 1}, PP: 21.46}, {Beatmap: Beatmap{ID: 5}, PP: 98.31}, {Beatmap: Beatmap{ID: 1}, PP: 51.08},
			},
			ExpectedCount: 2,
		},
	}
}

func TestRanking_AddScoreOrdered(t *testing.T) {
	for _, testCase := range getTestCases() {
		t.Run(testCase.Name, func(t *testing.T) {
			var rank Ranking
			for _, score := range testCase.Scores {
				rank.AddScore(score)
			}
			assertEqual(t, testCase.ExpectedCount, int(rank.count))
			for i := 1; i < int(rank.count); i++ {
				if rank.scores[i].PP > rank.scores[i-1].PP {
					t.Fatal("Rank is unordered")
				}
			}
		})
	}
}

func TestRanking_AddScore(t *testing.T) {
	var (
		scores = Scores{
			{Beatmap: Beatmap{ID: 1}, PP: 336.242},
			{Beatmap: Beatmap{ID: 2}, PP: 332.834},
			{Beatmap: Beatmap{ID: 3}, PP: 330.403},
			{Beatmap: Beatmap{ID: 4}, PP: 328.735},
			{Beatmap: Beatmap{ID: 5}, PP: 328.239},
		}

		arrExpectedTotalPP = []float64{
			336,
			316,
			298,
			282,
			267,
		}

		expectedTotalPP float64
	)
	log.Println(scores)

	var rank Ranking
	log.Println(&rank)

	for i, score := range scores {
		rank.AddScore(score)
		assertEqual(t, i+1, int(rank.count))

		for j := 1; j < int(rank.count); j++ {
			if rank.scores[j].PP > rank.scores[j-1].PP {
				t.Fatalf("Rank is unordered")
			}
		}

		expectedTotalPP += arrExpectedTotalPP[i]
		assertEqual(t, expectedTotalPP, rank.GetTotalPP())
	}
	log.Println(&rank)
}

func TestRanking_GetTotalPPSingle(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		const ExpectedTotalPP = float64(100)
		var rank = Ranking{
			count: 1,
			scores: [RankSize]Score{
				{PP: ExpectedTotalPP},
			},
		}

		assertEqual(t, ExpectedTotalPP, rank.GetTotalPP())
	})

	t.Run("Multiple", func(t *testing.T) {
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
	})
}

func BenchmarkRanking_AddScore(b *testing.B) {
	var testCases = getTestCases()
	for _, testCase := range testCases {
		b.Run(testCase.Name, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				var rank Ranking
				for _, score := range testCase.Scores {
					rank.AddScore(score)
				}
			}
		})
	}
}

func BenchmarkRanking_GetTotalPP(b *testing.B) {
	// Build ranks from test cases
	var testCases = getTestCases()
	var ranks = make([]Ranking, 0, len(testCases))
	for _, testCase := range testCases {
		var rank Ranking
		for _, score := range testCase.Scores {
			rank.AddScore(score)
		}
		ranks = append(ranks, rank)
	}

	// Benchmark rankings
	for i, rank := range ranks {
		b.Run(testCases[i].Name, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_ = rank.GetTotalPP()
			}
		})
	}
}
