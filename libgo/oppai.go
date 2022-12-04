package libgo

/*
#cgo CFLAGS: -I ../inc
#cgo LDFLAGS: -L ../libc -lperf -lstdc++
#include "./library.h"
*/
import "C"
import (
	"fmt"
	"log"
)

func GetPP() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from: ", r)
		}
	}()

	fmt.Println(C.test())
	C.hello()
}
