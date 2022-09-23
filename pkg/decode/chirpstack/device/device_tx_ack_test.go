//go:build !integration

package device

import (
	"fmt"
	"testing"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

func TestDeviceTXAck(t *testing.T) {
	t.Parallel()

	gatewayID := random.String(16)

	// Device TX ACK payloads, see deviceTXAck() for format description.
	tests := []struct {
		inp *integration.TxAckEvent
		res []*decode.Point
		err string
	}{
		// Device TX ACK.
		{
			&integration.TxAckEvent{}, []*decode.Point{
				{Attr: "raw_device", Value: `{}`},
				{Attr: "tx_queued", Value: true},
			}, "",
		},
		{
			&integration.TxAckEvent{
				GatewayId: gatewayID,
				TxInfo:    &gw.DownlinkTxInfo{Frequency: 902700000},
			}, []*decode.Point{
				{Attr: "raw_device", Value: fmt.Sprintf(`{"gatewayId":"%s",`+
					`"txInfo":{"frequency":902700000}}`, gatewayID)},
				{Attr: "tx_queued", Value: true},
				{Attr: "tx_gateway_id", Value: gatewayID},
				{Attr: "tx_frequency", Value: int32(902700000)},
			}, "",
		},
		// Device TX ACK bad length.
		{
			nil, nil, "cannot parse invalid wire-format data",
		},
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

			res, err := deviceTXAck(bInp)
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
