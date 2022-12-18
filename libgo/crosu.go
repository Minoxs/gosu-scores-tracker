package libgo

/*
	This package handles the bridge between Go and C code
*/

/*
#cgo CFLAGS: -I ../deps/crosu-pp/bindings
#cgo LDFLAGS: -L ../deps/crosu-pp/target/release -lcrosu_pp
#include "crosu.h"
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"unsafe"
)

type ModType C.u32
type Usize C.size_t

const (
	NF ModType = 1 << 0
	EZ         = 1 << 1
	TD         = 1 << 2
	HD         = 1 << 3
	HR         = 1 << 3
	DT         = 1 << 4
	RX         = 1 << 5
	HT         = 1 << 6
	FL         = 1 << 7
	SO         = 1 << 8
)

type Score struct {
	Combo  Usize
	N300   Usize
	N100   Usize
	N050   Usize
	Misses Usize
}

func GetPP() {
	buf, _ := ioutil.ReadFile("./test.osu")
	ptr := (*C.char)(unsafe.Pointer(&buf[0]))

	//var tst C.OsuDiffResult = C.GetOsuDifficultyAttributes(ptr, false, C.HR)
	//if tst.success {
	//	fmt.Printf("Star: %f, Aim: %f, Speed: %f\n", tst.attr.stars, tst.attr.aim, tst.attr.speed)
	//}

	var score = C.Score{
		combo: 476,
		n300:  281,
		n100:  37,
		n050:  4,
	}

	fmt.Println("PP: ", C.GetPPFromMap(ptr, false, C.HR, C.Osu, score))
}
