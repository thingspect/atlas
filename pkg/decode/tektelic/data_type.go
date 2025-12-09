package tektelic

import (
	"encoding/binary"

	"github.com/thingspect/atlas/pkg/decode"
)

const (
	identTypeDigital  = 0x00
	identTypeTempC    = 0x67
	identTypeHumidity = 0x68
	identTypeAnalogV  = 0xff
)

// typeDigital parses a Digital data type from a []byte according to the spec
// and returns the value, the remaining bytes, and an error value. Digital data
// types convey binary data.
func typeDigital(body []byte) (bool, []byte, error) {
	// Digital data type must be at least 3 bytes.
	if len(body) < 3 {
		return false, nil, decode.FormatErr("typeDigital", "bad length", body)
	}

	if body[1] != identTypeDigital {
		return false, nil, decode.FormatErr("typeDigital", "bad identifier",
			body)
	}

	// Parse presence.
	switch body[2] {
	case 0x00:
		return false, body[3:], nil
	case 0xff:
		return true, body[3:], nil
	default:
		return false, nil, decode.FormatErr("typeDigital", "unknown value",
			body)
	}
}

// typeTempC parses a Temperature data type from a []byte according to the spec
// and returns the value, the remaining bytes, and an error value.
func typeTempC(body []byte) (float64, []byte, error) {
	// Temperature data type must be at least 4 bytes.
	if len(body) < 4 {
		return 0, nil, decode.FormatErr("typeTempC", "bad length", body)
	}

	if body[1] != identTypeTempC {
		return 0, nil, decode.FormatErr("typeTempC", "bad identifier", body)
	}

	// Parse temperature.
	//nolint:gosec // Safe conversion for limited values.
	tempC := float64(int16(binary.BigEndian.Uint16(body[2:4]))) / 10

	return tempC, body[4:], nil
}

// typeHumidity parses a Humidity data type from a []byte according to the spec
// and returns the value, the remaining bytes, and an error value.
func typeHumidity(body []byte) (float64, []byte, error) {
	// Humidity data type must be at least 3 bytes.
	if len(body) < 3 {
		return 0, nil, decode.FormatErr("typeHumidity", "bad length", body)
	}

	if body[1] != identTypeHumidity {
		return 0, nil, decode.FormatErr("typeHumidity", "bad identifier", body)
	}

	// Parse hum.
	hum := float64(body[2]) / 2
	if hum > 100 {
		return 0, nil, decode.FormatErr("typeHumidity", "outside allowed range",
			body)
	}

	return hum, body[3:], nil
}

// typeAnalogV parses a Analog (V) data type from a []byte according to the spec
// and returns the value, the remaining bytes, and an error value.
func typeAnalogV(body []byte) (float64, []byte, error) {
	// Analog (V) data type must be at least 4 bytes.
	if len(body) < 4 {
		return 0, nil, decode.FormatErr("typeAnalogV", "bad length", body)
	}

	if body[1] != identTypeAnalogV {
		return 0, nil, decode.FormatErr("typeAnalogV", "bad identifier", body)
	}

	// Parse volt.
	volt := float64(binary.BigEndian.Uint16(body[2:4])) / 100

	return volt, body[4:], nil
}
