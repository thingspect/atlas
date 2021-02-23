// Package random provides functions for generating random values for tests.
// This package must not be used outside of tests.
package random

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// Bytes returns a new random []byte. This function must not be used outside of
// tests.
func Bytes(n uint) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// rand.Read should not fail.
		log.Fatalf("String rand.Read: %v", err)
	}

	return b
}

// String returns a new random hex string n characters in length. This function
// must not be used outside of tests.
func String(n uint) string {
	return hex.EncodeToString(Bytes(n/2 + 1))[:n]
}

// Email generates a random email at thingspect.com. This function must not be
// used outside of tests.
func Email() string {
	return String(10) + "@thingspect.com"
}
