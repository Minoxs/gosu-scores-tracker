package optimization

import "sync"

// Cache is a count-bounded cache with FIFO eviction.
// A maxCount of zero or less disables storage.
type Cache[K comparable, V any] struct {
	maxCount int

	lock  sync.RWMutex
	keys  []K
	cache map[K]V
}

func NewCache[K comparable, V any](maxCount int) *Cache[K, V] {
	return &Cache[K, V]{
		maxCount: maxCount,
		keys:     make([]K, 0),
		cache:    make(map[K]V),
	}
}

func (c *Cache[K, V]) Get(key K) (value V, found bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	value, found = c.cache[key]
	return
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.maxCount <= 0 {
		return
	}

	if _, exists := c.cache[key]; exists {
		c.cache[key] = value
		return
	}

	for len(c.keys) >= c.maxCount {
		var remove = c.keys[0]
		delete(c.cache, remove)
		c.keys = c.keys[1:]
	}

	c.keys = append(c.keys, key)
	c.cache[key] = value
}

func (c *Cache[K, V]) Count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.keys)
}
