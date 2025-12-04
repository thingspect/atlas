//go:build !integration

package auth

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	body := make([]byte, 10)
	_, err = rand.Read(body)
	require.NoError(t, err)

	tests := []struct {
		inpKey       []byte
		inpPlaintext []byte
		resLen       int
		err          error
	}{
		{key, body, 38, nil},
		{key, []byte{}, 28, nil},
		{[]byte{}, body, 0, ErrKeyLength},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can encrypt %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := Encrypt(test.inpKey, test.inpPlaintext)
			t.Logf("res, err: %x, %v", res, err)
			require.Len(t, res, test.resLen)
			require.Equal(t, test.err, err)
		})
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	body := make([]byte, 10)
	_, err = rand.Read(body)
	require.NoError(t, err)

	tests := []struct {
		inpEncKey     []byte
		inpDecKey     []byte
		inpPlaintext  []byte
		inpCiphertext []byte
		err           error
	}{
		{key, key, body, nil, nil},
		{key, []byte("notasecurekey"), nil, nil, ErrKeyLength},
		{key, key, nil, []byte{}, ErrMalformed},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can encrypt %+v", test), func(t *testing.T) {
			t.Parallel()

			resEnc, err := Encrypt(test.inpEncKey, test.inpPlaintext)
			t.Logf("resEnc, err: %x, %v", resEnc, err)
			require.NoError(t, err)

			var resDec []byte
			if test.inpCiphertext == nil {
				resDec, err = Decrypt(test.inpDecKey, resEnc)
			} else {
				resDec, err = Decrypt(test.inpDecKey, test.inpCiphertext)
			}
			t.Logf("resDec, err: %x, %#v", resDec, err)
			require.Equal(t, test.inpPlaintext, resDec)
			require.Equal(t, test.err, err)
		})
	}
}
