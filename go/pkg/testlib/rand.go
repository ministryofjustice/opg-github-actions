package testlib

import "math/rand"

func TestRandInRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
