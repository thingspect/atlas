// +build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/parse"
)

func TestDoor(t *testing.T) {
	t.Parallel()

	// Door payloads, see Door() for format description.
	tests := []struct {
		inp string
		res []*parse.Point
		err string
	}{
		// Reset.
		{"100001120102181c", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 0},
			{Attr: "hw_ver", Value: 18},
			{Attr: "ver", Value: "1.2"},
		}, ""},
		// Supervisory.
		{"1401080131", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 4},
			{Attr: "tamper", Value: false},
			{Attr: "battery", Value: 3.1},
		}, ""},
		// Tamper.
		{"1c0200", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 12},
			{Attr: "tamper", Value: true},
		}, ""},
		// Link Quality.
		{"1dfb010000", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 13},
			{Attr: "sub_band", Value: 1},
			{Attr: "dev_rssi", Value: 0},
			{Attr: "dev_snr", Value: 0},
		}, ""},
		// Door.
		{"190301", []*parse.Point{
			{Attr: "count", Value: 9},
			{Attr: "open", Value: true},
		}, ""},
		{"1a0300", []*parse.Point{
			{Attr: "count", Value: 10},
			{Attr: "open", Value: false},
		}, ""},
		// Door bad length.
		{"", nil, "door format bad length: "},
		// Door bad identifier.
		{"190401", nil, "door format bad identifier: 190401"},
		// Door unused trailing bytes.
		{"190301ff", []*parse.Point{
			{Attr: "count", Value: 9},
			{Attr: "open", Value: true},
		}, "door format unused trailing bytes: 190301ff"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)
			t.Logf("bInp: %x", bInp)

			res, err := Door(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err != "" {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
