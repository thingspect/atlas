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

	for i := range 5 {
		t.Run(fmt.Sprintf("Can error %v", i), func(t *testing.T) {
			t.Parallel()

			errStr := random.String(10)

			err := Error(errStr)
			t.Logf("err: %v", err)
			// Verify err implements error.
			var _ error = err

			require.Equal(t, errStr, err.Error())
			// Errors should compare exactly without reflection.
			require.True(t, err == Error(errStr)) //nolint:testifylint // Above.

			wrapErr := fmt.Errorf("%w: %s", err, random.String(10))
			t.Logf("wrapErr: %v", wrapErr)
			require.ErrorIs(t, wrapErr, err)
		})
	}
}
