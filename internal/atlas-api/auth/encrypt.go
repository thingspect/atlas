package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/thingspect/atlas/pkg/consterr"
)

// Errors returned due to encryption or decryption failures.
const (
	ErrKeyLength consterr.Error = "auth: incorrect key length"
	ErrMalformed consterr.Error = "auth: malformed ciphertext"
)

// Encrypt encrypts data using 256-bit AES-GCM, providing authenticated
// encryption with associated data.
//
// https://github.com/gtank/cryptopasta
//
// https://golang.org/pkg/crypto/cipher/#NewGCM
func Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrKeyLength
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM, providing authenticated
// encryption with associated data.
func Decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrKeyLength
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, ErrMalformed
	}

	return gcm.Open(nil, ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():], nil)
}
