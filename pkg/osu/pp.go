package osu

import (
	"github.com/minoxs/osu-phantom/pkg/osu/crosu"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
	"log"
)

func GetPP(score *player.Score) {
	if score.PP > 0 {
		return
	}

	var beatmap, err = DownloadBeatmap(score.Beatmap.ID)
	if err != nil {
		log.Println(err)
		return
	}

	score.PP = crosu.GetPPFromMap(
		beatmap,
		score.MaxCombo,
		score.Statistics.Count300,
		score.Statistics.Count100,
		score.Statistics.Count050,
		score.Statistics.CountMiss,
		crosu.ModTypeFromStringArray(score.Mods),
		crosu.GameModeFromString(score.Mode),
	)
}
