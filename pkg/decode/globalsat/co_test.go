//go:build !integration

package globalsat

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestCO(t *testing.T) {
	t.Parallel()

	// CO payloads, see CO() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// CO.
		{"020a3a10e80000", []*decode.Point{
			{Attr: "temp_c", Value: 26.2},
			{Attr: "temp_f", Value: 79.1},
			{Attr: "humidity_pct", Value: 43.28},
			{Attr: "co_ppm", Value: int32(0)},
		}, ""},
		{"020a1810e70000", []*decode.Point{
			{Attr: "temp_c", Value: 25.8},
			{Attr: "temp_f", Value: 78.5},
			{Attr: "humidity_pct", Value: 43.27},
			{Attr: "co_ppm", Value: int32(0)},
		}, ""},
		{"020a1810e79999", []*decode.Point{
			{Attr: "temp_c", Value: 25.8},
			{Attr: "temp_f", Value: 78.5},
			{Attr: "humidity_pct", Value: 43.27},
			{Attr: "co_ppm", Value: int32(39321)},
		}, ""},
		// CO bad length.
		{"000102030405", nil, "ls11x format bad length: 000102030405"},
		{"0001020304050607", nil, "ls11x format bad length: 0001020304050607"},
		// CO bad identifier.
		{"030a3a10e80000", []*decode.Point{
			{Attr: "temp_c", Value: 26.2},
			{Attr: "temp_f", Value: 79.1},
			{Attr: "humidity_pct", Value: 43.28},
		}, "co format bad identifier: 030a3a10e80000"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(test.inp)
			require.NoError(t, err)

			res, err := CO(bInp)
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
