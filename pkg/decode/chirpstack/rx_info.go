// Package chirpstack provides helper functions for use by ChirpStack decoder
// function subpackages.
package chirpstack

import (
	"encoding/hex"
	"sort"
	"strconv"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ParseRXInfo parses a gateway UplinkRXInfo payload according to the spec.
func ParseRXInfo(rxInfo *gw.UplinkRXInfo) []*decode.Point {
	if rxInfo == nil {
		return nil
	}

	var msgs []*decode.Point

	if rxInfo.Rssi != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "lora_rssi", Value: rxInfo.Rssi,
		})
	}
	if rxInfo.LoraSnr != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "snr", Value: rxInfo.LoraSnr,
		})
	}
	msgs = append(msgs, &decode.Point{
		Attr: "channel", Value: int32(rxInfo.Channel),
	})

	return msgs
}

// ParseRXInfos parses a gateway UplinkRXInfo slice according to the spec.
func ParseRXInfos(
	rxInfos []*gw.UplinkRXInfo,
) (*timestamppb.Timestamp, []*decode.Point) {
	msgTime := timestamppb.Now()

	if len(rxInfos) == 0 {
		return msgTime, nil
	}

	var msgs []*decode.Point

	// Sort rxInfos by strongest RSSI.
	sort.Slice(rxInfos, func(i, j int) bool {
		return rxInfos[i].Rssi > rxInfos[j].Rssi
	})

	if len(rxInfos[0].GatewayId) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "gateway_id", Value: hex.EncodeToString(rxInfos[0].GatewayId),
		})
	}

	// Populate time channel if it is provided by the gateway. Use it as msgTime
	// if it is accurate.
	if rxInfos[0].Time != nil {
		msgs = append(msgs, &decode.Point{
			Attr: "time", Value: strconv.FormatInt(rxInfos[0].Time.Seconds, 10),
		})

		ts := rxInfos[0].Time.AsTime()
		if ts.Before(msgTime.AsTime()) && time.Since(ts) < decode.ValidWindow {
			msgTime = rxInfos[0].Time
		}
	}

	msgs = append(msgs, ParseRXInfo(rxInfos[0])...)

	return msgTime, msgs
}
