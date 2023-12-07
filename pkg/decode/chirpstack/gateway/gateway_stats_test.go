//go:build !integration

package gateway

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGatewayStats(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)

	// Truncate to nearest second for compatibility with protojson.Format and
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
				GatewayId: uniqID, Time: pNow,
				RxPacketsReceived: 1, RxPacketsReceivedOk: 2,
				TxPacketsReceived: 3, TxPacketsEmitted: 4,
				Metadata: map[string]string{"aaa": "bbb"},
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: fmt.Sprintf(`{"gatewayId":"%s",`+
					`"time":"%s","rxPacketsReceived":1,`+
					`"rxPacketsReceivedOk":2,"txPacketsReceived":3,`+
					`"txPacketsEmitted":4,"metadata":{"aaa":"bbb"}}`, uniqID,
					now.Format(time.RFC3339Nano))},
				{Attr: "id", Value: uniqID},
				{
					Attr:  "gateway_time",
					Value: strconv.FormatInt(now.Unix(), 10),
				},
				{Attr: "rx_received", Value: int32(1)},
				{Attr: "rx_received_valid", Value: int32(2)},
				{Attr: "tx_received", Value: int32(3)},
				{Attr: "tx_transmitted", Value: int32(4)},
				{Attr: "aaa", Value: "bbb"},
			}, "",
		},
		// Gateway Stats bad length.
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
