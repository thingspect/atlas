// Package chirpstack provides helper functions for use by ChirpStack decoder
// function subpackages.
package chirpstack

import (
	"sort"
	"strconv"
	"time"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ParseRXInfo parses a gateway UplinkRxInfo payload according to the spec.
func ParseRXInfo(rxInfo *gw.UplinkRxInfo) []*decode.Point {
	if rxInfo == nil {
		return nil
	}

	// Preallocate at least as large of a slice as exists in rxInfo.Metadata.
	msgs := make([]*decode.Point, 0, len(rxInfo.Metadata))

	if rxInfo.Rssi != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "lora_rssi", Value: rxInfo.Rssi,
		})
	}
	if rxInfo.Snr != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "lora_snr", Value: float64(rxInfo.Snr),
		})
	}
	msgs = append(msgs, &decode.Point{
		Attr: "channel", Value: int32(rxInfo.Channel),
	})
	for k, v := range rxInfo.Metadata {
		msgs = append(msgs, &decode.Point{Attr: k, Value: v})
	}

	return msgs
}

// ParseRXInfos parses a gateway UplinkRxInfo slice according to the spec.
func ParseRXInfos(rxInfos []*gw.UplinkRxInfo) (
	*timestamppb.Timestamp, []*decode.Point,
) {
	msgTime := timestamppb.Now()

	if len(rxInfos) == 0 {
		return msgTime, nil
	}

	var msgs []*decode.Point

	// Sort rxInfos by strongest RSSI.
	sort.Slice(rxInfos, func(i, j int) bool {
		return rxInfos[i].Rssi > rxInfos[j].Rssi
	})

	if rxInfos[0].GatewayId != "" {
		msgs = append(msgs, &decode.Point{
			Attr: "gateway_id", Value: rxInfos[0].GatewayId,
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
