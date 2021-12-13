package osu

import "math"

func weight(rawPP float32, index int) float32 {
	return rawPP * float32(math.Pow(0.95, float64(index)))
}
