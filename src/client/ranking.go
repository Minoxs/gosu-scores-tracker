package client

import (
	"fmt"
	"math"

	"osu-phantom/src/osu"
)

const RankSize = 100

type Ranking struct {
	count  int8
	scores [RankSize]*osu.Score
}

func (r *Ranking) String() string {
	return fmt.Sprintf("Count=%d : Scores=%v : TotalPP=%.0f", r.count, r.scores[0:r.count], r.GetTotalPP())
}

func (r *Ranking) AddScore(s osu.Score) {
	// Calculate PP
	s.GetPP()

	// Find index to insert at
	var i int8 = 0
	for ; i < r.count; i++ {
		var com = r.scores[i]

		// If better score already exists, just leave
		if s.Beatmap.ID == com.Beatmap.ID && s.PP < com.PP {
			return
		}

		// Better score !
		if s.PP > com.PP {
			break
		}
	}

	// If score is out of the rankings, just leave
	if i >= RankSize {
		return
	}

	// Find index of last score to move
	var j = i
	for ; j < r.count; j++ {
		var com = r.scores[j]
		if s.Beatmap.ID == com.Beatmap.ID {
			// Score will be removed
			r.count--
			break
		}
	}

	// Increse count if ranks aren't full
	if r.count < RankSize {
		r.count++
	}

	// Make sure range is valid
	if j >= RankSize {
		j = RankSize - 1
	}
	// Move scores
	for k := j - 1; k >= i; k-- {
		r.scores[k+1] = r.scores[k]
	}

	// Add score
	r.scores[i] = &s
}

func (r *Ranking) GetTotalPP() (res float64) {
	res = 0
	var i int8 = 0
	for ; i < r.count; i++ {
		res += r.scores[i].PP * math.Pow(0.95, (float64)(i))
	}
	return
}
