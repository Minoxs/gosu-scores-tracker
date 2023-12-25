package optimization

import "sync"

//type Cache[Key any, Value any] interface {
//	Get(Key) Value
//	Set(Key, Value)
//}

const (
	Kb = 1024
	Mb = 1024 * Kb
)

type BeatmapCache struct {
	MaxUnitSize uint32
	CacheSize   uint32

	lock  sync.RWMutex
	size  int
	keys  []int
	cache map[int][]byte
}

func (b *BeatmapCache) Init() *BeatmapCache {
	b.size = 0
	b.keys = make([]int, 0)
	b.cache = make(map[int][]byte)

	return b
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
	}

	return true
}

func (b *BeatmapCache) Get(beatmapID int) (beatmap []byte, found bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	beatmap, found = b.cache[beatmapID]
	return
}

func (b *BeatmapCache) Set(beatmapID int, beatmap []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()

	var size = len(beatmap)
	if size > 0 && b.ensureSpace(size) {
		b.keys = append(b.keys, beatmapID)
		b.cache[beatmapID] = beatmap
		b.size += size
	}
}
