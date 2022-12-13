package libgo

/*
	This package handles the switch between Go and C code
*/

/*
#cgo CFLAGS: -I ../inc
#cgo LDFLAGS: -L ../libc -lcrosu_pp
#include "crosu.h"
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"unsafe"
)

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
