package phantom

import (
	"github.com/minoxs/osu-phantom/pkg/osu/optimization"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// DefaultCacheSize bounds how many distinct beatmaps a cached fetcher remembers.
const DefaultCacheSize = 10000

type mapMeta struct {
	beatmap player.Beatmap
	set     player.BeatmapSet
}

// NewCachedFetcher memoizes f by beatmap id, so a map fetched once for any tracked
// user's score is not fetched again. Only successful fetches are cached. A size of
// zero or less disables caching.
func NewCachedFetcher(f BeatmapFetcher, size int) BeatmapFetcher {
	cache := optimization.NewCache[int64, mapMeta](size)
	return func(id int64) (player.Beatmap, player.BeatmapSet, error) {
		if m, ok := cache.Get(id); ok {
			return m.beatmap, m.set, nil
		}
		bm, bs, err := f(id)
		if err != nil {
			return bm, bs, err
		}
		cache.Set(id, mapMeta{beatmap: bm, set: bs})
		return bm, bs, nil
	}
}
