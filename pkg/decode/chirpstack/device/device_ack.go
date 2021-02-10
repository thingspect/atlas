package device

import (
	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
)

const ackOK = "OK"
const ackTimeout = "DOWNLINK_TIMEOUT"

// deviceAck parses a device ACK payload from a []byte according to the spec.
func deviceAck(body []byte) ([]*decode.Point, error) {
	ackMsg := &as.AckEvent{}
	if err := proto.Unmarshal(body, ackMsg); err != nil {
		return nil, err
	}

	// Build raw device payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(ackMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*decode.Point{{Attr: "raw_device", Value: gw}}

	// Parse AckEvent. A false ack means it timed out.
	if ackMsg.Acknowledged {
		msgs = append(msgs, &decode.Point{Attr: "ack", Value: ackOK})
	} else {
		msgs = append(msgs, &decode.Point{Attr: "ack", Value: ackTimeout})
	}

	return msgs, nil
}
