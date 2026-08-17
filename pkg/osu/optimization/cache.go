package optimization

import (
	"log/slog"
	"sync"
)

//type Cache[Key any, Value any] interface {
//	Get(Key) Value
//	Set(Key, Value)
//}

type BeatmapCache struct {
	MaxUnitSize uint32
	CacheSize   uint32

	lock  sync.RWMutex
	size  int
	keys  []int64
	cache map[int64][]byte
}

func (b *BeatmapCache) Init() *BeatmapCache {
	b.size = 0
	b.keys = make([]int64, 0)
	b.cache = make(map[int64][]byte)

	return b
}

// SetLimits updates the cache bounds. Zero for either disables caching.
func (b *BeatmapCache) SetLimits(maxUnitSize, cacheSize uint32) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.MaxUnitSize = maxUnitSize
	b.CacheSize = cacheSize
}

func (b *BeatmapCache) ensureSpace(size int) bool {
	if size > int(b.MaxUnitSize) {
		return false
	}

	for b.size+size > int(b.CacheSize) {
		var remove = b.keys[0]
		b.size -= len(b.cache[remove])
		delete(b.cache, remove)
		b.keys = b.keys[1:]
		slog.Info("Removed from cache", "Beatmap.ID", remove, "Cache.Size", b.size, "Cache.Count", len(b.keys))
	}

	return true
}

func (b *BeatmapCache) Get(beatmapID int64) (beatmap []byte, found bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	beatmap, found = b.cache[beatmapID]
	slog.Debug("Get from cache", "Beatmap.ID", beatmapID, "Found", found)
	return
}

func (b *BeatmapCache) Set(beatmapID int64, beatmap []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()

	var size = len(beatmap)
	if size > 0 && b.ensureSpace(size) {
		b.keys = append(b.keys, beatmapID)
		b.cache[beatmapID] = beatmap
		b.size += size
		slog.Info("Added to cache", "Beatmap.ID", beatmapID, "Beatmap.Size", len(beatmap), "Cache.Size", b.size, "Cache.Count", len(b.keys))
	}
}

func (b *BeatmapCache) CurrentSize() int {
	return b.size
}
