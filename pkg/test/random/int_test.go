// +build !integration

package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntn(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			n := Intn(99)
			t.Logf("n: %v", n)
			require.GreaterOrEqual(t, n, 0)
			require.Less(t, n, 99)
		})
	}
}
