//go:build !integration

package decode

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCToF(t *testing.T) {
	t.Parallel()

	epsilon := math.Nextafter(1, 2) - 1

	tests := []struct {
		inp float64
		res float64
	}{
		{0, 32},
		{-1.234, 29.8},
		{37.89, 100.2},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can convert %+v", test), func(t *testing.T) {
			t.Parallel()

			res := CToF(test.inp)
			t.Logf("res: %v", res)
			require.InDelta(t, test.res, res, epsilon)
		})
	}
}
