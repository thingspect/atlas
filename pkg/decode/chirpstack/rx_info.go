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

// ParseRXInfo parses a gateway UplinkRxInfo payload according to the spec for
// gateway and device use.
func ParseRXInfo(rxInfo *gw.UplinkRxInfo) []*decode.Point {
	if rxInfo == nil {
		return nil
	}

	// Preallocate as large of a slice as exists in rxInfo.Metadata and channel.
	msgs := make([]*decode.Point, 0, len(rxInfo.GetMetadata())+1)

	if rxInfo.GetRssi() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "lora_rssi", Value: rxInfo.GetRssi(),
		})
	}
	if rxInfo.GetSnr() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "lora_snr", Value: float64(rxInfo.GetSnr()),
		})
	}
	msgs = append(msgs, &decode.Point{
		//nolint:gosec // Safe conversion for limited values.
		Attr: "channel", Value: int32(rxInfo.GetChannel()),
	})
	for k, v := range rxInfo.GetMetadata() {
		msgs = append(msgs, &decode.Point{Attr: k, Value: v})
	}

	return msgs
}

// ParseRXInfos parses a gateway UplinkRxInfo slice according to the spec for
// device use.
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
		return rxInfos[i].GetRssi() > rxInfos[j].GetRssi()
	})

	if rxInfos[0].GetGatewayId() != "" {
		msgs = append(msgs, &decode.Point{
			Attr: "gateway_id", Value: rxInfos[0].GetGatewayId(),
		})
	}

	// Populate time channel if it is provided by the gateway. Use it as msgTime
	// if it is accurate.
	if rxInfos[0].GetGwTime() != nil {
		msgs = append(msgs, &decode.Point{
			Attr:  "gateway_time",
			Value: strconv.FormatInt(rxInfos[0].GetGwTime().GetSeconds(), 10),
		})

		ts := rxInfos[0].GetGwTime().AsTime()
		if ts.Before(msgTime.AsTime()) && time.Since(ts) < decode.ValidWindow {
			msgTime = rxInfos[0].GetGwTime()
		}
	}

	msgs = append(msgs, ParseRXInfo(rxInfos[0])...)

	return msgTime, msgs
}
