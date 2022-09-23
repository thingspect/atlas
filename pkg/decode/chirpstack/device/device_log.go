package device

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// deviceLog parses a device Error payload from a []byte according to the
// spec.
func deviceLog(body []byte) ([]*decode.Point, error) {
	logMsg := &integration.LogEvent{}
	if err := proto.Unmarshal(body, logMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_device", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(logMsg), " ", "")}}

	// Parse ErrorEvent.
	msgs = append(msgs, &decode.Point{
		Attr: "log_level", Value: logMsg.Level.String(),
	})
	msgs = append(msgs, &decode.Point{
		Attr: "log_code", Value: logMsg.Code.String(),
	})
	if logMsg.Description != "" {
		msgs = append(msgs, &decode.Point{
			Attr: "log_desc", Value: logMsg.Description,
		})
	}

	return msgs, nil
}
