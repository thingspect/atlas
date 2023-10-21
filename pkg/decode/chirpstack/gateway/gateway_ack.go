package gateway

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// gatewayAck parses a gateway Downlink ACK payload from a []byte according to
// the spec.
func gatewayAck(body []byte) ([]*decode.Point, error) {
	ackMsg := &gw.DownlinkTxAck{}
	if err := proto.Unmarshal(body, ackMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(ackMsg), " ", "")}}

	// Parse DownlinkTXAckItems.
	for _, item := range ackMsg.GetItems() {
		msgs = append(msgs, &decode.Point{
			Attr: "ack", Value: item.GetStatus().String(),
		})
	}

	return msgs, nil
}
