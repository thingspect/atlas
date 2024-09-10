package globalsat

import (
	"encoding/binary"
	"math"

	"github.com/thingspect/atlas/pkg/decode"
)

// ls11x parses the common portion of an LS-11X payload from a []byte according
// to the spec.
func ls11x(body []byte) ([]*decode.Point, error) {
	// LS-11X payload must be 7 bytes.
	if len(body) != 7 {
		return nil, decode.ErrFormat("ls11x", "bad length", body)
	}

	// Parse temperature, rounded to one decimal digit.
	//nolint:gosec // Safe conversion for limited values.
	tempC := float64(int16(binary.BigEndian.Uint16(body[1:3]))) / 100
	msgs := []*decode.Point{{Attr: "temp_c", Value: math.Round(tempC*10) / 10}}
	msgs = append(msgs,
		&decode.Point{Attr: "temp_f", Value: decode.CToF(tempC)})

	// Parse humidity.
	hum := float64(binary.BigEndian.Uint16(body[3:5])) / 100
	if hum > 100 {
		return msgs, decode.ErrFormat("ls11x", "humidity outside allowed range",
			body)
	}

	msgs = append(msgs, &decode.Point{Attr: "humidity_pct", Value: hum})

	return msgs, nil
}
