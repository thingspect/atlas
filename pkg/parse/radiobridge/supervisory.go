package radiobridge

import (
	"encoding/binary"

	"github.com/thingspect/atlas/pkg/parse"
)

const identSupervisory = 0x01

// Supervisory error bitmap.
const (
	statusSuperTamper  = 0b00001000
	errorSuperReserved = 0b11100000
)

var errorSuper = []struct {
	flag byte
	code string
}{
	{0b00000001, "radio_comm"},
	{0b00000010, "battery_low"},
	{0b00000100, "last_downlink"},
	{0b00010000, "tamper_since_reset"},
}

// supervisory parses a Supervisory (status) payload from a []byte according to
// the spec. Supervisory payloads are used to confirm device operation and
// report error conditions.
func supervisory(body []byte) ([]*parse.Point, error) {
	const clearVolt = 0b00001111

	// Supervisory payload must be at least 5 bytes.
	if len(body) < 5 {
		return nil, parse.ErrFormat("supervisory", "bad length", body)
	}

	if body[1] != identSupervisory {
		return nil, parse.ErrFormat("supervisory", "bad identifier", body)
	}

	// Parse protocol and count.
	proto := int(body[0] >> 4)
	msgs := []*parse.Point{{Attr: "proto", Value: proto}}

	count := int(body[0] & clearProto)
	msgs = append(msgs, &parse.Point{Attr: "count", Value: count})

	// Parse error codes.
	errCodes := body[2]
	for _, err := range errorSuper {
		if errCodes&err.flag == err.flag {
			msgs = append(msgs, &parse.Point{Attr: "error", Value: err.code})
		}
	}

	if errCodes&errorSuperReserved > 0 {
		return msgs, parse.ErrFormat("supervisory", "bad error bitmap", body)
	}

	// Parse tamper status. Matching bit flag indicates closed/false state.
	msgs = append(msgs, &parse.Point{Attr: "tamper",
		Value: errCodes&statusSuperTamper != statusSuperTamper})

	// Parse battery level.
	vInt := body[4] >> 4
	vFract := body[4] & clearVolt
	batt := float64(vInt) + float64(vFract)/10
	msgs = append(msgs, &parse.Point{Attr: "battery", Value: batt})

	// Event count, 2 bytes.
	if len(body) >= 11 {
		msgs = append(msgs, &parse.Point{Attr: "total_count",
			Value: int(binary.BigEndian.Uint16(body[9:11]))})
	}

	// Check for any remaining bytes.
	if len(body) > 11 {
		return msgs, parse.ErrFormat("supervisory", "unused trailing bytes",
			body)
	}

	return msgs, nil
}
