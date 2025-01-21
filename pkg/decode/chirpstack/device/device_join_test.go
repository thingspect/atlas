//go:build !integration

package device

import (
	"fmt"
	"testing"

	"github.com/chirpstack/chirpstack/api/go/v4/common"
	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

func TestDeviceJoin(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)
	devAddr := random.String(16)

	// Device Join payloads, see deviceJoin() for format description.
	tests := []struct {
		inp       *integration.JoinEvent
		resPoints []*decode.Point
		err       string
	}{
		// Device Join.
		{
			&integration.JoinEvent{}, []*decode.Point{
				{Attr: "raw_device", Value: `{}`},
				{Attr: "join", Value: true},
			}, "",
		},
		{
			&integration.JoinEvent{
				DeviceInfo: &integration.DeviceInfo{
					DeviceProfileName: "1.0.2", DevEui: uniqID,
					DeviceClassEnabled: common.DeviceClass_CLASS_C,
				}, DevAddr: devAddr, RegionConfigId: "us915_0",
			}, []*decode.Point{
				{Attr: "raw_device", Value: fmt.Sprintf(`{"deviceInfo":`+
					`{"deviceProfileName":"1.0.2","devEui":"%s",`+
					`"deviceClassEnabled":"CLASS_C"},"devAddr":"%s",`+
					`"regionConfigId":"us915_0"}`, uniqID,
					devAddr)},
				{Attr: "join", Value: true},
				{Attr: "devaddr", Value: devAddr},
				{Attr: "region_config_id", Value: "us915_0"},
				{Attr: "id", Value: uniqID},
				{Attr: "lora_profile", Value: "1.0.2"},
				{Attr: "class", Value: "CLASS_C"},
			}, "",
		},
		// Device Join bad length.
		{
			nil, nil, "cannot parse invalid wire-format data",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp := []byte("aaa")
			if test.inp != nil {
				var err error
				bInp, err = proto.Marshal(test.inp)
				require.NoError(t, err)
			}

			res, err := deviceJoin(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, test.resPoints, res)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
