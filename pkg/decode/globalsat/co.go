package globalsat

import (
	"encoding/binary"

	"github.com/thingspect/atlas/pkg/decode"
)

const identCO = 0x02

// CO parses a CO payload from a []byte according to the spec. If a fatal
// error is encountered, it is returned along with any valid points.
func CO(body []byte) ([]*decode.Point, error) {
	msgs, err := ls11x(body)
	if err != nil {
		return msgs, err
	}

	if body[0] != identCO {
		return msgs, decode.FormatErr("co", "bad identifier", body)
	}

	// Parse CO.
	co := int32(binary.BigEndian.Uint16(body[5:7]))
	msgs = append(msgs, &decode.Point{Attr: "co_ppm", Value: co})

	return msgs, nil
}
