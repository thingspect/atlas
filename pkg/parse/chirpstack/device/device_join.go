package device

import (
	"encoding/hex"
	"sort"
	"strconv"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/parse"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// deviceJoin parses a device Join payload from a []byte according to the spec.
func deviceJoin(body []byte) ([]*parse.Point, *timestamppb.Timestamp, error) {
	joinMsg := &as.JoinEvent{}
	if err := proto.Unmarshal(body, joinMsg); err != nil {
		return nil, nil, err
	}

	// Build raw device payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(joinMsg)
	if err != nil {
		return nil, nil, err
	}
	msgs := []*parse.Point{{Attr: "raw_device", Value: gw}}

	// Parse JoinEvent.
	msgs = append(msgs, &parse.Point{Attr: "join", Value: true})
	if len(joinMsg.DevEui) != 0 {
		msgs = append(msgs, &parse.Point{Attr: "id",
			Value: hex.EncodeToString(joinMsg.DevEui)})
	}
	if len(joinMsg.DevAddr) != 0 {
		msgs = append(msgs, &parse.Point{Attr: "devaddr",
			Value: hex.EncodeToString(joinMsg.DevAddr)})
	}

	// Parse UplinkRXInfos.
	joinTime := timestamppb.Now()
	if len(joinMsg.RxInfo) > 0 {
		// Sort joinMsg.RxInfo by strongest RSSI.
		sort.Slice(joinMsg.RxInfo, func(i, j int) bool {
			return joinMsg.RxInfo[i].Rssi > joinMsg.RxInfo[j].Rssi
		})

		if len(joinMsg.RxInfo[0].GatewayId) != 0 {
			msgs = append(msgs, &parse.Point{Attr: "gateway_id",
				Value: hex.EncodeToString(joinMsg.RxInfo[0].GatewayId)})
		}

		// Populate time channel if it is provided by the gateway. Use it as
		// joinTime if it is accurate.
		if joinMsg.RxInfo[0].Time != nil {
			msgs = append(msgs, &parse.Point{Attr: "time",
				Value: strconv.FormatInt(joinMsg.RxInfo[0].Time.Seconds, 10)})

			ts := joinMsg.RxInfo[0].Time.AsTime()
			if ts.Before(joinTime.AsTime()) &&
				time.Since(ts) < parse.ValidWindow {
				joinTime = joinMsg.RxInfo[0].Time
			}
		}

		if joinMsg.RxInfo[0].Rssi != 0 {
			msgs = append(msgs, &parse.Point{Attr: "lora_rssi",
				Value: int(joinMsg.RxInfo[0].Rssi)})
		}
		if joinMsg.RxInfo[0].LoraSnr != 0 {
			msgs = append(msgs, &parse.Point{Attr: "snr",
				Value: joinMsg.RxInfo[0].LoraSnr})
		}
	}

	// Parse UplinkTXInfo.
	if joinMsg.TxInfo != nil && joinMsg.TxInfo.Frequency != 0 {
		msgs = append(msgs, &parse.Point{Attr: "frequency",
			Value: int32(joinMsg.TxInfo.Frequency)})
	}

	// Parse JoinEvent data rate.
	msgs = append(msgs, &parse.Point{Attr: "data_rate",
		Value: int32(joinMsg.Dr)})

	return msgs, joinTime, nil
}
