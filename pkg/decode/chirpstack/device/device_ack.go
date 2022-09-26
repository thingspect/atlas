package device

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	ackOK      = "OK"
	ackTimeout = "DOWNLINK_TIMEOUT"
)

// deviceAck parses a device ACK payload from a []byte according to the spec.
func deviceAck(body []byte) ([]*decode.Point, error) {
	ackMsg := &integration.AckEvent{}
	if err := proto.Unmarshal(body, ackMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_device", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(ackMsg), " ", "")}}

	// Parse AckEvent. A false ack means it timed out.
	if ackMsg.Acknowledged {
		msgs = append(msgs, &decode.Point{Attr: "ack", Value: ackOK})
	} else {
		msgs = append(msgs, &decode.Point{Attr: "ack", Value: ackTimeout})
	}

	return msgs, nil
}
