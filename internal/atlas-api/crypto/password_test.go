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
		t.Run(fmt.Sprintf("Can check %+v", test), func(t *testing.T) {
			t.Parallel()

			err := CheckPass(test.inp)
			t.Logf("err: %v", err)
			require.Equal(t, test.err, err)
		})
	}
}

func TestHashPass(t *testing.T) {
	t.Parallel()

	pass := random.String(10)

	h1, err := HashPass(pass)
	t.Logf("h1, err: %s, %v", h1, err)
	require.NoError(t, err)
	require.Len(t, h1, 60)

	h2, err := HashPass(pass)
	t.Logf("h2, err: %s, %v", h2, err)
	require.NoError(t, err)
	require.Len(t, h2, 60)

	require.NotEqual(t, h1, h2)
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
