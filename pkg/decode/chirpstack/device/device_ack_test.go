// +build !integration

package device

import (
	"fmt"
	"testing"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestDeviceAck(t *testing.T) {
	t.Parallel()

	// Device ACK payloads, see deviceAck() for format description.
	tests := []struct {
		inp *as.AckEvent
		res []*decode.Point
		err string
	}{
		// Device ACK.
		{&as.AckEvent{}, []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ack", Value: ackTimeout},
		}, ""},
		{&as.AckEvent{Acknowledged: true}, []*decode.Point{
			{Attr: "raw_device", Value: `{"acknowledged":true}`},
			{Attr: "ack", Value: ackOK},
		}, ""},
		// Device ACK bad length.
		{nil, nil, "unexpected EOF"},
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
			if lTest.err != "" {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
