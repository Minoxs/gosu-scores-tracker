package utils

import (
	"os"
	"strconv"
)

type StringObject string

func ParseInt[T ConInteger](s string, d T) (T, error) {
	c, err := strconv.Atoi(s)
	if err != nil {
		return d, err
	}
	return T(c), nil
}

func GetEnv(key string) StringObject {
	return StringObject(os.Getenv(key))
}

func (so StringObject) String() string {
	return string(so)
}

func (so StringObject) Integer(def int) (res int) {
	//strconv.Atoi()
	res, _ = ParseInt(so.String(), def)
	return res
}
