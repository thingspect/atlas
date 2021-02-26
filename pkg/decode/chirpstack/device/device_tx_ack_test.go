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

func TestDeviceTxAck(t *testing.T) {
	t.Parallel()

	// Device TX ACK payloads, see deviceTxAck() for format description.
	tests := []struct {
		inp *as.TxAckEvent
		res []*decode.Point
		err string
	}{
		// Device TX ACK.
		{&as.TxAckEvent{}, []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ack_gateway_tx", Value: true},
		}, ""},
		// Device TX ACK bad length.
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

			res, err := deviceTxAck(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err != "" {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
