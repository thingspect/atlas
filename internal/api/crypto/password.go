// Package crypto provides cryptography functions.
package crypto

import (
	"strings"

	"github.com/thingspect/atlas/pkg/consterr"
	"golang.org/x/crypto/bcrypt"
)

const ErrWeakPass consterr.Error = "weak password, see NIST password guidelines"

// CheckPass checks whether a password is weak or disallowed.
func CheckPass(pass string) error {
	if len(pass) < 10 ||
		strings.Contains(strings.ToLower(pass), "thingsp") ||
		strings.Contains(weakPasswords, strings.ToLower(pass)) {
		return ErrWeakPass
	}

	return nil
}

// HashPass returns a hash of a password.
func HashPass(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

// CompareHashPass compares a hashed password with its possible plaintext
// equivalent. It returns nil on success or an error on failure.
func CompareHashPass(hash []byte, pass string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(pass))
}
