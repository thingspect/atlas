// +build !unit

package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	tests := []struct {
		inp string
		err string
	}{
		// Success.
		{testConfig.PostgresURI, ""},
		// Database does not exist.
		{fmt.Sprintf("%s_not_exist", testConfig.PostgresURI), "does not exist"},
		// Wrong port.
		{"postgres://localhost:5433/db", "connect: connection refused"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can connect %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := New(lTest.inp)
			t.Logf("res, err: %+v, %v", res, err)
			if lTest.err == "" {
				require.NotNil(t, res, "Should match result")
			} else {
				require.Contains(t, err.Error(), lTest.err,
					"Should match error")
			}
		})
	}
}
