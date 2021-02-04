// +build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/parse"
)

func TestLinkQuality(t *testing.T) {
	t.Parallel()

	// Link Quality payloads, see linkQuality() for format description.
	tests := []struct {
		inp string
		res []*parse.Point
		err string
	}{
		// Link Quality.
		{"1dfb010000", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 13},
			{Attr: "sub_band", Value: 1},
			{Attr: "dev_rssi", Value: 0},
			{Attr: "dev_snr", Value: 0},
		}, ""},
		{"1dfb01ca0b", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 13},
			{Attr: "sub_band", Value: 1},
			{Attr: "dev_rssi", Value: -54},
			{Attr: "dev_snr", Value: 11},
		}, ""},
		// Link Quality bad length.
		{"", nil, "link quality format bad length: "},
		// Link Quality bad identifier.
		{"1dfc010000", nil, "link quality format bad identifier: 1dfc010000"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)
			t.Logf("bInp: %x", bInp)

			res, err := linkQuality(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err != "" {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
