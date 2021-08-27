//go:build !integration

package gateway

import (
	"fmt"
	"testing"

	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestGatewayAck(t *testing.T) {
	t.Parallel()

	// Gateway ACK payloads, see gatewayAck() for format description.
	tests := []struct {
		inp *gw.DownlinkTXAck
		res []*decode.Point
		err string
	}{
		// Gateway ACK.
		{
			&gw.DownlinkTXAck{}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{}`},
			}, "",
		},
		{
			&gw.DownlinkTXAck{
				Items: []*gw.DownlinkTXAckItem{{Status: gw.TxAckStatus_OK}},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"items":[{"status":"OK"}]}`},
				{Attr: "ack", Value: "OK"},
			}, "",
		},
		{
			&gw.DownlinkTXAck{
				Items: []*gw.DownlinkTXAckItem{
					{Status: gw.TxAckStatus_TOO_LATE},
				},
			}, []*decode.Point{
				{
					Attr:  "raw_gateway",
					Value: `{"items":[{"status":` + `"TOO_LATE"}]}`,
				}, {
					Attr: "ack", Value: "TOO_LATE",
				},
			}, "",
		},
		// Gateway ACK bad length.
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

			res, err := gatewayAck(bInp)
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
