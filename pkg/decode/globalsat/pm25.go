package globalsat

import (
	"encoding/binary"

	"github.com/thingspect/atlas/pkg/decode"
)

const identPM25 = 0x03

// PM25 parses a PM2.5 payload from a []byte according to the spec. If a fatal
// error is encountered, it is returned along with any valid points.
func PM25(body []byte) ([]*decode.Point, error) {
	msgs, err := ls11x(body)
	if err != nil {
		return msgs, err
	}

	if body[0] != identPM25 {
		return msgs, decode.FormatErr("pm25", "bad identifier", body)
	}

	// Parse PM25.
	pm25 := int32(binary.BigEndian.Uint16(body[5:7]))
	msgs = append(msgs, &decode.Point{Attr: "pm25_ugm3", Value: pm25})

	return msgs, nil
}
