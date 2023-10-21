package device

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// deviceStatus parses a device Status payload from a []byte according to the
// spec.
func deviceStatus(body []byte) ([]*decode.Point, error) {
	statusMsg := &integration.StatusEvent{}
	if err := proto.Unmarshal(body, statusMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_device", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(statusMsg), " ", "")}}

	// Parse StatusEvent.
	if statusMsg.GetMargin() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "lora_snr_margin", Value: statusMsg.GetMargin(),
		})
	}
	msgs = append(msgs, &decode.Point{
		Attr: "ext_power", Value: statusMsg.GetExternalPowerSource(),
	})
	if statusMsg.GetBatteryLevelUnavailable() {
		msgs = append(msgs, &decode.Point{
			Attr: "battery_unavail", Value: statusMsg.GetBatteryLevelUnavailable(),
		})
	}
	if statusMsg.GetBatteryLevel() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "battery_pct", Value: float64(statusMsg.GetBatteryLevel()),
		})
	}

	return msgs, nil
}
