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

func TestGatewayUp(t *testing.T) {
	t.Parallel()

	// Gateway Uplink payloads, see gatewayUp() for format description.
	tests := []struct {
		inp *gw.UplinkFrame
		res []*decode.Point
		err string
	}{
		// Gateway Uplink.
		{
			&gw.UplinkFrame{}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{}`},
			}, "",
		},
		{
			&gw.UplinkFrame{RxInfo: &gw.UplinkRxInfo{}}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"rxInfo":{}}`},
				{Attr: "channel", Value: int32(0)},
			}, "",
		},
		{
			&gw.UplinkFrame{
				TxInfo: &gw.UplinkTxInfo{
					Frequency: 902700000, Modulation: &gw.Modulation{
						Parameters: &gw.Modulation_Lora{
							Lora: &gw.LoraModulationInfo{SpreadingFactor: 7},
						},
					},
				}, RxInfo: &gw.UplinkRxInfo{Rssi: -74, Snr: 7, Channel: 2},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"txInfo":{"frequency":` +
					`902700000,"modulation":{"lora":{"spreadingFactor":7}}},` +
					`"rxInfo":{"rssi":-74,"snr":7,"channel":2}}`},
				{Attr: "frequency", Value: int32(902700000)},
				{Attr: "sf", Value: int32(7)},
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "lora_snr", Value: float64(7)},
				{Attr: "channel", Value: int32(2)},
			}, "",
		},
		// Gateway Uplink bad length.
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

			res, err := gatewayUp(bInp)
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
