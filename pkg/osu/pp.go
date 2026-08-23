package osu

import (
	"log/slog"

	"github.com/minoxs/osu-phantom/pkg/osu/crosu"
	"github.com/minoxs/osu-phantom/pkg/osu/optimization"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// attrKey identifies cached difficulty attributes.
// Attributes depend on mods, so the same map with different mods gets separate entries.
type attrKey struct {
	ID   int64
	Mods crosu.ModType
}

// attributeCache stores osu!standard difficulty attributes so pp can be recomputed
// without re-downloading or re-parsing the map file.
var attributeCache = optimization.NewCache[attrKey, crosu.OsuDifficultyAttributes](10000)

func (c *Client) GetPP(score *player.Score) {
	if score.PP > 0 {
		return
	}

	var mode = crosu.GameModeFromString(score.Mode())
	var mods = crosu.ModTypeFromStringArray(score.Mods.Acronyms())

	slog.Info("Calculating score pp with crosu", "BeatmapID", score.Beatmap.ID, "Title", score.BeatmapSet.Title, "Mods", score.Mods.Acronyms(), "Mode", score.Mode())

	if mode == crosu.OSU {
		score.PP = c.getOsuPP(score, mods)
		return
	}

	var beatmap, err = c.DownloadBeatmap(score.Beatmap.ID)
	if err != nil {
		slog.Error("Error downloading beatmap", "BeatmapID", score.Beatmap.ID, "Title", score.BeatmapSet.Title)
		return
	}

	score.PP = crosu.GetPPFromMap(
		beatmap,
		score.MaxCombo,
		score.Statistics.Great,
		score.Statistics.Ok,
		score.Statistics.Meh,
		score.Statistics.Miss,
		mods,
		mode,
	)
}

func (c *Client) getOsuPP(score *player.Score, mods crosu.ModType) float64 {
	var key = attrKey{ID: score.Beatmap.ID, Mods: mods}

	var attr, found = attributeCache.Get(key)
	if !found {
		var beatmap, err = c.DownloadBeatmap(score.Beatmap.ID)
		if err != nil {
			slog.Error("Error downloading beatmap", "BeatmapID", score.Beatmap.ID, "Title", score.BeatmapSet.Title)
			return 0
		}

		var ok bool
		attr, ok = crosu.GetOsuDifficultyAttributes(beatmap, mods)
		if !ok {
			slog.Error("Error calculating difficulty attributes", "BeatmapID", score.Beatmap.ID, "Title", score.BeatmapSet.Title)
			return 0
		}

		attributeCache.Set(key, attr)
	}

	return crosu.GetOsuPP(
		attr,
		mods,
		score.MaxCombo,
		score.Statistics.Great,
		score.Statistics.Ok,
		score.Statistics.Meh,
		score.Statistics.Miss,
	)
}
