//go:build !integration

package device

import (
	"fmt"
	"testing"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/proto"
)

func TestDeviceStatus(t *testing.T) {
	t.Parallel()

	// Device Status payloads, see deviceStatus() for format description.
	tests := []struct {
		inp *integration.StatusEvent
		res []*decode.Point
		err string
	}{
		// Device Status.
		{&integration.StatusEvent{}, []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "ext_power", Value: false},
		}, ""},
		{&integration.StatusEvent{
			Margin: 7, BatteryLevelUnavailable: true,
		}, []*decode.Point{
			{Attr: "raw_device", Value: `{"margin":7,` +
				`"batteryLevelUnavailable":true}`},
			{Attr: "lora_snr_margin", Value: int32(7)},
			{Attr: "ext_power", Value: false},
			{Attr: "battery_unavail", Value: true},
		}, ""},
		{&integration.StatusEvent{Margin: 7, BatteryLevel: 99}, []*decode.Point{
			{Attr: "raw_device", Value: `{"margin":7,"batteryLevel":99}`},
			{Attr: "lora_snr_margin", Value: int32(7)},
			{Attr: "ext_power", Value: false},
			{Attr: "battery_pct", Value: float64(99)},
		}, ""},
		// Device Status bad length.
		{nil, nil, "cannot parse invalid wire-format data"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp := []byte("aaa")
			if lTest.inp != nil {
				var err error
				bInp, err = proto.Marshal(lTest.inp)
				require.NoError(t, err)
			}

			res, err := deviceStatus(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
