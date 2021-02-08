package gateway

import (
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/parse"
)

// gatewayUp parses a gateway Uplink payload from a []byte according to the
// spec.
func gatewayUp(body []byte) ([]*parse.Point, error) {
	upMsg := &gw.UplinkFrame{}
	if err := proto.Unmarshal(body, upMsg); err != nil {
		return nil, err
	}

	// Build raw gateway payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(upMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*parse.Point{{Attr: "raw_gateway", Value: gw}}

	// Parse UplinkTXInfo.
	if upMsg.TxInfo != nil {
		if upMsg.TxInfo.Frequency != 0 {
			msgs = append(msgs, &parse.Point{Attr: "frequency",
				Value: int32(upMsg.TxInfo.Frequency)})
		}

		mod := upMsg.TxInfo.GetLoraModulationInfo()
		if mod != nil && mod.SpreadingFactor != 0 {
			msgs = append(msgs, &parse.Point{Attr: "spread_factor",
				Value: int32(mod.SpreadingFactor)})
		}
	}

	// Parse UplinkRXInfo.
	if upMsg.RxInfo != nil {
		if upMsg.RxInfo.Rssi != 0 {
			msgs = append(msgs, &parse.Point{Attr: "lora_rssi",
				Value: upMsg.RxInfo.Rssi})
		}

		if upMsg.RxInfo.LoraSnr != 0 {
			msgs = append(msgs, &parse.Point{Attr: "snr",
				Value: upMsg.RxInfo.LoraSnr})
		}

		if upMsg.RxInfo.Channel != 0 {
			msgs = append(msgs, &parse.Point{Attr: "channel",
				Value: int32(upMsg.RxInfo.Channel)})
		}
	}

	return msgs, nil
}
