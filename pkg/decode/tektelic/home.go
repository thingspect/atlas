package tektelic

import (
	"github.com/thingspect/atlas/pkg/decode"
)

// Home parses a Home Sensor payload from a []byte according to the spec. Points
// are built from successful parse results. If a fatal error is encountered, it
// is returned along with any valid points.
func Home(body []byte) ([]*decode.Point, error) {
	// Home Sensor payloads must be at least 3 bytes.
	if len(body) < 3 {
		return nil, decode.FormatErr("home", "bad length", body)
	}

	var msgs, lMsgs []*decode.Point
	var err error

	// Parse home.
	for len(body) >= 3 {
		switch body[0] {
		case identChanMotion:
			lMsgs, body, err = chanMotion(body)
		case identChanTempC:
			lMsgs, body, err = chanTempC(body)
		case identChanHumidity:
			lMsgs, body, err = chanHumidity(body)
		case identChanBatteryV:
			lMsgs, body, err = chanBatteryV(body)
		default:
			return msgs, decode.FormatErr("home", "bad identifier", body)
		}

		// Store valid points.
		msgs = append(msgs, lMsgs...)

		if err != nil {
			return msgs, err
		}
	}

	// Check for any remaining bytes.
	if len(body) > 0 {
		return msgs, decode.FormatErr("home", "unused trailing bytes", body)
	}

	return msgs, nil
}
