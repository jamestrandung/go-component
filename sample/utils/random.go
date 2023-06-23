package utils

import (
	"math/rand"
	"time"
)

var s = rand.NewSource(time.Now().UnixNano())
var r = rand.New(s)

// FlipCoin returns true/false with 50/50 probability.
func FlipCoin() bool {
	return r.Intn(100) < 50
}
