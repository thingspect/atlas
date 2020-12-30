// +build !integration,!race

package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"golang.org/x/crypto/bcrypt"
)

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

	require.NoError(t, CompareHashPass(hash, pass))
	require.Equal(t, bcrypt.ErrMismatchedHashAndPassword,
		CompareHashPass(hash, random.String(10)))
}
