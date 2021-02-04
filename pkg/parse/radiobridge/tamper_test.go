// +build !integration

package radiobridge

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/parse"
)

func TestTamper(t *testing.T) {
	t.Parallel()

	// Tamper payloads, see tamper() for format description.
	tests := []struct {
		inp string
		res []*parse.Point
		err string
	}{
		// Tamper.
		{"1c0200", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 12},
			{Attr: "tamper", Value: true},
		}, ""},
		{"1d0201", []*parse.Point{
			{Attr: "proto", Value: 1},
			{Attr: "count", Value: 13},
			{Attr: "tamper", Value: false},
		}, ""},
		// Tamper bad length.
		{"", nil, "tamper format bad length: "},
		// Tamper bad identifier.
		{"1c0300", nil, "tamper format bad identifier: 1c0300"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)
			t.Logf("bInp: %x", bInp)

			res, err := tamper(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err != "" {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
