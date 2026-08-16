package optimization

import "testing"

func TestCacheGetSet(t *testing.T) {
	var cache = NewCache[int, string](3)

	cache.Set(1, "a")
	cache.Set(2, "b")

	var value, found = cache.Get(1)
	assertEqual(t, true, found)
	assertEqual(t, "a", value)

	value, found = cache.Get(2)
	assertEqual(t, true, found)
	assertEqual(t, "b", value)

	_, found = cache.Get(3)
	assertEqual(t, false, found)
}

func TestCacheEviction(t *testing.T) {
	var cache = NewCache[int, int](2)

	cache.Set(1, 10)
	cache.Set(2, 20)
	cache.Set(3, 30)

	assertEqual(t, 2, cache.Count())

	// Oldest key is evicted first
	var _, found = cache.Get(1)
	assertEqual(t, false, found)

	var value bool
	_, value = cache.Get(2)
	assertEqual(t, true, value)
	_, value = cache.Get(3)
	assertEqual(t, true, value)
}

func TestCacheOverwriteKeepsCount(t *testing.T) {
	var cache = NewCache[int, int](2)

	cache.Set(1, 10)
	cache.Set(1, 11)

	assertEqual(t, 1, cache.Count())

	var value, found = cache.Get(1)
	assertEqual(t, true, found)
	assertEqual(t, 11, value)
}

func TestCacheDisabled(t *testing.T) {
	var cache = NewCache[int, int](0)

	cache.Set(1, 10)

	var _, found = cache.Get(1)
	assertEqual(t, false, found)
	assertEqual(t, 0, cache.Count())
}
