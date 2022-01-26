//go:build !integration

package chirpstack

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseRXInfo(t *testing.T) {
	t.Parallel()

	// Gateway UplinkRXInfo payloads, see ParseRXInfo() for format description.
	tests := []struct {
		inp *gw.UplinkRXInfo
		res []*decode.Point
	}{
		// Gateway UplinkRXInfo.
		{
			&gw.UplinkRXInfo{}, []*decode.Point{
				{Attr: "channel", Value: int32(0)},
			},
		},
		{
			&gw.UplinkRXInfo{
				Rssi: -74, LoraSnr: 7.8, Channel: 2,
			}, []*decode.Point{
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "snr", Value: 7.8},
				{Attr: "channel", Value: int32(2)},
			},
		},
		// Gateway UplinkRXInfo bad length.
		{
			nil, nil,
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := ParseRXInfo(lTest.inp)
			t.Logf("res: %v", res)
			require.Equal(t, lTest.res, res)
		})
	}
}

func TestParseRXInfos(t *testing.T) {
	t.Parallel()

	gatewayID := random.String(16)
	bGatewayID, err := hex.DecodeString(gatewayID)
	require.NoError(t, err)
	t.Logf("bGatewayID: %x", bGatewayID)

	b64GatewayID := base64.StdEncoding.EncodeToString(bGatewayID)
	t.Logf("b64GatewayID: %v", b64GatewayID)

	now := time.Now().UTC().Add(-15 * time.Minute)
	tsNow := timestamppb.New(now)
	bad := time.Now().UTC().Add(time.Minute)
	tsBad := timestamppb.New(bad)

	// Gateway UplinkRXInfo slices, see ParseRXInfos() for format description.
	tests := []struct {
		inp       []*gw.UplinkRXInfo
		resTS     *timestamppb.Timestamp
		resPoints []*decode.Point
	}{
		// Gateway UplinkRXInfos.
		{
			[]*gw.UplinkRXInfo{
				{},
			}, nil, []*decode.Point{
				{Attr: "channel", Value: int32(0)},
			},
		},
		{
			[]*gw.UplinkRXInfo{
				{GatewayId: []byte("aaa"), Time: tsNow, Rssi: -80, LoraSnr: 1},
				{GatewayId: bGatewayID, Time: tsNow, Rssi: -74, LoraSnr: 7.8},
			}, tsNow, []*decode.Point{
				{Attr: "gateway_id", Value: gatewayID},
				{Attr: "time", Value: strconv.FormatInt(now.Unix(), 10)},
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "snr", Value: 7.8},
				{Attr: "channel", Value: int32(0)},
			},
		},
		{
			[]*gw.UplinkRXInfo{
				{GatewayId: bGatewayID, Time: tsBad, Rssi: -74, LoraSnr: 7.8},
			}, nil, []*decode.Point{
				{Attr: "gateway_id", Value: gatewayID},
				{Attr: "time", Value: strconv.FormatInt(bad.Unix(), 10)},
				{Attr: "lora_rssi", Value: int32(-74)},
				{Attr: "snr", Value: 7.8},
				{Attr: "channel", Value: int32(0)},
			},
		},
		// Gateway UplinkRXInfo bad length.
		{
			nil, nil, nil,
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			ts, res := ParseRXInfos(lTest.inp)
			t.Logf("ts, res: %v, %#v", ts, res)
			require.Equal(t, lTest.resPoints, res)
			if lTest.resTS == nil {
				require.WithinDuration(t, time.Now(), ts.AsTime(),
					2*time.Second)
			} else if !proto.Equal(lTest.resTS, ts) {
				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.resTS, ts)
			}
		})
	}
}
