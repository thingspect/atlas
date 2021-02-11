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
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// deviceUp parses a device Uplink payload from a []byte according to the spec.
// Points, a timestamp, and a data []byte are built from successful parse
// results. If a fatal error is encountered, it is returned along with any valid
// points.
func deviceUp(body []byte) ([]*decode.Point, *timestamppb.Timestamp, []byte,
	error) {
	upMsg := &as.UplinkEvent{}
	if err := proto.Unmarshal(body, upMsg); err != nil {
		return nil, nil, nil, err
	}

	// Build raw device and data payloads for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(upMsg)
	if err != nil {
		return nil, nil, upMsg.Data, err
	}
	msgs := []*decode.Point{{Attr: "raw_device", Value: gw}}

	if upMsg.Data != nil {
		msgs = append(msgs, &decode.Point{Attr: "raw_data",
			Value: hex.EncodeToString(upMsg.Data)})
	}

	// Parse UplinkRXInfos.
	upTime := timestamppb.Now()
	if len(upMsg.RxInfo) > 0 {
		// Sort upMsg.RxInfo by strongest RSSI.
		sort.Slice(upMsg.RxInfo, func(i, j int) bool {
			return upMsg.RxInfo[i].Rssi > upMsg.RxInfo[j].Rssi
		})

		if len(upMsg.RxInfo[0].GatewayId) != 0 {
			msgs = append(msgs, &decode.Point{Attr: "gateway_id",
				Value: hex.EncodeToString(upMsg.RxInfo[0].GatewayId)})
		}

		// Populate time channel if it is provided by the gateway. Use it as
		// upTime if it is accurate.
		if upMsg.RxInfo[0].Time != nil {
			msgs = append(msgs, &decode.Point{Attr: "time",
				Value: strconv.FormatInt(upMsg.RxInfo[0].Time.Seconds, 10)})

			ts := upMsg.RxInfo[0].Time.AsTime()
			if ts.Before(upTime.AsTime()) &&
				time.Since(ts) < decode.ValidWindow {
				upTime = upMsg.RxInfo[0].Time
			}
		}

		if upMsg.RxInfo[0].Rssi != 0 {
			msgs = append(msgs, &decode.Point{Attr: "lora_rssi",
				Value: int(upMsg.RxInfo[0].Rssi)})
		}
		if upMsg.RxInfo[0].LoraSnr != 0 {
			msgs = append(msgs, &decode.Point{Attr: "snr",
				Value: upMsg.RxInfo[0].LoraSnr})
		}
		msgs = append(msgs, &decode.Point{Attr: "channel",
			Value: int32(upMsg.RxInfo[0].Channel)})
	}

	// Parse UplinkTXInfo.
	if upMsg.TxInfo != nil && upMsg.TxInfo.Frequency != 0 {
		msgs = append(msgs, &decode.Point{Attr: "frequency",
			Value: int32(upMsg.TxInfo.Frequency)})
	}

	// Parse UplinkEvent.
	msgs = append(msgs, &decode.Point{Attr: "adr", Value: upMsg.Adr})
	msgs = append(msgs, &decode.Point{Attr: "data_rate", Value: int32(upMsg.Dr)})
	msgs = append(msgs, &decode.Point{Attr: "confirmed",
		Value: upMsg.ConfirmedUplink})

	return msgs, upTime, upMsg.Data, nil
}
