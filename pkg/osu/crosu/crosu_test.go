package crosu

import (
	"os"
	"testing"
)

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assert(t *testing.T, msg string, b bool) {
	if !b {
		t.Fatalf(msg)
	}
}

func TestGetPPFromMap(t *testing.T) {
	const testFilePath = "../../../deps/crosu-pp/tests/files/2785319.osu"

	var f, err = os.ReadFile(testFilePath)
	assertNoErr(t, err)

	var pp = GetPPFromMap(f, 909, 909, 0, 0, 0, ModType(0), OSU)
	assert(t, "PP calculation failed", pp > 0)
	assert(t, "PP is too low", pp > 200)
}
