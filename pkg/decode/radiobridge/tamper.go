package radiobridge

import (
	"github.com/thingspect/atlas/pkg/decode"
)

const identTamper = 0x02

// tamper parses a Tamper payload from a []byte according to the spec. Tamper
// payloads are used to indicate whether the tamper switch is opened or closed.
func tamper(body []byte) ([]*decode.Point, error) {
	// Tamper payload must be 3 bytes.
	if len(body) != 3 {
		return nil, decode.FormatErr("tamper", "bad length", body)
	}

	if body[1] != identTamper {
		return nil, decode.FormatErr("tamper", "bad identifier", body)
	}

	msgs := make([]*decode.Point, 0, 3)

	// Parse protocol and count.
	proto := int32(body[0] >> 4)
	msgs = append(msgs, &decode.Point{Attr: "proto", Value: proto})

	count := int32(body[0] & clearProto)
	msgs = append(msgs, &decode.Point{Attr: "count", Value: count})

	// Parse tamper status.
	msgs = append(msgs, &decode.Point{Attr: "tamper", Value: body[2] == 0x00})

	return msgs, nil
}
