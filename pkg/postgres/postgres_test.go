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
		{testConfig.PgURI, ""},
		// Database does not exist.
		{fmt.Sprintf("%s_not_exist", testConfig.PgURI), "does not exist"},
		// Wrong port.
		{"postgres://127.0.0.1:5433/db", "connect: connection refused"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can connect %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := New(lTest.inp)
			t.Logf("res, err: %+v, %#v", res, err)
			if lTest.err == "" {
				require.NotNil(t, res)
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
