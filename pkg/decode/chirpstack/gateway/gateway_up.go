package gateway

import (
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	//nolint:staticcheck // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	//nolint:staticcheck // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/decode/chirpstack"
)

// gatewayUp parses a gateway Uplink payload from a []byte according to the
// spec.
func gatewayUp(body []byte) ([]*decode.Point, error) {
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
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: gw}}

	// Parse UplinkTXInfo.
	if upMsg.TxInfo != nil {
		if upMsg.TxInfo.Frequency != 0 {
			msgs = append(msgs, &decode.Point{
				Attr: "frequency", Value: int32(upMsg.TxInfo.Frequency),
			})
		}

		mod := upMsg.TxInfo.GetLoraModulationInfo()
		if mod != nil && mod.SpreadingFactor != 0 {
			msgs = append(msgs, &decode.Point{
				Attr: "sf", Value: int32(mod.SpreadingFactor),
			})
		}
	}

	// Parse UplinkRXInfo.
	msgs = append(msgs, chirpstack.ParseRXInfo(upMsg.RxInfo)...)

	return msgs, nil
}
