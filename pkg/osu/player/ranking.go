package player

import (
	"fmt"
	"math"
)

const RankSize = 100

type Ranking struct {
	count  int8
	scores [RankSize]Score
}

func (r *Ranking) String() string {
	return fmt.Sprintf("Count=%d : TotalPP=%.0f : Scores=%v", r.count, r.GetTotalPP(), Scores(r.scores[:r.count]))
}

func (r *Ranking) findPosition(s Score) (valid bool, pos int8) {
	for pos = 0; pos < r.count; pos++ {
		var com = r.scores[pos]

		// Return early if the score was already added or
		// A better score on the same map already exists
		if s.ID == com.ID || s.Beatmap.ID == com.Beatmap.ID && s.PP < com.PP {
			return
		}

		// Better score !
		if s.PP > com.PP {
			break
		}
	}
	valid = pos < RankSize
	return
}

func (r *Ranking) insertScore(pos int8, s Score) {
	// Find index of last score to move
	var j = pos
	for ; j < r.count; j++ {
		var com = r.scores[j]
		if s.Beatmap.ID == com.Beatmap.ID {
			// Score will be removed
			r.count--
			break
		}
	}

	// Increase count if ranks aren't full
	if r.count < RankSize {
		r.count++
	}

	// Make sure range is valid
	if j >= RankSize {
		j = RankSize - 1
	}
	// Move scores
	for k := j - 1; k >= pos; k-- {
		r.scores[k+1] = r.scores[k]
	}

	// Add score
	r.scores[pos] = s
}

func (r *Ranking) AddScore(s Score) {
	var valid, pos = r.findPosition(s)
	if valid {
		r.insertScore(pos, s)
	}
}

func (r *Ranking) GetTotalPP() (res float64) {
	res = 0
	var i int8 = 0
	for ; i < r.count; i++ {
		res += r.scores[i].PP * math.Pow(0.95, (float64)(i))
	}
	res = math.Floor(res)
	return
}
