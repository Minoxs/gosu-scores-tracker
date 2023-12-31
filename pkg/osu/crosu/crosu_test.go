package crosu

import (
	"os"
	"testing"
)

type fatal interface {
	Fatal(args ...any)
}

func assertNoErr(t fatal, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assert(t fatal, msg string, b bool) {
	if !b {
		t.Fatal(msg)
	}
}

const testFilePath = "../../../deps/crosu-pp/tests/files/2785319.osu"

func TestGetPPFromMap(t *testing.T) {
	var f, err = os.ReadFile(testFilePath)
	assertNoErr(t, err)

	var pp = GetPPFromMap(f, 909, 909, 0, 0, 0, ModType(0), OSU)
	assert(t, "PP calculation failed", pp > 0)
	assert(t, "PP is too low", pp > 200)
}

func TestGetPPFromFile(t *testing.T) {
	var pp = GetPPFromFile(testFilePath, 909, 909, 0, 0, 0, ModType(0), OSU)
	assert(t, "PP calculation failed", pp > 0)
	assert(t, "PP is too low", pp > 200)
}

func BenchmarkGetPPFromDLL(b *testing.B) {
	var f, err = os.ReadFile(testFilePath)
	assertNoErr(b, err)

	b.Run("PP from Map", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = GetPPFromMap(f, 909, 909, 0, 0, 0, ModType(0), OSU)
		}
	})

	b.Run("PP from File", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = GetPPFromFile(testFilePath, 909, 909, 0, 0, 0, ModType(0), OSU)
		}
	})

	b.Run("PP from Map in Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = GetPPFromMap(f, 909, 909, 0, 0, 0, ModType(0), OSU)
			}
		})
	})

	b.Run("PP from File in Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = GetPPFromFile(testFilePath, 909, 909, 0, 0, 0, ModType(0), OSU)
			}
		})
	})
}
