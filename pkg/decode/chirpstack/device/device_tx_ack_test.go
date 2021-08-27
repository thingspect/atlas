//go:build !integration

package device

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestDeviceTxAck(t *testing.T) {
	t.Parallel()

	gatewayID := random.String(16)
	bGatewayID, err := hex.DecodeString(gatewayID)
	require.NoError(t, err)
	t.Logf("bGatewayID: %x", bGatewayID)

	b64GatewayID := base64.StdEncoding.EncodeToString(bGatewayID)
	t.Logf("b64GatewayID: %v", b64GatewayID)

	// Device TX ACK payloads, see deviceTxAck() for format description.
	tests := []struct {
		inp *as.TxAckEvent
		res []*decode.Point
		err string
	}{
		// Device TX ACK.
		{
			&as.TxAckEvent{}, []*decode.Point{
				{Attr: "raw_device", Value: `{}`},
				{Attr: "ack_gateway_tx", Value: true},
			}, "",
		},
		{
			&as.TxAckEvent{
				GatewayId: bGatewayID,
				TxInfo:    &gw.DownlinkTXInfo{Frequency: 902700000},
			}, []*decode.Point{
				{Attr: "raw_device", Value: fmt.Sprintf(`{"gatewayID":"%s",`+
					`"txInfo":{"frequency":902700000}}`, b64GatewayID)},
				{Attr: "ack_gateway_tx", Value: true},
				{Attr: "gateway_id", Value: gatewayID},
				{Attr: "frequency", Value: int32(902700000)},
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

			res, err := deviceTxAck(bInp)
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
