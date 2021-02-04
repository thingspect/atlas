package radiobridge

import (
	"github.com/thingspect/atlas/pkg/parse"
)

const identDoor = 0x03

// Door parses a Door payload from a []byte according to the spec. Door payloads
// are used to indicate open and close events. Points are built from successful
// parse results. If a fatal error is encountered, it is returned along with any
// valid points.
func Door(body []byte) ([]*parse.Point, error) {
	// Door and children payloads must be at least 3 bytes.
	if len(body) < 3 {
		return nil, parse.ErrFormat("door", "bad length", body)
	}

	switch body[1] {
	case identDoor:
		break
	case identReset:
		return reset(body)
	case identSupervisory:
		return supervisory(body)
	case identTamper:
		return tamper(body)
	case identLinkQuality:
		return linkQuality(body)
	default:
		return nil, parse.ErrFormat("door", "bad identifier", body)
	}

	// Parse count.
	count := int(body[0] & clearProto)
	msgs := []*parse.Point{{Attr: "count", Value: count}}

	// Parse open status.
	msgs = append(msgs, &parse.Point{Attr: "open", Value: body[2] == 0x01})

	// Check for any remaining bytes.
	if len(body) > 3 {
		return msgs, parse.ErrFormat("door", "unused trailing bytes", body)
	}

	return msgs, nil
}
