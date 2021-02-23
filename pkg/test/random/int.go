package random

import (
	"crypto/rand"
	"log"
	"math/big"
)

// Intn returns, as an int, a non-negative cryptographically secure number in
// [0,n) from rand.Reader. This function must not be used outside of tests.
func Intn(n int) int {
	rn, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		log.Fatalf("Intn rand.Int: %v", err)
	}

	// It is safe to cast to Int because the original bounding value was Int.
	return int(rn.Int64())
}
