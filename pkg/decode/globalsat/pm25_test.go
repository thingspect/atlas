//go:build !integration

package globalsat

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestPM25(t *testing.T) {
	t.Parallel()

	// PM2.5 payloads, see PM25() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// PM25.
		{"03096113950088", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
			{Attr: "pm25_ugm3", Value: int32(136)},
		}, ""},
		{"03096113950000", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
			{Attr: "pm25_ugm3", Value: int32(0)},
		}, ""},
		{"03096113959999", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
			{Attr: "pm25_ugm3", Value: int32(39321)},
		}, ""},
		// PM25 bad length.
		{"000102030405", nil, "ls11x format bad length: 000102030405"},
		{"0001020304050607", nil, "ls11x format bad length: 0001020304050607"},
		// PM25 bad identifier.
		{"04096113950292", []*decode.Point{
			{Attr: "temp_c", Value: float64(24)},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
		}, "pm25 format bad identifier: 04096113950292"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(test.inp)
			require.NoError(t, err)

			res, err := PM25(bInp)
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
