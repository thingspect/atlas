// +build !integration

package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Parallel()

	for i := 0; i < 10; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			s1 := String(uint(lTest))
			s2 := String(uint(lTest))
			t.Logf("s1, s2: %v, %v", s1, s2)

			require.Len(t, s1, lTest)
			require.Len(t, s2, lTest)
			// Collisions on 1- and 2-character strings are common.
			if lTest > 2 {
				require.NotEqual(t, s1, s2)
			}
		})
	}
}
