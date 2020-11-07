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

			require.Len(t, s1, lTest, "Should be correct length")
			require.Len(t, s2, lTest, "Should be correct length")
			// Collisions on 1-character strings are common.
			if lTest > 1 {
				require.NotEqual(t, s1, s2, "Should not be equal")
			}
		})
	}
}
