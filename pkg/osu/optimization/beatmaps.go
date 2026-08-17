package optimization

// beatmaps is the process-wide beatmap file cache, shared across tracked users so
// a map downloaded for one is reused for the next. It is owned here, in the
// package that defines the cache, rather than by callers. Disabled until
// ConfigureBeatmaps sets non-zero limits.
var beatmaps = (&BeatmapCache{}).Init()

// ConfigureBeatmaps sets the shared beatmap cache bounds. Zero disables caching.
func ConfigureBeatmaps(maxUnitSize, cacheSize uint32) {
	beatmaps.SetLimits(maxUnitSize, cacheSize)
}

// GetBeatmap returns a cached beatmap file and whether it was present.
func GetBeatmap(id int64) ([]byte, bool) {
	return beatmaps.Get(id)
}

// PutBeatmap stores a beatmap file in the shared cache.
func PutBeatmap(id int64, data []byte) {
	beatmaps.Set(id, data)
}
