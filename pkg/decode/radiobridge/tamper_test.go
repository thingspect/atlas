//go:build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestTamper(t *testing.T) {
	t.Parallel()

	// Tamper payloads, see tamper() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		err string
	}{
		// Tamper.
		{"1c0200", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(12)},
			{Attr: "tamper", Value: true},
		}, ""},
		{"1d0201", []*decode.Point{
			{Attr: "proto", Value: int32(1)},
			{Attr: "count", Value: int32(13)},
			{Attr: "tamper", Value: false},
		}, ""},
		// Tamper bad length.
		{"0001", nil, "tamper format bad length: 0001"},
		{"00010203", nil, "tamper format bad length: 00010203"},
		// Tamper bad identifier.
		{"1c0300", nil, "tamper format bad identifier: 1c0300"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(test.inp)
			require.NoError(t, err)

			res, err := tamper(bInp)
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
