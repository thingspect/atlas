//go:build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestLinkQuality(t *testing.T) {
	t.Parallel()

	// Link Quality payloads, see linkQuality() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// Link Quality.
		{"1dfb010000", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(13)},
			{Attr: "sub_band", Value: int32(1)},
			{Attr: "device_rssi", Value: int32(0)},
			{Attr: "device_snr", Value: int32(0)},
		}, ""},
		{"1dfb01ca0b", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(13)},
			{Attr: "sub_band", Value: int32(1)},
			{Attr: "device_rssi", Value: int32(-54)},
			{Attr: "device_snr", Value: int32(11)},
		}, ""},
		// Link Quality bad length.
		{"00010203", nil, "link quality format bad length: 00010203"},
		{"000102030405", nil, "link quality format bad length: 000102030405"},
		// Link Quality bad identifier.
		{"1dfc010000", nil, "link quality format bad identifier: 1dfc010000"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(test.inp)
			require.NoError(t, err)

			res, err := linkQuality(bInp)
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
