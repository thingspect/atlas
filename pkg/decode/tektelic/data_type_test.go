//go:build !integration

package tektelic

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeDigital(t *testing.T) {
	t.Parallel()

	// Digital data types, see typeDigital() for format description.
	tests := []struct {
		inp string
		res bool
		rem []byte
		err string
	}{
		// Digital.
		{"0a0000", false, []byte{}, ""},
		{"0a00ff", true, []byte{}, ""},
		{"0a00000a00ff", false, []byte{0x0a, 0x00, 0xff}, ""},
		// Digital bad length.
		{"0a", false, nil, "typeDigital format bad length: 0a"},
		// Digital bad identifier.
		{"0a0100", false, nil, "typeDigital format bad identifier: 0a0100"},
		// Digital unknown value.
		{"0a00fe", false, nil, "typeDigital format unknown value: 0a00fe"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := typeDigital(bInp)
			t.Logf("res, rem, err: %v, %v, %v", res, rem, err)
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

func TestTypeTempC(t *testing.T) {
	t.Parallel()

	epsilon := math.Nextafter(1, 2) - 1

	// Temperature data types, see typeTempC() for format description.
	tests := []struct {
		inp string
		res float64
		rem []byte
		err string
	}{
		// Temperature.
		{"0367000a", 1, []byte{}, ""},
		{"036700ca", 20.2, []byte{}, ""},
		{"0367fff0", -1.6, []byte{}, ""},
		{"036700c4", 19.6, []byte{}, ""},
		{"036700c404687f", 19.6, []byte{0x04, 0x68, 0x7f}, ""},
		// Temperature bad length.
		{"03", 0, nil, "typeTempC format bad length: 03"},
		// Temperature bad identifier.
		{"036800c4", 0, nil, "typeTempC format bad identifier: 036800c4"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := typeTempC(bInp)
			t.Logf("res, rem, err: %v, %v, %v", res, rem, err)
			require.InDelta(t, lTest.res, res, epsilon)
			require.Equal(t, lTest.rem, rem)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}

func TestTypeHumidity(t *testing.T) {
	t.Parallel()

	epsilon := math.Nextafter(1, 2) - 1

	// Humidity data types, see typeHumidity() for format description.
	tests := []struct {
		inp string
		res float64
		rem []byte
		err string
	}{
		// Humidity.
		{"046814", 10, []byte{}, ""},
		{"04687f", 63.5, []byte{}, ""},
		{"04687f00ff0138", 63.5, []byte{0x00, 0xff, 0x01, 0x38}, ""},
		// Humidity bad length.
		{"04", 0, nil, "typeHumidity format bad length: 04"},
		// Humidity bad identifier.
		{"04697f", 0, nil, "typeHumidity format bad identifier: 04697f"},
		// Humidity outside allowed range.
		{"0468f0", 0, nil, "typeHumidity format outside allowed range: 0468f0"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := typeHumidity(bInp)
			t.Logf("res, rem, err: %v, %v, %v", res, rem, err)
			require.InDelta(t, lTest.res, res, epsilon)
			require.Equal(t, lTest.rem, rem)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}

func TestTypeAnalogV(t *testing.T) {
	t.Parallel()

	epsilon := math.Nextafter(1, 2) - 1

	// Analog (V) data types, see typeAnalogV() for format description.
	tests := []struct {
		inp string
		res float64
		rem []byte
		err string
	}{
		// Analog (V).
		{"00ff012c", 3.0, []byte{}, ""},
		{"00ff0138", 3.12, []byte{}, ""},
		{"00ff013804687f", 3.12, []byte{0x04, 0x68, 0x7f}, ""},
		// Analog (V) bad length.
		{"00", 0, nil, "typeAnalogV format bad length: 00"},
		// Analog (V) bad identifier.
		{"00fe0138", 0, nil, "typeAnalogV format bad identifier: 00fe0138"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp, err := hex.DecodeString(lTest.inp)
			require.NoError(t, err)

			res, rem, err := typeAnalogV(bInp)
			t.Logf("res, rem, err: %v, %v, %v", res, rem, err)
			require.InDelta(t, lTest.res, res, epsilon)
			require.Equal(t, lTest.rem, rem)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
