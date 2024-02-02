//go:build !integration

package tektelic

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestChanMotion(t *testing.T) {
	t.Parallel()

	// Motion data channels, see chanMotion() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		rem []byte
		err string
	}{
		// Motion.
		{"0a0000", []*decode.Point{
			{Attr: "motion", Value: false},
		}, []byte{}, ""},
		{"0a00ff", []*decode.Point{
			{Attr: "motion", Value: true},
		}, []byte{}, ""},
		{"0a00000a00ff", []*decode.Point{
			{Attr: "motion", Value: false},
		}, []byte{0x0a, 0x00, 0xff}, ""},
		// Motion bad length.
		{"0a", nil, nil, "chanMotion format bad length: 0a"},
		// Motion bad identifier.
		{"0b0000", nil, nil, "chanMotion format bad identifier: 0b0000"},
		// Digital bad identifier.
		{"0a0100", nil, nil, "typeDigital format bad identifier: 0a0100"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := chanMotion(bInp)
			t.Logf("res, rem, err: %#v, %v, %v", res, rem, err)
			require.Equal(t, lTest.res, res)
			require.Equal(t, lTest.rem, rem)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}

func TestChanTempC(t *testing.T) {
	t.Parallel()

	// Temperature data channels, see chanTempC() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		rem []byte
		err string
	}{
		// Temperature.
		{"0367000a", []*decode.Point{
			{Attr: "temp_c", Value: 1.0},
			{Attr: "temp_f", Value: 33.8},
		}, []byte{}, ""},
		{"036700ca", []*decode.Point{
			{Attr: "temp_c", Value: 20.2},
			{Attr: "temp_f", Value: 68.4},
		}, []byte{}, ""},
		{"0367fff0", []*decode.Point{
			{Attr: "temp_c", Value: -1.6},
			{Attr: "temp_f", Value: 29.1},
		}, []byte{}, ""},
		{"036700c4", []*decode.Point{
			{Attr: "temp_c", Value: 19.6},
			{Attr: "temp_f", Value: 67.3},
		}, []byte{}, ""},
		{"036700c404687f", []*decode.Point{
			{Attr: "temp_c", Value: 19.6},
			{Attr: "temp_f", Value: 67.3},
		}, []byte{0x04, 0x68, 0x7f}, ""},
		// Temperature bad length.
		{"03", nil, nil, "chanTempC format bad length: 03"},
		// Temperature bad identifier.
		{"046700c4", nil, nil, "chanTempC format bad identifier: 046700c4"},
		// Temperature data type bad identifier.
		{"036800c4", nil, nil, "typeTempC format bad identifier: 036800c4"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := chanTempC(bInp)
			t.Logf("res, rem, err: %#v, %v, %v", res, rem, err)
			require.Equal(t, lTest.res, res)
			require.Equal(t, lTest.rem, rem)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}

func TestChanHumidity(t *testing.T) {
	t.Parallel()

	// Humidity data channels, see chanHumidity() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		rem []byte
		err string
	}{
		// Humidity.
		{"046814", []*decode.Point{
			{Attr: "humidity_pct", Value: 10.0},
		}, []byte{}, ""},
		{"04687f", []*decode.Point{
			{Attr: "humidity_pct", Value: 63.5},
		}, []byte{}, ""},
		{"04687f00ff0138", []*decode.Point{
			{Attr: "humidity_pct", Value: 63.5},
		}, []byte{0x00, 0xff, 0x01, 0x38}, ""},
		// Humidity bad length.
		{"04", nil, nil, "chanHumidity format bad length: 04"},
		// Humidity bad identifier.
		{"05687f", nil, nil, "chanHumidity format bad identifier: 05687f"},
		// Humidity data type bad identifier.
		{"04697f", nil, nil, "typeHumidity format bad identifier: 04697f"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := chanHumidity(bInp)
			t.Logf("res, rem, err: %#v, %v, %v", res, rem, err)
			require.Equal(t, lTest.res, res)
			require.Equal(t, lTest.rem, rem)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}

func TestChanBatteryV(t *testing.T) {
	t.Parallel()

	// Battery (V) data channels, see chanBatteryV() for format description.
	tests := []struct {
		inp string
		res []*decode.Point
		rem []byte
		err string
	}{
		// Battery (V).
		{"00ff012c", []*decode.Point{
			{Attr: "battery_v", Value: 3.0},
		}, []byte{}, ""},
		{"00ff0138", []*decode.Point{
			{Attr: "battery_v", Value: 3.12},
		}, []byte{}, ""},
		{"00ff013804687f", []*decode.Point{
			{Attr: "battery_v", Value: 3.12},
		}, []byte{0x04, 0x68, 0x7f}, ""},
		// Battery (V) bad length.
		{"00", nil, nil, "chanBatteryV format bad length: 00"},
		// Battery (V) bad identifier.
		{"01ff0138", nil, nil, "chanBatteryV format bad identifier: 01ff0138"},
		// Analog (V) bad identifier.
		{"00fe0138", nil, nil, "typeAnalogV format bad identifier: 00fe0138"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := chanBatteryV(bInp)
			t.Logf("res, rem, err: %#v, %v, %v", res, rem, err)
			require.Equal(t, lTest.res, res)
			require.Equal(t, lTest.rem, rem)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
