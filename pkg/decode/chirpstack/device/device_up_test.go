//go:build !integration

package device

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDeviceUp(t *testing.T) {
	t.Parallel()

	gatewayID := random.String(16)

	// Truncate to nearest second for compatibility with protojson.Format and
	// time.RFC3339Nano formatting.
	now := time.Now().UTC().Add(-15 * time.Minute).Truncate(time.Second)
	tsNow := timestamppb.New(now)

	bData := random.Bytes(10)
	b64Data := base64.StdEncoding.EncodeToString(bData)
	t.Logf("b64Data: %v", b64Data)

	// Device Uplink payloads, see deviceUp() for format description.
	tests := []struct {
		inp       *integration.UplinkEvent
		resPoints []*decode.Point
		resTime   time.Time
		resData   []byte
		err       string
	}{
		// Device Uplink.
		{
			&integration.UplinkEvent{}, []*decode.Point{
				{Attr: "raw_device", Value: `{}`},
				{Attr: "adr", Value: false},
				{Attr: "data_rate", Value: int32(0)},
				{Attr: "confirmed", Value: false},
			}, time.Now(), nil, "",
		},
		{
			&integration.UplinkEvent{RxInfo: []*gw.UplinkRxInfo{{}}},
			[]*decode.Point{
				{Attr: "raw_device", Value: `{"rxInfo":[{}]}`},
				{Attr: "channel", Value: int32(0)},
				{Attr: "adr", Value: false},
				{Attr: "data_rate", Value: int32(0)},
				{Attr: "confirmed", Value: false},
			},
			time.Now(), nil, "",
		},
		{
			&integration.UplinkEvent{
				RxInfo: []*gw.UplinkRxInfo{
					{GatewayId: "aaa", GwTime: tsNow, Rssi: -80, Snr: 1},
					{GatewayId: gatewayID, GwTime: tsNow, Rssi: -74, Snr: 7},
				}, TxInfo: &gw.UplinkTxInfo{Frequency: 902700000}, Adr: true,
				Dr: 3, Data: bData, Confirmed: true, RegionConfigId: "us915_0",
			},
			[]*decode.Point{
				{Attr: "raw_device", Value: fmt.Sprintf(`{"adr":true,"dr":3,`+
					`"confirmed":true,"data":"%s","rxInfo":[{"gatewayId":`+
					`"aaa","gwTime":"%s","rssi":-80,"snr":1},{"gatewayId":`+
					`"%s","gwTime":"%s","rssi":-74,"snr":7}],"txInfo":{`+
					`"frequency":902700000},"regionConfigId":"us915_0"}`,
					b64Data, now.Format(time.RFC3339Nano), gatewayID,
					now.Format(time.RFC3339Nano))},
				{Attr: "raw_data", Value: hex.EncodeToString(bData)},
				{Attr: "gateway_id", Value: gatewayID},
				{
					Attr:  "gateway_time",
					Value: strconv.FormatInt(now.Unix(), 10),
				},
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "lora_snr", Value: float64(7)},
				{Attr: "channel", Value: int32(0)},
				{Attr: "frequency", Value: int32(902700000)},
				{Attr: "adr", Value: true},
				{Attr: "data_rate", Value: int32(3)},
				{Attr: "confirmed", Value: true},
				{Attr: "region_config_id", Value: "us915_0"},
			},
			now, bData, "",
		},
		// Device Uplink bad length.
		{
			nil, nil, time.Time{}, nil, "cannot parse invalid wire-format data",
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

			res, ts, data, err := deviceUp(bInp)
			t.Logf("res, ts, data, err: %#v, %v, %x, %v", res, ts, data, err)
			require.Equal(t, test.resPoints, res)
			if !test.resTime.IsZero() {
				require.WithinDuration(t, test.resTime, ts.AsTime(),
					2*time.Second)
			}
			require.Equal(t, test.resData, data)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
