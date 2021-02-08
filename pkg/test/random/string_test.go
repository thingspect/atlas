// +build !integration

package random

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytes(t *testing.T) {
	t.Parallel()

	for i := 5; i < 15; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			b1 := Bytes(uint(lTest))
			b2 := Bytes(uint(lTest))
			t.Logf("b1, b2: %x, %x", b1, b2)

			require.Len(t, b1, lTest)
			require.Len(t, b2, lTest)
			require.NotEqual(t, b1, b2)
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	for i := 5; i < 15; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			s1 := String(uint(lTest))
			s2 := String(uint(lTest))
			t.Logf("s1, s2: %v, %v", s1, s2)

			require.Len(t, s1, lTest)
			require.Len(t, s2, lTest)
			require.NotEqual(t, s1, s2)
		})
	}
}

func TestEmail(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can generate %v", lTest), func(t *testing.T) {
			t.Parallel()

			e1 := Email()
			e2 := Email()
			t.Logf("e1, e2: %v, %v", e1, e2)

			require.True(t, strings.HasSuffix(e1, "@thingspect.com"))
			require.True(t, strings.HasSuffix(e2, "@thingspect.com"))
			require.NotEqual(t, e1, e2)
		})
	}
}
