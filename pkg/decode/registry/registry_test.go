// +build !integration

package registry

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestNew(t *testing.T) {
	t.Parallel()

	reg := New()
	t.Logf("reg: %#v", reg)
	require.NotNil(t, reg)
}

func TestDecode(t *testing.T) {
	t.Parallel()

	reg := New()

	tests := []struct {
		inpDecoder api.Decoder
		inpBody    string
		res        []*decode.Point
		err        error
	}{
		// Decoder.
		{api.Decoder_RAW, "", nil, nil},
		{api.Decoder_RADIO_BRIDGE_DOOR_V2, "190301", []*decode.Point{
			{Attr: "count", Value: int32(9)},
			{Attr: "open", Value: true},
		}, nil},
		// Decoder function not found.
		{api.Decoder(999), "", nil, fmt.Errorf("%w: 999", ErrNotFound)},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInpBody, err := hex.DecodeString(lTest.inpBody)
			require.NoError(t, err)

			res, err := reg.Decode(lTest.inpDecoder, bInpBody)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			require.Equal(t, lTest.err, err)
		})
	}
}
