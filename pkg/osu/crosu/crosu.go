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
import (
	"strings"
	"unsafe"
)

// TODO CHECK IF CROSU NEEDS UPDATING
type ModType C.u32
type Usize C.size_t

const (
	NF ModType = 1 << 0
	EZ         = 1 << 1
	TD         = 1 << 2
	HD         = 1 << 3
	HR         = 1 << 4
	DT         = 1 << 5
	RX         = 1 << 6
	HT         = 1 << 7
	FL         = 1 << 8
	SO         = 1 << 9
)

func (m ModType) FromString(s string) ModType {
	switch s {
	case "NF":
		return NF
	case "EZ":
		return EZ
	case "TD":
		return TD
	case "HD":
		return HD
	case "HR":
		return HR
	case "DT":
		return DT
	case "RX":
		return RX
	case "HT":
		return HT
	case "FL":
		return FL
	case "SO":
		return SO
	default:
		return 0
	}
}

func ModTypeFromStringArray(arr []string) (res ModType) {
	res = 0

	for _, s := range arr {
		res |= res.FromString(s)
	}

	return
}

func (m ModType) ToC() C.u32 {
	return (C.u32)(m)
}

type GameMode C.GameMode

const (
	OSU   GameMode = 0
	TAIKO          = 1
	CATCH          = 2
	MANIA          = 3
)

func GameModeFromString(s string) GameMode {
	switch strings.ToLower(s) {
	case "osu":
		return OSU
	case "taiko":
		return TAIKO
	case "catch":
		return CATCH
	case "mania":
		return MANIA
	default:
		return OSU
	}
}

func (m GameMode) ToC() C.GameMode {
	return (C.GameMode)(m)
}

func GetPPFromMap(beatmap []byte, combo, n300, n100, n050, misses int, mods ModType, mode GameMode) float64 {
	if beatmap[len(beatmap)-1] != 0 {
		beatmap = append(beatmap, 0)
	}
	ptr := (*C.char)(unsafe.Pointer(&beatmap[0]))

	var score = C.Score{
		combo:  (C.size_t)(combo),
		n300:   (C.size_t)(n300),
		n100:   (C.size_t)(n100),
		n050:   (C.size_t)(n050),
		misses: (C.size_t)(misses),
	}

	return (float64)(C.GetPPFromMap(ptr, false, mods.ToC(), mode.ToC(), score))
}
