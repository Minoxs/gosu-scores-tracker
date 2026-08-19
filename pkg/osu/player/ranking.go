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

func (r *Ranking) findPosition(s *Score) (valid bool, pos int8) {
	for pos = 0; pos < r.count; pos++ {
		// Return early if the score was already added or
		// A better score on the same map already exists
		if s.ID == r.scores[pos].ID || s.Beatmap.ID == r.scores[pos].Beatmap.ID && s.PP <= r.scores[pos].PP {
			return
		}

		// Better score !
		if s.PP > r.scores[pos].PP {
			break
		}
	}
	valid = pos < RankSize
	return
}

func (r *Ranking) insertScore(pos int8, s *Score) {
	// Find index of last score to move
	var j = pos
	for ; j < r.count; j++ {
		if s.Beatmap.ID == r.scores[j].Beatmap.ID {
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
	copy(r.scores[pos+1:j+1], r.scores[pos:j])

	// Add score
	r.scores[pos] = *s
}

func (r *Ranking) AddScore(s Score) (rank int, added bool) {
	var valid, idx = r.findPosition(&s)
	if valid {
		r.insertScore(idx, &s)
	}
	return int(idx + 1), valid
}

func (r *Ranking) GetTotalPP() (res float64) {
	res = 0
	for i := int8(0); i < r.count; i++ {
		res += r.scores[i].PP * math.Pow(0.95, (float64)(i))
	}
	res = math.Floor(res)
	return
}

func (r Ranking) Scores() Scores {
	return r.scores[:r.count]
}

func (r Ranking) Count() int {
	return int(r.count)
}
