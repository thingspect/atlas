package tektelic

import (
	"github.com/thingspect/atlas/pkg/decode"
)

const (
	identChanMotion   = 0x0a
	identChanTempC    = 0x03
	identChanHumidity = 0x04
	identChanBatteryV = 0x00
)

// chanMotion parses a Motion data channel from a []byte according to the spec
// and returns the points, the remaining bytes, and an error value.
func chanMotion(body []byte) ([]*decode.Point, []byte, error) {
	// Motion data channel must be at least 3 bytes.
	if len(body) < 3 {
		return nil, nil, decode.FormatErr("chanMotion", "bad length", body)
	}

	if body[0] != identChanMotion {
		return nil, nil, decode.FormatErr("chanMotion", "bad identifier", body)
	}

	// Parse motion.
	motion, rem, err := typeDigital(body)
	if err != nil {
		return nil, nil, err
	}

	return []*decode.Point{{Attr: "motion", Value: motion}}, rem, nil
}

// chanTempC parses a Temperature data channel from a []byte according to the
// spec and returns the points, the remaining bytes, and an error value.
func chanTempC(body []byte) ([]*decode.Point, []byte, error) {
	// Temperature data channel must be at least 4 bytes.
	if len(body) < 4 {
		return nil, nil, decode.FormatErr("chanTempC", "bad length", body)
	}

	if body[0] != identChanTempC {
		return nil, nil, decode.FormatErr("chanTempC", "bad identifier", body)
	}

	// Parse temperature.
	tempC, rem, err := typeTempC(body)
	if err != nil {
		return nil, nil, err
	}

	return []*decode.Point{
		{Attr: "temp_c", Value: tempC},
		{Attr: "temp_f", Value: decode.CToF(tempC)},
	}, rem, nil
}

// chanHumidity parses a Humidity data channel from a []byte according to the
// spec and returns the points, the remaining bytes, and an error value.
func chanHumidity(body []byte) ([]*decode.Point, []byte, error) {
	// Humidity data channel must be at least 3 bytes.
	if len(body) < 3 {
		return nil, nil, decode.FormatErr("chanHumidity", "bad length", body)
	}

	if body[0] != identChanHumidity {
		return nil, nil, decode.FormatErr("chanHumidity", "bad identifier",
			body)
	}

	// Parse hum.
	hum, rem, err := typeHumidity(body)
	if err != nil {
		return nil, nil, err
	}

	return []*decode.Point{{Attr: "humidity_pct", Value: hum}}, rem, nil
}

// chanBatteryV parses a Battery (V) data channel from a []byte according to the
// spec and returns the points, the remaining bytes, and an error value.
func chanBatteryV(body []byte) ([]*decode.Point, []byte, error) {
	// Battery (V) data channel must be at least 4 bytes.
	if len(body) < 4 {
		return nil, nil, decode.FormatErr("chanBatteryV", "bad length", body)
	}

	if body[0] != identChanBatteryV {
		return nil, nil, decode.FormatErr("chanBatteryV", "bad identifier",
			body)
	}

	// Parse volt.
	volt, rem, err := typeAnalogV(body)
	if err != nil {
		return nil, nil, err
	}

	return []*decode.Point{{Attr: "battery_v", Value: volt}}, rem, nil
}
