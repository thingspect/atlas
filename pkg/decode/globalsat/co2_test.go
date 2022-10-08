//go:build !integration

package globalsat

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestCO2(t *testing.T) {
	t.Parallel()

	// CO2 payloads, see CO2() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// CO2.
		{"01096113950292", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
			{Attr: "co2_ppm", Value: int32(658)},
		}, ""},
		{"01096113950000", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
			{Attr: "co2_ppm", Value: int32(0)},
		}, ""},
		{"01096113959999", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
			{Attr: "co2_ppm", Value: int32(39321)},
		}, ""},
		// CO2 bad length.
		{"000102030405", nil, "ls11x format bad length: 000102030405"},
		{"0001020304050607", nil, "ls11x format bad length: 0001020304050607"},
		// CO2 bad identifier.
		{"02096113950292", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
		}, "co2 format bad identifier: 02096113950292"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, err := CO2(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
