package utils

import (
	"log"
)

func PanicHandler(funcName string) {
	if err := recover(); err != nil {
		log.Printf("%s : recovered from panic : %v\n", funcName, err)
	}
}
