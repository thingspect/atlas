package device

import (
	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/parse"
)

// deviceTxAck parses a device TX ACK payload from a []byte according to the
// spec.
func deviceTxAck(body []byte) ([]*parse.Point, error) {
	txAckMsg := &as.TxAckEvent{}
	if err := proto.Unmarshal(body, txAckMsg); err != nil {
		return nil, err
	}

	// Build raw device payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(txAckMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*parse.Point{{Attr: "raw_device", Value: gw}}

	// Parse TxAckEvent.
	msgs = append(msgs, &parse.Point{Attr: "ack_gateway_tx", Value: true})

	return msgs, nil
}
