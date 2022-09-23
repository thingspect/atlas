package gateway

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/decode/chirpstack"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// gatewayUp parses a gateway Uplink payload from a []byte according to the
// spec.
func gatewayUp(body []byte) ([]*decode.Point, error) {
	upMsg := &gw.UplinkFrame{}
	if err := proto.Unmarshal(body, upMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(upMsg), " ", "")}}

	// Parse UplinkTXInfo.
	if upMsg.TxInfo != nil {
		if upMsg.TxInfo.Frequency != 0 {
			msgs = append(msgs, &decode.Point{
				Attr: "frequency", Value: int32(upMsg.TxInfo.Frequency),
			})
		}

		if upMsg.TxInfo.Modulation != nil {
			mod := upMsg.TxInfo.Modulation.GetLora()
			if mod != nil && mod.SpreadingFactor != 0 {
				msgs = append(msgs, &decode.Point{
					Attr: "sf", Value: int32(mod.SpreadingFactor),
				})
			}
		}
	}

	// Parse UplinkRXInfo.
	msgs = append(msgs, chirpstack.ParseRXInfo(upMsg.RxInfo)...)

	return msgs, nil
}
