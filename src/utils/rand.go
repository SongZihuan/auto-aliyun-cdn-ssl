package utils

import (
	"math/rand"
	"time"
)

var r *rand.Rand = nil

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Rand() *rand.Rand {
	if r == nil {
		panic("nil Rand")
	}

	return r
}
