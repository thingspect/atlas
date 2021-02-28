// +build !integration

package crypto

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
		lTest := test

		t.Run(fmt.Sprintf("Can encrypt %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := Encrypt(lTest.inpKey, lTest.inpPlaintext)
			t.Logf("res, err: %x, %v", res, err)
			require.Len(t, res, lTest.resLen)
			require.Equal(t, lTest.err, err)
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
		lTest := test

		t.Run(fmt.Sprintf("Can encrypt %+v", lTest), func(t *testing.T) {
			t.Parallel()

			resEnc, err := Encrypt(lTest.inpEncKey, lTest.inpPlaintext)
			t.Logf("resEnc, err: %x, %v", resEnc, err)
			require.NoError(t, err)

			var resDec []byte
			if lTest.inpCiphertext == nil {
				resDec, err = Decrypt(lTest.inpDecKey, resEnc)
			} else {
				resDec, err = Decrypt(lTest.inpDecKey, lTest.inpCiphertext)
			}
			t.Logf("resDec, err: %x, %#v", resDec, err)
			require.Equal(t, lTest.inpPlaintext, resDec)
			require.Equal(t, lTest.err, err)
		})
	}
}
