package device

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// deviceJoin parses a device Join payload from a []byte according to the spec.
func deviceJoin(body []byte) ([]*decode.Point, error) {
	joinMsg := &integration.JoinEvent{}
	if err := proto.Unmarshal(body, joinMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_device", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(joinMsg), " ", "")}}

	// Parse JoinEvent.
	msgs = append(msgs, &decode.Point{Attr: "join", Value: true})
	if joinMsg.DevAddr != "" {
		msgs = append(msgs, &decode.Point{
			Attr: "devaddr", Value: joinMsg.DevAddr,
		})
	}

	// Parse DeviceInfo.
	if joinMsg.DeviceInfo != nil {
		if joinMsg.DeviceInfo.DevEui != "" {
			msgs = append(msgs, &decode.Point{
				Attr: "id", Value: joinMsg.DeviceInfo.DevEui,
			})
		}
		if joinMsg.DeviceInfo.DeviceProfileName != "" {
			msgs = append(msgs, &decode.Point{
				Attr:  "lora_profile",
				Value: joinMsg.DeviceInfo.DeviceProfileName,
			})
		}
		msgs = append(msgs, &decode.Point{
			Attr:  "class",
			Value: joinMsg.DeviceInfo.DeviceClassEnabled.String(),
		})
	}

	return msgs, nil
}
