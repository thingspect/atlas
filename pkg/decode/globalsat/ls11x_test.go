//go:build !integration

package globalsat

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestLs11x(t *testing.T) {
	t.Parallel()

	// LS-11X payloads, see ls11x() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// LS-11X.
		{"01096113950292", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
		}, ""},
		{"020a3a10e80000", []*decode.Point{
			{Attr: "temp_c", Value: 26.2},
			{Attr: "temp_f", Value: 79.1},
			{Attr: "humidity_pct", Value: 43.28},
		}, ""},
		{"020a1810e70000", []*decode.Point{
			{Attr: "temp_c", Value: 25.8},
			{Attr: "temp_f", Value: 78.5},
			{Attr: "humidity_pct", Value: 43.27},
		}, ""},
		// LS-11X bad length.
		{"000102030405", nil, "ls11x format bad length: 000102030405"},
		{"0001020304050607", nil, "ls11x format bad length: 0001020304050607"},
		// LS-11X humidity outside allowed range.
		{"020a18f0e70000", []*decode.Point{
			{Attr: "temp_c", Value: 25.8},
			{Attr: "temp_f", Value: 78.5},
		}, "ls11x format humidity outside allowed range: 020a18f0e70000"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(test.inp)
			require.NoError(t, err)

			res, err := ls11x(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, test.res, res)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}
