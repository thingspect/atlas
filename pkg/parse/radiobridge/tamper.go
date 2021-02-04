package radiobridge

import (
	"github.com/thingspect/atlas/pkg/parse"
)

const identTamper = 0x02

// tamper parses a Tamper payload from a []byte according to the spec. Tamper
// payloads are used to indicate whether the tamper switch is opened or closed.
func tamper(body []byte) ([]*parse.Point, error) {
	// Tamper payload must be 3 bytes.
	if len(body) != 3 {
		return nil, parse.ErrFormat("tamper", "bad length", body)
	}

	if body[1] != identTamper {
		return nil, parse.ErrFormat("tamper", "bad identifier", body)
	}

	// Parse protocol and count.
	proto := int(body[0] >> 4)
	msgs := []*parse.Point{{Attr: "proto", Value: proto}}

	count := int(body[0] & clearProto)
	msgs = append(msgs, &parse.Point{Attr: "count", Value: count})

	// Parse tamper status.
	msgs = append(msgs, &parse.Point{Attr: "tamper", Value: body[2] == 0x00})

	return msgs, nil
}
