package util

import (
	"math/rand"
	"time"
)

func RandomFloat(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

func RandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0
}
