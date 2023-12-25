package osu

import (
	"fmt"
	"log"
	"osu-phantom/pkg/osu/calculator"
	"time"
)

type statistics struct {
	Count050  int `json:"count_50"`
	Count100  int `json:"count_100"`
	Count300  int `json:"count_300"`
	CountGeki int `json:"count_geki"`
	CountKatu int `json:"count_katu"`
	CountMiss int `json:"count_miss"`
}

type Beatmap struct {
	ID            int     `json:"id"`
	Version       string  `json:"version"`
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
}

type BeatmapSet struct {
	Title string `json:"title"`
}

type Score struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	CreatedAt  time.Time  `json:"created_at"`
	Accuracy   float32    `json:"accuracy"`
	Mods       []string   `json:"mods"`
	Score      int        `json:"score"`
	MaxCombo   int        `json:"max_combo"`
	Passed     bool       `json:"passed"`
	Statistics statistics `json:"statistics"`
	PP         float64    `json:"pp"`
	Mode       string     `json:"mode"`
	Beatmap    Beatmap    `json:"beatmap"`
	BeatmapSet BeatmapSet `json:"beatmapSet"`
}

func (s *statistics) String() string {
	return fmt.Sprintf("300=%d : 100=%d : 50=%d : Miss=%d", s.Count300, s.Count100, s.Count050, s.CountMiss)
}

type Scores []Score

func (s Scores) String() (res string) {
	res = ""
	for _, score := range s {
		res += score.String() + ",\n"
	}

	if len(res) > 0 {
		return "[\n" + res[:len(res)-2] + "\n]"
	} else {
		return "[]"
	}
}

func (s *Score) String() string {
	return fmt.Sprintf(
		"{ ID=%d : BeatmapID=%d : Title=%s : Diff=%s : Mode=%s : Mods=%v : Score=%d : PP=%.0f : MaxCombo=%d : Acc=%.2f : %s }",
		s.ID,
		s.Beatmap.ID,
		s.BeatmapSet.Title,
		s.Beatmap.Version,
		s.Mode,
		s.Mods,
		s.Score,
		s.PP,
		s.MaxCombo,
		s.Accuracy,
		s.Statistics.String(),
	)
}

// TODO CHECK IF CROSU NEEDS UPDATING

// GetPP returns PP value of score
// If PP from play is 0, will download Beatmap and calculate manually
func (s *Score) GetPP() float64 {
	if s.PP == 0 {
		// Download Beatmap
		var beatmap, err = DownloadBeatmap(s.Beatmap.ID)
		if err != nil {
			log.Println(err)
			return 0
		}

		// Calculate pp
		s.PP = calculator.GetPPFromMap(
			beatmap,
			s.MaxCombo,
			s.Statistics.Count300,
			s.Statistics.Count100,
			s.Statistics.Count050,
			s.Statistics.CountMiss,
			calculator.ModType(0).FromStringArray(s.Mods),
			calculator.GameMode(0).FromString(s.Mode),
		)

		log.Printf("CalculatedPP=%.4f\n", s.PP)
	}

	return s.PP
}
