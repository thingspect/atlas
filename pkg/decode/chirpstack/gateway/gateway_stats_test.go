//go:build !integration

package gateway

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	//nolint:staticcheck // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGatewayStats(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)
	bUniqID, err := hex.DecodeString(uniqID)
	require.NoError(t, err)
	t.Logf("bUniqID: %x", bUniqID)

	b64UniqID := base64.StdEncoding.EncodeToString(bUniqID)
	t.Logf("b64UniqID: %v", b64UniqID)

	// Truncate to nearest second for compatibility with jsonpb.Marshaler and
	// time.RFC3339Nano formatting.
	now := time.Now().UTC().Add(-15 * time.Minute).Truncate(time.Second)
	pNow := timestamppb.New(now)

	// Gateway Stats payloads, see gatewayStats() for format description.
	tests := []struct {
		inp *gw.GatewayStats
		res []*decode.Point
		err string
	}{
		// Gateway Stats.
		{
			&gw.GatewayStats{}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{}`},
			}, "",
		},
		{
			&gw.GatewayStats{
				GatewayId: bUniqID, Ip: "127.0.0.1", Time: pNow,
				RxPacketsReceived: 1, RxPacketsReceivedOk: 2,
				TxPacketsReceived: 3, TxPacketsEmitted: 4,
				MetaData: map[string]string{"aaa": "bbb"},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: fmt.Sprintf(`{"gatewayID":"%s",`+
					`"ip":"127.0.0.1","time":"%s","rxPacketsReceived":1,`+
					`"rxPacketsReceivedOK":2,"txPacketsReceived":3,`+
					`"txPacketsEmitted":4,"metaData":{"aaa":"bbb"}}`, b64UniqID,
					now.Format(time.RFC3339Nano))},
				{Attr: "id", Value: uniqID},
				{Attr: "ip", Value: "127.0.0.1"},
				{Attr: "time", Value: strconv.FormatInt(now.Unix(), 10)},
				{Attr: "rx_received", Value: int32(1)},
				{Attr: "rx_received_valid", Value: int32(2)},
				{Attr: "tx_received", Value: int32(3)},
				{Attr: "tx_transmitted", Value: int32(4)},
				{Attr: "aaa", Value: "bbb"},
			}, "",
		},
		{
			&gw.GatewayStats{
				MetaData: map[string]string{"str_cell_status": "DISCONNECTED"},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"metaData":{"str_cell_status":` +
					`"DISCONNECTED"}}`},
				{Attr: "cell_status", Value: "DISCONNECTED"},
			}, "",
		},
		{
			&gw.GatewayStats{
				MetaData: map[string]string{"int_hour": "1612391935"},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"metaData":{"int_hour":` +
					`"1612391935"}}`},
				{Attr: "hour", Value: int32(1612391935)},
			}, "",
		},
		{
			&gw.GatewayStats{
				MetaData: map[string]string{"fl64_uptime": "13379004.3"},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"metaData":{"fl64_uptime":` +
					`"13379004.3"}}`},
				{Attr: "uptime", Value: 13379004.3},
			}, "",
		},
		// Gateway Stats bad length.
		{
			nil, nil, "cannot parse invalid wire-format data",
		},
		// Gateway Stats bad metadata int conversion.
		{
			&gw.GatewayStats{MetaData: map[string]string{"int_hour": "aaa"}},
			[]*decode.Point{
				{Attr: "raw_gateway", Value: `{"metaData":{"int_hour":"aaa"}}`},
			},
			`strconv.ParseInt: parsing "aaa": invalid syntax`,
		},
		// Gateway Stats bad metadata float64 conversion.
		{
			&gw.GatewayStats{MetaData: map[string]string{"fl64_time": "bbb"}},
			[]*decode.Point{
				{Attr: "raw_gateway", Value: `{"metaData":{"fl64_time":` +
					`"bbb"}}`},
			},
			`strconv.ParseFloat: parsing "bbb": invalid syntax`,
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

			res, err := gatewayStats(bInp)
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
