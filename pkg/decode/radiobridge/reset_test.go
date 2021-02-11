// +build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestReset(t *testing.T) {
	t.Parallel()

	// Reset payloads, see reset() for format description.
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
		{"100001127fff181c", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(0)},
			{Attr: "hw_ver", Value: int32(18)},
			{Attr: "ver", Value: "127.255"},
		}, ""},
		{"100001128823181c", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(0)},
			{Attr: "hw_ver", Value: int32(18)},
			{Attr: "ver", Value: "2.1.3"},
		}, ""},
		{"100001128801181c", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(0)},
			{Attr: "hw_ver", Value: int32(18)},
			{Attr: "ver", Value: "2.0.1"},
		}, ""},
		{"10000112ffff181c", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(0)},
			{Attr: "hw_ver", Value: int32(18)},
			{Attr: "ver", Value: "31.31.31"},
		}, ""},
		// Reset bad length.
		{"", nil, "reset format bad length: "},
		// Reset bad identifier.
		{"100101120102181c", nil, "reset format bad identifier: " +
			"100101120102181c"},
		// Reset unused trailing bytes.
		{"100001120102181cff", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(0)},
			{Attr: "hw_ver", Value: int32(18)},
			{Attr: "ver", Value: "1.2"},
		}, "reset format unused trailing bytes: 100001120102181cff"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, err := reset(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err != "" {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
