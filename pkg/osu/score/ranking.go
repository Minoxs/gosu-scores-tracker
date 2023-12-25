package score

import (
	"fmt"
	"math"
)

const RankSize = 100

// TODO MAKE SOME UNIT TESTS FOR RANKING (I THINK THEY'RE WRONG)

type Ranking struct {
	count  int8
	scores [RankSize]Score
}

func (r *Ranking) String() string {
	return fmt.Sprintf("Count=%d : TotalPP=%.0f : Scores=%v", r.count, r.GetTotalPP(), Scores(r.scores[:r.count]))
}

func (r *Ranking) AddScore(s Score) {
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

	// Increase count if ranks aren't full
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
	r.scores[i] = s
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
