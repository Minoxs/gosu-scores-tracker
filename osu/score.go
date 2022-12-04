package osu

import (
	"fmt"
	"math"
)

type statistics struct {
	Count050  int `json:"count_50"`
	Count100  int `json:"count_100"`
	Count300  int `json:"count_300"`
	CountGeki int `json:"count_geki"`
	CountKatu int `json:"count_katu"`
	CountMiss int `json:"count_miss"`
}

type beatmap struct {
	ID            int     `json:"id"`
	StarRating    float32 `json:"difficulty_rating"`
	Status        string  `json:"status"`
	TotalLength   int     `json:"total_length"`
	OD            int     `json:"od"`
	Ar            float32 `json:"ar"`
	BPM           int     `json:"bpm"`
	CountCircles  int     `json:"count_circles"`
	CountSliders  int     `json:"count_sliders"`
	CountSpinners int     `json:"count_spinners"`
	CS            float32 `json:"cs"`
	HP            float32 `json:"drain"`
	HitLength     int     `json:"hit_length"`
	URL           string  `json:"url"`
}

func (b beatmap) maxCombo() int {
	return b.CountCircles + b.CountSliders + b.CountSpinners
}

type Score struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	Accuracy   float32    `json:"accuracy"`
	Mods       []string   `json:"mods"`
	Score      int        `json:"score"`
	MaxCombo   int        `json:"max_combo"`
	Passed     bool       `json:"passed"`
	Statistics statistics `json:"statistics"`
	PP         float32    `json:"pp"`
	Mode       string     `json:"mode"`
	Beatmap    beatmap    `json:"beatmap"`
}

func (s *Score) GetPP() float32 {
	if s.PP == 0 {
		s.PP = s.computeTotalValue()
	}
	return s.PP
}

func (s *Score) ContainsMod(mod BeatmapMod) bool {
	return containsMod(s.Mods, mod)
}

func (s *Score) GetAR() float64 {
	// TODO ADD MOD MODIFIERS
	return float64(s.Beatmap.Ar)
}

func (s *Score) GetCS() float64 {
	// TODO ADD MOD MODIFIERS
	return float64(s.Beatmap.CS)
}

func (s *Score) GetOD() float64 {
	// TODO ADD MOD MODIFIERS
	return float64(s.Beatmap.OD)
}

// "Translating" from https://github.com/ppy/osu-performance/blob/master/src/performance/osu/OsuScore.cpp

func (s *Score) totalHits() int {
	return s.Statistics.Count050 + s.Statistics.Count100 + s.Statistics.Count300 + s.Statistics.CountMiss
}

func (s *Score) successfulHits() int {
	return s.Statistics.Count050 + s.Statistics.Count100 + s.Statistics.Count300
}

func (s *Score) computeEffectiveMissCount() int {
	maxCombo := float64(s.MaxCombo)
	comboBasedMissCount := float64(0)
	if s.Beatmap.CountSliders > 0 {
		fullComboThreshold := float64(s.Beatmap.maxCombo()) - 0.1*float64(s.Beatmap.CountSliders)
		if maxCombo < fullComboThreshold {
			comboBasedMissCount = fullComboThreshold / math.Max(1.0, maxCombo)
		}
	}

	comboBasedMissCount = math.Min(comboBasedMissCount, float64(s.totalHits()))
	return maxInt(s.Statistics.CountMiss, int(math.Floor(comboBasedMissCount)))
}

// TODO COMPUTE ACC VALUE
// TODO COMPUTE SPEEDE VALUE
// TODO COMPUTE FLESHLIGHT

func (s *Score) computeAimValue() float64 {
	rawAim := s.GetCS() // TODO make sure AIM = CS

	aimValue := math.Pow(5.0*math.Max(1.0, rawAim/0.0675)-4.0, 3.0) / 100000.0

	// Longer maps are worth more
	numTotalHits := float64(s.totalHits())
	lengthBonus := 0.95 + 0.4*math.Min(1.0, numTotalHits/2000.0)
	if numTotalHits > 2000 {
		lengthBonus += math.Log10(numTotalHits/2000.0) * 0.5
	}
	aimValue *= lengthBonus

	// Penalize misses
	effectiveMissCount := s.computeEffectiveMissCount()
	if effectiveMissCount > 0 {
		aimValue *= 0.97 * math.Pow(1.0-math.Pow(float64(effectiveMissCount)/float64(numTotalHits), 0.775), float64(effectiveMissCount))
	}

	// Combo scaling
	maxCombo := float64(s.Beatmap.maxCombo())
	if maxCombo > 0 {
		aimValue *= math.Min(math.Pow(float64(s.MaxCombo), 0.8)/math.Pow(maxCombo, 0.8), 1.0)
	}

	// AR
	approachRate := s.GetAR()
	approachRateFactor := 0.0
	if approachRate > 10.33 {
		approachRateFactor = 0.3 * (approachRate - 10.33)
	} else if approachRate < 8.0 {
		approachRateFactor = 0.1 * (8.0 - approachRate)
	}
	aimValue *= 1.0 + approachRateFactor*lengthBonus

	if s.ContainsMod(HD) {
		aimValue *= 1.0 + 0.04*(12.0-approachRate)
	}

	// Assume 15% of sliders are difficult (what???)
	if s.Beatmap.CountSliders > 0 {
		// TODO FIGURE OUT WHAT THE FUCK IS HAPPENING
		//estimateDifficultSliders := float64(s.Beatmap.CountSliders) * 0.15
		//estimateSliderEndsDropped := math.Min(math.Max(math.Min(float64(s.Statistics.Count100+s.Statistics.Count050+s.Statistics.CountMiss), maxCombo-float64(s.MaxCombo)), 0.0), estimateDifficultSliders)
		//sliderFactor :=
	}

	aimValue *= float64(s.Accuracy)*0.98 + math.Pow(s.GetOD(), 2)/2500.0
	return aimValue
}

func (s *Score) computeTotalValue() float32 {
	// Return 0 if there is an unranked mod
	for _, mod := range s.Mods {
		if !BeatmapMod(mod).Ranked() {
			return 0
		}
	}

	multiplier := float64(1.12)
	if s.ContainsMod(NF) {
		multiplier *= math.Max(0.9, 1.0-0.02*float64(s.computeEffectiveMissCount()))
	}
	if s.ContainsMod(SO) {
		multiplier *= 1.0 - math.Pow(float64(s.Beatmap.CountSpinners)/float64(s.totalHits()), 0.85)
	}

	score := math.Pow(s.computeAimValue(), 1.1)
	return float32(math.Pow(score, 1.0/1.1)) * float32(multiplier)
}

type Scores []Score

func (s Scores) String() (res string) {
	res = ""
	for _, score := range s {
		res += fmt.Sprintf(
			"ID=%d\n"+
				"Accuracy=%f\n"+
				"Mods=%v\n"+
				"Score=%d\n"+
				"PP=%.0f\n"+
				"300=%d\n"+
				"100=%d\n"+
				"050=%d\n"+
				"MIS=%d\n"+
				"COM=%d\n"+
				"URL=%s\n"+
				"-\n",
			score.ID,
			score.Accuracy,
			score.Mods,
			score.Score,
			score.PP,
			score.Statistics.Count300,
			score.Statistics.Count100,
			score.Statistics.Count050,
			score.Statistics.CountMiss,
			score.MaxCombo,
			score.Beatmap.URL,
		)
	}

	if len(res) > 0 {
		return "[\n" + res[:len(res)-2] + "]"
	} else {
		return "[]"
	}
}
