package gateway

import (
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
)

// gatewayConn parses a gateway connection state payload from a []byte according
// to the spec.
func gatewayConn(body []byte) ([]*decode.Point, error) {
	connMsg := &gw.ConnState{}
	if err := proto.Unmarshal(body, connMsg); err != nil {
		return nil, err
	}

	// Build raw gateway payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(connMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: gw}}

	// Parse ConnState.
	msgs = append(msgs, &decode.Point{
		Attr: "conn", Value: connMsg.State.String(),
	})

	return msgs, nil
}
