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
	if upMsg.GetTxInfo() != nil {
		if upMsg.GetTxInfo().GetFrequency() != 0 {
			//nolint:gosec // Safe conversion for limited values.
			msgs = append(msgs, &decode.Point{
				Attr:  "frequency",
				Value: int32(upMsg.GetTxInfo().GetFrequency()),
			})
		}

		if upMsg.GetTxInfo().GetModulation() != nil {
			mod := upMsg.GetTxInfo().GetModulation().GetLora()
			if mod.GetSpreadingFactor() != 0 {
				//nolint:gosec // Safe conversion for limited values.
				msgs = append(msgs, &decode.Point{
					Attr: "sf", Value: int32(mod.GetSpreadingFactor()),
				})
			}
		}
	}

	// Parse UplinkRXInfo.
	msgs = append(msgs, chirpstack.ParseRXInfo(upMsg.GetRxInfo())...)

	return msgs, nil
}
