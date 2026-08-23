package phantom

import (
	"github.com/minoxs/osu-phantom/pkg/osu"
	"github.com/minoxs/osu-phantom/pkg/osu/optimization"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// DefaultCacheSize bounds how many distinct beatmaps an osu BeatmapProvider
// remembers.
const DefaultCacheSize = 10000

type mapMeta struct {
	beatmap player.Beatmap
	set     player.BeatmapSet
}

// osuBeatmaps is the production BeatmapProvider: it fetches beatmaps through
// osu-phantom, memoized by id so a map fetched for one tracked user's score is not
// fetched again. Its client reserves at the priority the caller chose, so tracked
// enrichment is not starved by lower-priority traffic.
type osuBeatmaps struct {
	provider AuthProvider
	cache    *optimization.Cache[int64, mapMeta]
	client   *osu.Client
}

// NewOsuBeatmapProvider builds a BeatmapProvider over osu-phantom, caching up to
// cacheSize distinct beatmaps and reserving its fetches at priority. A cacheSize of
// zero or less uses DefaultCacheSize.
func NewOsuBeatmapProvider(provider AuthProvider, cacheSize int, priority osu.Priority) BeatmapProvider {
	if cacheSize <= 0 {
		cacheSize = DefaultCacheSize
	}
	return &osuBeatmaps{
		provider: provider,
		cache:    optimization.NewCache[int64, mapMeta](cacheSize),
		client:   osu.NewClient(priority),
	}
}

func (b *osuBeatmaps) Beatmap(id int64) (player.Beatmap, player.BeatmapSet, error) {
	if m, ok := b.cache.Get(id); ok {
		return m.beatmap, m.set, nil
	}
	bm, bs, err := b.client.GetBeatmap(b.provider.GetToken(), id)
	if err != nil {
		return bm, bs, err
	}
	b.cache.Set(id, mapMeta{beatmap: bm, set: bs})
	return bm, bs, nil
}

func (b *osuBeatmaps) ResolvePP(score *player.Score) {
	b.client.GetPP(score)
}
