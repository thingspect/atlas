package gateway

import (
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
)

// gatewayAck parses a gateway Downlink ACK payload from a []byte according to
// the spec.
func gatewayAck(body []byte) ([]*decode.Point, error) {
	ackMsg := &gw.DownlinkTXAck{}
	if err := proto.Unmarshal(body, ackMsg); err != nil {
		return nil, err
	}

	// Build raw gateway payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(ackMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: gw}}

	// Parse DownlinkTXAckItems.
	for _, item := range ackMsg.Items {
		msgs = append(msgs, &decode.Point{Attr: "ack",
			Value: item.Status.String()})
	}

	return msgs, nil
}
