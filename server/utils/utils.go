package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomDuration(min, max int) time.Duration {
	if min > max || min < 1 {
		panic("Invalid range of time")
	}
	rand.Seed(time.Now().UnixNano())

	randomSec := rand.Intn(max-min+1) + min

	return time.Duration(randomSec) * time.Second
}
