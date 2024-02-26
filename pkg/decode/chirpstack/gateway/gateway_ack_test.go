//go:build !integration

package gateway

import (
	"fmt"
	"testing"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/proto"
)

func TestGatewayAck(t *testing.T) {
	t.Parallel()

	// Gateway ACK payloads, see gatewayAck() for format description.
	tests := []struct {
		inp *gw.DownlinkTxAck
		res []*decode.Point
		err string
	}{
		// Gateway ACK.
		{
			&gw.DownlinkTxAck{}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{}`},
			}, "",
		},
		{
			&gw.DownlinkTxAck{
				Items: []*gw.DownlinkTxAckItem{{Status: gw.TxAckStatus_OK}},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"items":[{"status":"OK"}]}`},
				{Attr: "ack", Value: "OK"},
			}, "",
		},
		{
			&gw.DownlinkTxAck{
				Items: []*gw.DownlinkTxAckItem{
					{Status: gw.TxAckStatus_TOO_LATE},
				},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"items":[{"status":` +
					`"TOO_LATE"}]}`},
				{Attr: "ack", Value: "TOO_LATE"},
			}, "",
		},
		// Gateway ACK bad length.
		{
			nil, nil, "cannot parse invalid wire-format data",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp := []byte("aaa")
			if test.inp != nil {
				var err error
				bInp, err = proto.Marshal(test.inp)
				require.NoError(t, err)
			}

			res, err := gatewayAck(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, test.res, res)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
