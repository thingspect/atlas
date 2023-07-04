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
				}, DevAddr: devAddr,
			}, []*decode.Point{
				{Attr: "raw_device", Value: fmt.Sprintf(`{"deviceInfo":`+
					`{"deviceProfileName":"1.0.2","devEui":"%s",`+
					`"deviceClassEnabled":"CLASS_C"},"devAddr":"%s"}`, uniqID,
					devAddr)},
				{Attr: "join", Value: true},
				{Attr: "devaddr", Value: devAddr},
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
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInp := []byte("aaa")
			if lTest.inp != nil {
				var err error
				bInp, err = proto.Marshal(lTest.inp)
				require.NoError(t, err)
			}

			res, err := deviceJoin(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.resPoints, res)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
