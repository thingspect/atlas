//go:build !integration

package consterr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestOrg(t *testing.T) {
	t.Parallel()

	for i := 0; i < 5; i++ {
		lTest := i

		t.Run(fmt.Sprintf("Can error %v", lTest), func(t *testing.T) {
			t.Parallel()

			errStr := random.String(10)

			err := Error(errStr)
			t.Logf("err: %v", err)
			// Verify err implements error.
			var _ error = err

			require.Equal(t, errStr, err.Error())
			// Errors should compare exactly without reflection.
			require.True(t, err == Error(errStr))

			wrapErr := fmt.Errorf("%w: %s", err, random.String(10))
			t.Logf("wrapErr: %v", wrapErr)
			require.ErrorIs(t, wrapErr, err)
		})
	}
}
