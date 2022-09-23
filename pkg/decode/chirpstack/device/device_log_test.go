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

func TestDeviceError(t *testing.T) {
	t.Parallel()

	// Device Error payloads, see deviceError() for format description.
	tests := []struct {
		inp *integration.LogEvent
		res []*decode.Point
		err string
	}{
		// Device Error.
		{
			&integration.LogEvent{}, []*decode.Point{
				{Attr: "raw_device", Value: `{}`},
				{Attr: "log_level", Value: "INFO"},
				{Attr: "log_code", Value: "UNKNOWN"},
			}, "",
		},
		{
			&integration.LogEvent{
				Level: integration.LogLevel_WARNING,
				Code:  integration.LogCode_OTAA, Description: "OTAA_ERROR",
			},
			[]*decode.Point{
				{
					Attr: "raw_device",
					Value: `{"level":"WARNING","code":"OTAA","description":` +
						`"OTAA_ERROR"}`,
				},
				{Attr: "log_level", Value: "WARNING"},
				{Attr: "log_code", Value: "OTAA"},
				{Attr: "log_desc", Value: "OTAA_ERROR"},
			},
			"",
		},
		// Device Error bad length.
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

			res, err := deviceLog(bInp)
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
