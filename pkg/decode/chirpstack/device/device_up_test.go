//go:build !integration

package device

//nolint:staticcheck // third-party dependency
import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDeviceUp(t *testing.T) {
	t.Parallel()

	gatewayID := random.String(16)
	bGatewayID, err := hex.DecodeString(gatewayID)
	require.NoError(t, err)
	t.Logf("bGatewayID: %x", bGatewayID)

	b64GatewayID := base64.StdEncoding.EncodeToString(bGatewayID)
	t.Logf("b64GatewayID: %v", b64GatewayID)

	// Truncate to nearest second for compatibility with jsonpb.Marshaler and
	// time.RFC3339Nano formatting.
	now := time.Now().UTC().Add(-15 * time.Minute).Truncate(time.Second)
	tsNow := timestamppb.New(now)

	bData := random.Bytes(10)
	b64Data := base64.StdEncoding.EncodeToString(bData)
	t.Logf("b64Data: %v", b64Data)

	// Device Uplink payloads, see deviceUp() for format description.
	tests := []struct {
		inp       *as.UplinkEvent
		resPoints []*decode.Point
		resTime   time.Time
		resData   []byte
		err       string
	}{
		// Device Uplink.
		{
			&as.UplinkEvent{RxInfo: []*gw.UplinkRXInfo{{}}}, []*decode.Point{
				{Attr: "raw_device", Value: `{"rxInfo":[{}]}`},
				{Attr: "channel", Value: int32(0)},
				{Attr: "adr", Value: false},
				{Attr: "data_rate", Value: int32(0)},
				{Attr: "confirmed", Value: false},
			}, time.Now(), nil, "",
		},
		{
			&as.UplinkEvent{
				RxInfo: []*gw.UplinkRXInfo{
					{
						GatewayId: []byte("aaa"), Time: tsNow, Rssi: -80,
						LoraSnr: 1,
					}, {
						GatewayId: bGatewayID, Time: tsNow, Rssi: -74,
						LoraSnr: 7.8,
					},
				}, TxInfo: &gw.UplinkTXInfo{Frequency: 902700000}, Adr: true,
				Dr: 3, Data: bData, ConfirmedUplink: true,
			}, []*decode.Point{
				{Attr: "raw_device", Value: fmt.Sprintf(`{"rxInfo":[{`+
					`"gatewayID":"YWFh","time":"%s","rssi":-80,"loRaSNR":1},{`+
					`"gatewayID":"%s","time":"%s","rssi":-74,"loRaSNR":7.8}],`+
					`"txInfo":{"frequency":902700000},"adr":true,"dr":3,`+
					`"data":"%s","confirmedUplink":true}`,
					now.Format(time.RFC3339Nano), b64GatewayID,
					now.Format(time.RFC3339Nano), b64Data)},
				{Attr: "raw_data", Value: hex.EncodeToString(bData)},
				{Attr: "gateway_id", Value: gatewayID},
				{Attr: "time", Value: strconv.FormatInt(now.Unix(), 10)},
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "snr", Value: 7.8},
				{Attr: "channel", Value: int32(0)},
				{Attr: "frequency", Value: int32(902700000)},
				{Attr: "adr", Value: true},
				{Attr: "data_rate", Value: int32(3)},
				{Attr: "confirmed", Value: true},
			}, now, bData, "",
		},
		// Device Uplink bad length.
		{
			nil, nil, time.Time{}, nil, "cannot parse invalid wire-format data",
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

			res, ts, data, err := deviceUp(bInp)
			t.Logf("res, ts, data, err: %#v, %v, %x, %v", res, ts, data, err)
			require.Equal(t, lTest.resPoints, res)
			if !lTest.resTime.IsZero() {
				require.WithinDuration(t, lTest.resTime, ts.AsTime(),
					2*time.Second)
			}
			require.Equal(t, lTest.resData, data)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
