// +build !integration

package device

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/parse"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDeviceJoin(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)
	bUniqID, err := hex.DecodeString(uniqID)
	require.NoError(t, err)
	t.Logf("bUniqID: %x", bUniqID)

	b64UniqID := base64.StdEncoding.EncodeToString(bUniqID)
	t.Logf("b64UniqID: %v", b64UniqID)

	devAddr := random.String(16)
	bDevAddr, err := hex.DecodeString(devAddr)
	require.NoError(t, err)
	t.Logf("bDevAddr: %x", bDevAddr)

	b64DevAddr := base64.StdEncoding.EncodeToString(bDevAddr)
	t.Logf("b64DevAddr: %v", b64DevAddr)

	gatewayID := random.String(16)
	bGatewayID, err := hex.DecodeString(gatewayID)
	require.NoError(t, err)
	t.Logf("bGatewayID: %x", bGatewayID)

	b64GatewayID := base64.StdEncoding.EncodeToString(bGatewayID)
	t.Logf("b64GatewayID: %v", b64GatewayID)

	// Truncate to nearest second for compatibility with jsonpb.Marshaler and
	// time.RFC3339Nano formatting.
	now := time.Now().UTC().Add(-15 * time.Minute).Truncate(time.Second)
	pNow := timestamppb.New(now)

	// Device Join payloads, see deviceJoin() for format description.
	tests := []struct {
		inp       *as.JoinEvent
		resPoints []*parse.Point
		resTime   time.Time
		err       string
	}{
		// Device Join.
		{&as.JoinEvent{RxInfo: []*gw.UplinkRXInfo{{}}},
			[]*parse.Point{
				{Attr: "raw_device", Value: `{"rxInfo":[{}]}`},
				{Attr: "join", Value: true},
				{Attr: "data_rate", Value: 0},
			}, time.Now(), ""},
		{&as.JoinEvent{DevEui: bUniqID, DevAddr: bDevAddr,
			RxInfo: []*gw.UplinkRXInfo{{GatewayId: []byte("aaa"), Time: pNow,
				Rssi: -80, LoraSnr: 1}, {GatewayId: bGatewayID, Time: pNow,
				Rssi: -74, LoraSnr: 7.8}}, TxInfo: &gw.UplinkTXInfo{
				Frequency: 902700000}, Dr: 3}, []*parse.Point{
			{Attr: "raw_device", Value: fmt.Sprintf(`{"devEUI":"%s","devAddr":`+
				`"%s","rxInfo":[{"gatewayID":"YWFh","time":"%s","rssi":-80,`+
				`"loRaSNR":1},{"gatewayID":"%s","time":"%s","rssi":-74,`+
				`"loRaSNR":7.8}],"txInfo":{"frequency":902700000},"dr":3}`,
				b64UniqID, b64DevAddr, now.Format(time.RFC3339Nano),
				b64GatewayID, now.Format(time.RFC3339Nano))},
			{Attr: "join", Value: true},
			{Attr: "id", Value: uniqID},
			{Attr: "devaddr", Value: devAddr},
			{Attr: "gateway_id", Value: gatewayID},
			{Attr: "time", Value: int(now.Unix())},
			{Attr: "rssi", Value: -74},
			{Attr: "snr", Value: 7.8},
			{Attr: "frequency", Value: 902700000},
			{Attr: "data_rate", Value: 3},
		}, now, ""},
		// Device Join bad length.
		{nil, nil, time.Time{}, "unexpected EOF"},
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

			res, ts, err := deviceJoin(bInp)
			t.Logf("res, ts, err: %#v, %v, %v", res, ts, err)
			require.Equal(t, lTest.resPoints, res)
			require.WithinDuration(t, lTest.resTime, ts, 2*time.Second)
			if lTest.err != "" {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
