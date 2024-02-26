//go:build !integration

package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntn(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			n := Intn(99)
			t.Logf("n: %v", n)
			require.GreaterOrEqual(t, n, 0)
			require.Less(t, n, 99)
		})
	}
}
