package random

import (
	"math/rand"
	"time"
)

func GetRandomNumber(b int) int {
	// Handle panic in case b is 0 or negative
	if b <= 0 {
		return 0
	}

	return rand.New(rand.NewSource(time.Now().Unix())).Intn(b)
}
