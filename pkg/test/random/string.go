// Package random provides functions for generating random values for tests.
package random

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// String returns a new random hex string n characters in length. Care should be
// taken when used outside of tests.
func String(n uint) string {
	b := make([]byte, n/2+1)
	if _, err := rand.Read(b); err != nil {
		// rand.Read should not fail.
		log.Fatalf("String rand.Read: %v", err)
	}

	return hex.EncodeToString(b)[:n]
}

// Email generates a random email at thingspect.com.
func Email() string {
	return String(10) + "@thingspect.com"
}
