//go:build !unit

package dao

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/config"
)

func TestNewPgDB(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	tests := []struct {
		inp string
		err string
	}{
		// Success.
		{testConfig.PgURI, ""},
		// Bad URI.
		{"://", "missing protocol scheme"},
		// Database does not exist.
		{testConfig.PgURI + "_not_exist", "does not exist"},
		// Wrong port.
		{"postgres://127.0.0.1:5433/db", "connect: connection refused"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can connect %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := NewPgDB(test.inp)
			t.Logf("res, err: %+v, %#v", res, err)
			if test.err == "" {
				require.NotNil(t, res)
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
