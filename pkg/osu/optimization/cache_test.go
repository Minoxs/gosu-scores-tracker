package optimization

import (
	"log"
	"testing"
)

func assertEqual[T comparable](t *testing.T, val1 T, val2 T) {
	if val1 != val2 {
		t.Fatalf("Values are not equal\nExpected=%v\nActual=%v\n", val1, val2)
	}
}

func assertStrictEqual(t *testing.T, arr1 []byte, arr2 []byte) {
	if len(arr1) != len(arr2) {
		t.Fatalf("Arrays don't match\nExpected=%v\nActual=%v\n", arr1, arr2)
	}

	for i, _ := range arr1 {
		if arr1[i] != arr2[i] {
			t.Fatalf("Arrays don't match\nExpected=%s\nActual=%s\n", arr1, arr2)
		}
	}
}

func TestGetSet(t *testing.T) {
	var cache = BeatmapCache{
		MaxUnitSize: 3,
		CacheSize:   6,
	}
	cache.Init()

	var (
		test1 = []byte{10, 100, 110}
		test2 = []byte{5, 2}
		test3 = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	)

	var (
		expected []byte
		actual   []byte
		found    bool
	)

	// Add 6 bytes to cache, reaching maximum
	expected = test1
	cache.Set(1, expected)
	log.Printf("%d : CacheSize=%d\n", 1, cache.CurrentSize())
	cache.Set(2, expected)
	log.Printf("%d : CacheSize=%d\n", 2, cache.CurrentSize())

	// This should NOT affect the cache since it's over the limit
	cache.Set(3, test3)
	log.Printf("%d : CacheSize=%d\n", 3, cache.CurrentSize())

	// Look for 1
	actual, found = cache.Get(1)
	assertEqual(t, true, found)
	assertStrictEqual(t, expected, actual)

	// Look for 2
	actual, found = cache.Get(2)
	assertEqual(t, true, found)
	assertStrictEqual(t, expected, actual)

	// 3 should not be in the cache
	_, found = cache.Get(3)
	assertEqual(t, false, found)

	// Adding 2 more bytes should remove ID=1
	expected = test2
	cache.Set(4, expected)
	log.Printf("%d : CacheSize=%d\n", 4, cache.CurrentSize())

	actual, found = cache.Get(4)
	assertEqual(t, true, found)
	assertStrictEqual(t, expected, actual)

	// Should be removed
	_, found = cache.Get(1)
	assertEqual(t, false, found)

	// 2 should not be affected
	actual, found = cache.Get(2)
	assertEqual(t, true, found)
	assertStrictEqual(t, test1, actual)

	// Increase cache size
	cache.MaxUnitSize = 10
	cache.CacheSize = 10

	// This should be the only thing in the cache now
	expected = test3
	cache.Set(10, expected)
	log.Printf("%d : CacheSize=%d\n", 5, cache.CurrentSize())

	actual, found = cache.Get(10)
	assertEqual(t, true, found)
	assertStrictEqual(t, expected, actual)

	for i := 0; i < 10; i++ {
		_, found = cache.Get(i)
		assertEqual(t, false, found)
	}
}
