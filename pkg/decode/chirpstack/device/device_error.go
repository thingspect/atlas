package device

import (
	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
)

// deviceError parses a device Error payload from a []byte according to the
// spec.
func deviceError(body []byte) ([]*decode.Point, error) {
	errMsg := &as.ErrorEvent{}
	if err := proto.Unmarshal(body, errMsg); err != nil {
		return nil, err
	}

	// Build raw device payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(errMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*decode.Point{{Attr: "raw_device", Value: gw}}

	// Parse ErrorEvent.
	msgs = append(msgs, &decode.Point{Attr: "error_type",
		Value: errMsg.Type.String()})
	if errMsg.Error != "" {
		msgs = append(msgs, &decode.Point{Attr: "error", Value: errMsg.Error})
	}

	return msgs, nil
}
