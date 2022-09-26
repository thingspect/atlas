//go:build !integration

package device

import (
	"fmt"
	"testing"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/proto"
)

func TestDeviceAck(t *testing.T) {
	t.Parallel()

	// Device ACK payloads, see deviceAck() for format description.
	tests := []struct {
		inp *integration.AckEvent
		res []*decode.Point
		err string
	}{
		// Device ACK.
		{&integration.AckEvent{}, []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ack", Value: ackTimeout},
		}, ""},
		{&integration.AckEvent{Acknowledged: true}, []*decode.Point{
			{Attr: "raw_device", Value: `{"acknowledged":true}`},
			{Attr: "ack", Value: ackOK},
		}, ""},
		// Device ACK bad length.
		{nil, nil, "cannot parse invalid wire-format data"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp := []byte("aaa")
			if lTest.inp != nil {
				var err error
				bInp, err = proto.Marshal(lTest.inp)
				require.NoError(t, err)
			}

			res, err := deviceAck(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
