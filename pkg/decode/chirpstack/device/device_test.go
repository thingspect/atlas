// +build !integration

package device

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestDevice(t *testing.T) {
	t.Parallel()

	badEvent := random.String(10)

	// Trivial device payloads, see Device() for format description. Parsers
	// are exercised more thoroughly in their respective tests.
	tests := []struct {
		inpEvent string
		inpBody  string
		res      []*decode.Point
		err      error
	}{
		// Device.
		{"up", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "adr", Value: false},
			{Attr: "data_rate", Value: int32(0)},
			{Attr: "confirmed", Value: false},
		}, nil},
		{"join", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "join", Value: true},
			{Attr: "data_rate", Value: int32(0)},
		}, nil},
		{"ack", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ack", Value: ackTimeout},
		}, nil},
		{"error", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "error_type", Value: "UNKNOWN"},
		}, nil},
		{"txack", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ack_gateway_tx", Value: true},
		}, nil},
		// Device unknown event type.
		{badEvent, "", nil, fmt.Errorf("%w: %s, %x", decode.ErrUnknownEvent,
			badEvent, []byte{})},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInpBody, err := hex.DecodeString(lTest.inpBody)
			require.NoError(t, err)

			res, ts, data, err := Device(lTest.inpEvent, bInpBody)
			t.Logf("res, ts, data, err: %#v, %v, %x, %v", res, ts, data, err)
			require.Equal(t, lTest.res, res)
			if ts != nil {
				require.WithinDuration(t, time.Now(), ts.AsTime(),
					2*time.Second)
			}
			require.Nil(t, data)
			require.Equal(t, lTest.err, err)
		})
	}
}
