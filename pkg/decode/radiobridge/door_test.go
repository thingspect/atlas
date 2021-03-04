// +build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestDoor(t *testing.T) {
	t.Parallel()

	// Door payloads, see Door() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// Reset.
		{"100001120102181c", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(0)},
			{Attr: "hw_ver", Value: int32(18)},
			{Attr: "ver", Value: "1.2"},
		}, ""},
		// Supervisory.
		{"1401080131", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(4)},
			{Attr: "tamper", Value: false},
			{Attr: "battery", Value: 3.1},
		}, ""},
		// Tamper.
		{"1c0200", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(12)},
			{Attr: "tamper", Value: true},
		}, ""},
		// Link Quality.
		{"1dfb010000", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(13)},
			{Attr: "sub_band", Value: int32(1)},
			{Attr: "device_rssi", Value: int32(0)},
			{Attr: "device_snr", Value: int32(0)},
		}, ""},
		// Door.
		{"190301", []*decode.Point{
			{Attr: "count", Value: int32(9)},
			{Attr: "open", Value: true},
		}, ""},
		{"1a0300", []*decode.Point{
			{Attr: "count", Value: int32(10)},
			{Attr: "open", Value: false},
		}, ""},
		// Door bad length.
		{"", nil, "door format bad length: "},
		// Door bad identifier.
		{"190401", nil, "door format bad identifier: 190401"},
		// Door unused trailing bytes.
		{"190301ff", []*decode.Point{
			{Attr: "count", Value: int32(9)},
			{Attr: "open", Value: true},
		}, "door format unused trailing bytes: 190301ff"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, err := Door(bInp)
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
