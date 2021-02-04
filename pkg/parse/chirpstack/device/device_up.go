package device

import (
	"encoding/hex"
	"sort"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/parse"
)

// deviceUp parses a device Uplink payload from a []byte according to the spec.
// Points, a time, and a data []byte are built from successful parse results. If
// a fatal error is encountered, it is returned along with any valid points.
func deviceUp(body []byte) ([]*parse.Point, time.Time, []byte, error) {
	upMsg := &as.UplinkEvent{}
	if err := proto.Unmarshal(body, upMsg); err != nil {
		return nil, time.Time{}, nil, err
	}

	// Build raw device payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(upMsg)
	if err != nil {
		return nil, time.Time{}, upMsg.Data, err
	}
	msgs := []*parse.Point{{Attr: "raw_device", Value: gw}}

	// Parse UplinkRXInfos.
	upTime := time.Now()
	if len(upMsg.RxInfo) > 0 {
		// Sort upMsg.RxInfo by strongest RSSI.
		sort.Slice(upMsg.RxInfo, func(i, j int) bool {
			return upMsg.RxInfo[i].Rssi > upMsg.RxInfo[j].Rssi
		})

		if len(upMsg.RxInfo[0].GatewayId) != 0 {
			msgs = append(msgs, &parse.Point{Attr: "gateway_id",
				Value: hex.EncodeToString(upMsg.RxInfo[0].GatewayId)})
		}

		// Populate time channel if it is provided by the gateway. Use it as
		// upTime if it is accurate.
		if upMsg.RxInfo[0].Time != nil {
			ts := upMsg.RxInfo[0].Time.AsTime()
			msgs = append(msgs, &parse.Point{Attr: "time",
				Value: int(ts.Unix())})
			if ts.Before(upTime) && time.Since(ts) < parse.ValidWindow {
				upTime = ts
			}
		}

		if upMsg.RxInfo[0].Rssi != 0 {
			msgs = append(msgs, &parse.Point{Attr: "rssi",
				Value: int(upMsg.RxInfo[0].Rssi)})
		}
		if upMsg.RxInfo[0].LoraSnr != 0 {
			msgs = append(msgs, &parse.Point{Attr: "snr",
				Value: upMsg.RxInfo[0].LoraSnr})
		}
	}

	// Parse UplinkTXInfo.
	if upMsg.TxInfo != nil && upMsg.TxInfo.Frequency != 0 {
		msgs = append(msgs, &parse.Point{Attr: "frequency",
			Value: int(upMsg.TxInfo.Frequency)})
	}

	// Parse UplinkEvent.
	msgs = append(msgs, &parse.Point{Attr: "adr", Value: upMsg.Adr})
	msgs = append(msgs, &parse.Point{Attr: "data_rate", Value: int(upMsg.Dr)})
	msgs = append(msgs, &parse.Point{Attr: "confirmed",
		Value: upMsg.ConfirmedUplink})

	return msgs, upTime.UTC(), upMsg.Data, nil
}
