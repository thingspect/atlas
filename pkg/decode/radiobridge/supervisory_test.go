//go:build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestSupervisory(t *testing.T) {
	t.Parallel()

	// Supervisory payloads, see supervisory() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// Supervisory.
		{"1401080131", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(4)},
			{Attr: "tamper", Value: false},
			{Attr: "battery_volt", Value: 3.1},
		}, ""},
		{"1401170127", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(4)},
			{Attr: "error", Value: "radio_comm"},
			{Attr: "error", Value: "battery_low"},
			{Attr: "error", Value: "last_downlink"},
			{Attr: "error", Value: "tamper_since_reset"},
			{Attr: "tamper", Value: true},
			{Attr: "battery_volt", Value: 2.7},
		}, ""},
		// Supervisory event count.
		{"1701080130ffffffff1234", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(7)},
			{Attr: "tamper", Value: false},
			{Attr: "battery_volt", Value: float64(3)},
			{Attr: "total_count", Value: int32(4660)},
		}, ""},
		{"1401080132000000000002", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(4)},
			{Attr: "tamper", Value: false},
			{Attr: "battery_volt", Value: 3.2},
			{Attr: "total_count", Value: int32(2)},
		}, ""},
		// Supervisory bad length.
		{"", nil, "supervisory format bad length: "},
		// Supervisory bad identifier.
		{"1402080131", nil, "supervisory format bad identifier: 1402080131"},
		// Supervisory bad error bitmap.
		{"1401400131", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(4)},
		}, "supervisory format bad error bitmap: 1401400131"},
		// Supervisory event count unused trailing bytes.
		{"1701080130ffffffff1234ff", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(7)},
			{Attr: "tamper", Value: false},
			{Attr: "battery_volt", Value: float64(3)},
			{Attr: "total_count", Value: int32(4660)},
		}, "supervisory format unused trailing bytes: " +
			"1701080130ffffffff1234ff"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, err := supervisory(bInp)
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
