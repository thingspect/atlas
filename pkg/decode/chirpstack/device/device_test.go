//go:build !integration

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

func TestParse(t *testing.T) {
	t.Parallel()

	badEvent := random.String(10)

	// Trivial device payloads, see Parse() for format description. Parsers
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
		}, nil},
		{"ack", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ack", Value: ackTimeout},
		}, nil},
		{"log", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "log_level", Value: "INFO"},
			{Attr: "log_code", Value: "UNKNOWN"},
		}, nil},
		{"txack", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "tx_queued", Value: true},
		}, nil},
		{"status", "", []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ext_power", Value: false},
		}, nil},
		// Device unknown event type.
		{badEvent, "", nil, fmt.Errorf("%w: %s, %x", decode.ErrUnknownEvent,
			badEvent, []byte{})},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInpBody, err := hex.DecodeString(test.inpBody)
			require.NoError(t, err)

			res, ts, data, err := Parse(test.inpEvent, bInpBody)
			t.Logf("res, ts, data, err: %#v, %v, %x, %v", res, ts, data, err)
			require.Equal(t, test.res, res)
			if ts != nil {
				require.WithinDuration(t, time.Now(), ts.AsTime(),
					2*time.Second)
			}
			require.Nil(t, data)
			require.Equal(t, test.err, err)
		})
	}
}
