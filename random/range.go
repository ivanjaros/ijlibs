package random

import (
	"math/rand"
)

func Integer() int {
	return rand.Int()
}

func IntegerRange(min, max int) int {
	return rand.Intn(max-min) + min
}
