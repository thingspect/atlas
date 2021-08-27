//go:build !integration

package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPass(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inp string
		err error
	}{
		{random.String(20), nil},
		{random.String(5), ErrWeakPass},
		{"Thingsp3ct", ErrWeakPass},
		{"1234567890", ErrWeakPass},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can check %+v", lTest), func(t *testing.T) {
			t.Parallel()

			err := CheckPass(lTest.inp)
			t.Logf("err: %v", err)
			require.Equal(t, lTest.err, err)
		})
	}
}

func TestHashPass(t *testing.T) {
	t.Parallel()

	pass := random.String(10)
	hashChan := make(chan []byte)

	for i := 0; i < 2; i++ {
		go func() {
			h, err := HashPass(pass)
			t.Logf("h, err: %s, %v", h, err)
			require.NoError(t, err)
			require.Len(t, h, 60)

			hashChan <- h
		}()
	}

	require.NotEqual(t, <-hashChan, <-hashChan)
}

func TestCompareHashPass(t *testing.T) {
	t.Parallel()

	pass := random.String(10)

	hash, err := HashPass(pass)
	t.Logf("hash, err: %s, %v", hash, err)
	require.NoError(t, err)

	t.Run("Can compare correct pass", func(t *testing.T) {
		t.Parallel()

		require.NoError(t, CompareHashPass(hash, pass))
	})

	t.Run("Can compare incorrect pass", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, bcrypt.ErrMismatchedHashAndPassword,
			CompareHashPass(hash, random.String(10)))
	})

	require.Equal(t, bcrypt.ErrHashTooShort, CompareHashPass([]byte{}, pass))
	require.Equal(t, bcrypt.ErrHashTooShort, CompareHashPass(nil, pass))
}
