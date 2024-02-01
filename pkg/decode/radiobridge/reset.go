package radiobridge

import (
	"encoding/binary"
	"fmt"

	"github.com/thingspect/atlas/pkg/decode"
)

const identReset = 0x00

// reset parses a Reset (boot) payload from a []byte according to the spec.
// Reset payloads are used to confirm device settings on boot.
func reset(body []byte) ([]*decode.Point, error) {
	const fwFormat = 0b10000000

	// Reset payload must be at least 8 bytes.
	if len(body) < 8 {
		return nil, decode.ErrFormat("reset", "bad length", body)
	}

	if body[1] != identReset {
		return nil, decode.ErrFormat("reset", "bad identifier", body)
	}

	// Parse protocol, count, and hardware version.
	proto := int32(body[0] >> 4)
	msgs := []*decode.Point{{Attr: "proto", Value: proto}}

	count := int32(body[0] & clearProto)
	msgs = append(msgs, &decode.Point{Attr: "count", Value: count})
	msgs = append(msgs, &decode.Point{Attr: "hw_ver", Value: int32(body[3])})

	// Parse firmware version.
	var firmware string
	// Format determined by MSBit.
	if body[4]&fwFormat == 0 {
		major := body[4] &^ fwFormat
		minor := body[5]
		firmware = fmt.Sprintf("%d.%d", major, minor)
	} else {
		major := body[4] &^ fwFormat >> 2
		minor := binary.BigEndian.Uint16(body[4:6]) >> 5 & 0b0000000000011111
		build := body[5] & 0b00011111
		firmware = fmt.Sprintf("%d.%d.%d", major, minor, build)
	}
	msgs = append(msgs, &decode.Point{Attr: "ver", Value: firmware})

	// Check for any remaining bytes.
	if len(body) > 8 {
		return msgs, decode.ErrFormat("reset", "unused trailing bytes", body)
	}

	return msgs, nil
}
