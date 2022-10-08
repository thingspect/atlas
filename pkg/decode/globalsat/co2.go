package globalsat

import (
	"encoding/binary"

	"github.com/thingspect/atlas/pkg/decode"
)

const identCO2 = 0x01

// CO2 parses a CO2 payload from a []byte according to the spec. If a fatal
// error is encountered, it is returned along with any valid points.
func CO2(body []byte) ([]*decode.Point, error) {
	msgs, err := ls11x(body)
	if err != nil {
		return msgs, err
	}

	if body[0] != identCO2 {
		return msgs, decode.ErrFormat("co2", "bad identifier", body)
	}

	// Parse CO2.
	co2 := int32(binary.BigEndian.Uint16(body[5:7]))
	msgs = append(msgs, &decode.Point{Attr: "co2_ppm", Value: co2})

	return msgs, nil
}
