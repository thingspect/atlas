package radiobridge

import (
	"encoding/binary"
	"fmt"

	"github.com/thingspect/atlas/pkg/parse"
)

const identReset = 0x00

// reset parses an Reset (boot) payload from a []byte according to the spec.
// Reset payloads are used to confirm device settings on boot.
func reset(body []byte) ([]*parse.Point, error) {
	const fwFormat = 0b10000000

	// Reset payload must be at least 8 bytes.
	if len(body) < 8 {
		return nil, parse.ErrFormat("reset", "bad length", body)
	}

	if body[1] != identReset {
		return nil, parse.ErrFormat("reset", "bad identifier", body)
	}

	// Parse protocol, count, and hardware version.
	proto := int(body[0] >> 4)
	msgs := []*parse.Point{{Attr: "proto", Value: proto}}

	count := int(body[0] & clearProto)
	msgs = append(msgs, &parse.Point{Attr: "count", Value: count})
	msgs = append(msgs, &parse.Point{Attr: "hw_ver", Value: int(body[3])})

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
	msgs = append(msgs, &parse.Point{Attr: "ver", Value: firmware})

	// Check for any remaining bytes.
	if len(body) > 8 {
		return msgs, parse.ErrFormat("reset", "unused trailing bytes", body)
	}

	return msgs, nil
}
