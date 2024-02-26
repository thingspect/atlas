//go:build !integration

package registry

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/proto/go/api"
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
		{api.Decoder_GLOBALSAT_CO2, "01096113950292", []*decode.Point{
			{Attr: "temp_c", Value: 24.0},
			{Attr: "temp_f", Value: 75.2},
			{Attr: "humidity_pct", Value: 50.13},
			{Attr: "co2_ppm", Value: int32(658)},
		}, nil},
		{api.Decoder_TEKTELIC_HOME, "036700c404687f00ff0138", []*decode.Point{
			{Attr: "temp_c", Value: 19.6},
			{Attr: "temp_f", Value: 67.3},
			{Attr: "humidity_pct", Value: 63.5},
			{Attr: "battery_v", Value: 3.12},
		}, nil},
		// Decoder function not found.
		{api.Decoder(999), "", nil, fmt.Errorf("%w: 999", ErrNotFound)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can decode %+v", test), func(t *testing.T) {
			t.Parallel()

			bInpBody, err := hex.DecodeString(test.inpBody)
			require.NoError(t, err)

			res, err := reg.Decode(test.inpDecoder, bInpBody)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, test.res, res)
			require.Equal(t, test.err, err)
		})
	}
}
