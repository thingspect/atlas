package gateway

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// gatewayConn parses a gateway connection state payload from a []byte according
// to the spec.
func gatewayConn(body []byte) ([]*decode.Point, error) {
	connMsg := &gw.ConnState{}
	if err := proto.Unmarshal(body, connMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(connMsg), " ", "")}}

	// Parse ConnState.
	msgs = append(msgs, &decode.Point{
		Attr: "conn", Value: connMsg.GetState().String(),
	})

	return msgs, nil
}
