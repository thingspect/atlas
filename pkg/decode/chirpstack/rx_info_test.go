//go:build !integration

package chirpstack

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseRXInfo(t *testing.T) {
	t.Parallel()

	// Gateway UplinkRxInfo payloads, see ParseRXInfo() for format description.
	tests := []struct {
		inp *gw.UplinkRxInfo
		res []*decode.Point
	}{
		// Gateway UplinkRxInfo.
		{
			&gw.UplinkRxInfo{}, []*decode.Point{
				{Attr: "channel", Value: int32(0)},
			},
		},
		{
			&gw.UplinkRxInfo{
				Rssi: -74, Snr: 7, Channel: 2,
				Metadata: map[string]string{"aaa": "bbb"},
			}, []*decode.Point{
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "lora_snr", Value: float64(7)},
				{Attr: "channel", Value: int32(2)},
				{Attr: "aaa", Value: "bbb"},
			},
		},
		// Gateway UplinkRxInfo bad length.
		{
			nil, nil,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			res := ParseRXInfo(test.inp)
			t.Logf("res: %v", res)
			require.Equal(t, test.res, res)
		})
	}
}

func TestParseRXInfos(t *testing.T) {
	t.Parallel()

	gatewayID := random.String(16)

	now := time.Now().UTC().Add(-15 * time.Minute)
	tsNow := timestamppb.New(now)
	bad := time.Now().UTC().Add(time.Minute)
	tsBad := timestamppb.New(bad)

	// Gateway UplinkRxInfo slices, see ParseRXInfos() for format description.
	tests := []struct {
		inp       []*gw.UplinkRxInfo
		resTS     *timestamppb.Timestamp
		resPoints []*decode.Point
	}{
		// Gateway UplinkRxInfos.
		{
			[]*gw.UplinkRxInfo{
				{},
			}, nil, []*decode.Point{
				{Attr: "channel", Value: int32(0)},
			},
		},
		{
			[]*gw.UplinkRxInfo{
				{GatewayId: "aaa", GwTime: tsNow, Rssi: -80, Snr: 1},
				{
					GatewayId: gatewayID, GwTime: tsNow, Rssi: -74, Snr: 7,
					Metadata: map[string]string{"aaa": "bbb"},
				},
			}, tsNow, []*decode.Point{
				{Attr: "gateway_id", Value: gatewayID},
				{
					Attr:  "gateway_time",
					Value: strconv.FormatInt(now.Unix(), 10),
				},
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "lora_snr", Value: float64(7)},
				{Attr: "channel", Value: int32(0)},
				{Attr: "aaa", Value: "bbb"},
			},
		},
		{
			[]*gw.UplinkRxInfo{
				{GatewayId: gatewayID, GwTime: tsBad, Rssi: -74, Snr: 7},
			}, nil, []*decode.Point{
				{Attr: "gateway_id", Value: gatewayID},
				{
					Attr:  "gateway_time",
					Value: strconv.FormatInt(bad.Unix(), 10),
				},
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "lora_snr", Value: float64(7)},
				{Attr: "channel", Value: int32(0)},
			},
		},
		// Gateway UplinkRxInfo bad length.
		{
			nil, nil, nil,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			ts, res := ParseRXInfos(test.inp)
			t.Logf("ts, res: %v, %#v", ts, res)
			require.Equal(t, test.resPoints, res)
			if test.resTS == nil {
				require.WithinDuration(t, time.Now(), ts.AsTime(),
					2*time.Second)
			} else {
				require.EqualExportedValues(t, test.resTS, ts)
			}
		})
	}
}
