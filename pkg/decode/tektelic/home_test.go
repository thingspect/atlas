//go:build !integration

package tektelic

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestHome(t *testing.T) {
	t.Parallel()

	// Home Sensor payloads, see Home() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// Motion.
		{"0a0000", []*decode.Point{{Attr: "motion", Value: false}}, ""},
		{"0a00ff", []*decode.Point{{Attr: "motion", Value: true}}, ""},
		// Temperature.
		{"0367000a", []*decode.Point{
			{Attr: "temp_c", Value: 1.0},
			{Attr: "temp_f", Value: 33.8},
		}, ""},
		{"036700c4", []*decode.Point{
			{Attr: "temp_c", Value: 19.6},
			{Attr: "temp_f", Value: 67.3},
		}, ""},
		// Humidity.
		{"046814", []*decode.Point{{Attr: "humidity_pct", Value: 10.0}}, ""},
		{"04687f", []*decode.Point{{Attr: "humidity_pct", Value: 63.5}}, ""},
		// Battery (V).
		{"00ff012c", []*decode.Point{{Attr: "battery_v", Value: 3.0}}, ""},
		{"00ff0138", []*decode.Point{{Attr: "battery_v", Value: 3.12}}, ""},
		// Temperature, Humidity, and Battery (V).
		{"036700c404687f00ff0138", []*decode.Point{
			{Attr: "temp_c", Value: 19.6},
			{Attr: "temp_f", Value: 67.3},
			{Attr: "humidity_pct", Value: 63.5},
			{Attr: "battery_v", Value: 3.12},
		}, ""},
		{"036700d004688000ff0139", []*decode.Point{
			{Attr: "temp_c", Value: 20.8},
			{Attr: "temp_f", Value: 69.4},
			{Attr: "humidity_pct", Value: 64.0},
			{Attr: "battery_v", Value: 3.13},
		}, ""},
		// Home bad length.
		{"0a", nil, "home format bad length: 0a"},
		// Home bad identifier.
		{"ff0000", nil, "home format bad identifier: ff0000"},
		// Digital bad identifier.
		{"0a0100", nil, "typeDigital format bad identifier: 0a0100"},
		// Home unused trailing bytes.
		{"0a0000ff", []*decode.Point{{Attr: "motion", Value: false}}, "home format unused trailing bytes: ff"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(test.inp)
			require.NoError(t, err)

			res, err := Home(bInp)
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
