// Package crypto provides cryptography functions.
package crypto

import "golang.org/x/crypto/bcrypt"

// HashPass returns a hash of a password.
func HashPass(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

// CompareHashPass compares a hashed password with its possible plaintext
// equivalent. It returns nil on success or an error on failure.
func CompareHashPass(hash []byte, pass string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(pass))
}
