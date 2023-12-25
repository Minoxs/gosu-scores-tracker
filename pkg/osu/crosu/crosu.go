package crosu

/*
	This package handles the bridge between Go and C code
*/

/*
#cgo CFLAGS: -I ../../../deps/crosu-pp/bindings
#cgo LDFLAGS: -L ../../../deps/crosu-pp/target/release -lcrosu_pp
#include "crosu.h"
*/
import "C"
import "unsafe"

type (
	GameMode C.GameMode
	ModType  C.u32
)

func (m ModType) ToC() C.u32 {
	return (C.u32)(m)
}

func (m GameMode) ToC() C.GameMode {
	return (C.GameMode)(m)
}

// TODO CHECK IF CROSU NEEDS UPDATING

func GetPPFromMap(beatmap []byte, combo, n300, n100, n050, misses int, mods ModType, mode GameMode) float64 {
	if beatmap[len(beatmap)-1] != 0 {
		beatmap = append(beatmap, 0)
	}

	var ptr = (*C.char)(unsafe.Pointer(&beatmap[0]))
	var score = C.Score{
		combo:  (C.size_t)(combo),
		n300:   (C.size_t)(n300),
		n100:   (C.size_t)(n100),
		n050:   (C.size_t)(n050),
		misses: (C.size_t)(misses),
	}

	return (float64)(C.GetPPFromMap(ptr, false, mods.ToC(), mode.ToC(), score))
}
