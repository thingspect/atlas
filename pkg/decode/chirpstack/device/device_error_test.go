// +build !integration

package device

import (
	"fmt"
	"testing"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestDeviceError(t *testing.T) {
	t.Parallel()

	// Device Error payloads, see deviceError() for format description.
	tests := []struct {
		inp *as.ErrorEvent
		res []*decode.Point
		err string
	}{
		// Device Error.
		{&as.ErrorEvent{}, []*decode.Point{
			{Attr: "raw_device", Value: `{}`},
			{Attr: "error_type", Value: "UNKNOWN"},
		}, ""},
		{&as.ErrorEvent{Type: as.ErrorType_OTAA,
			Error: "OTAA_ERROR"}, []*decode.Point{
			{Attr: "raw_device", Value: `{"type":"OTAA","error":"OTAA_ERROR"}`},
			{Attr: "error_type", Value: "OTAA"},
			{Attr: "error", Value: "OTAA_ERROR"},
		}, ""},
		// Device Error bad length.
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

			res, err := deviceError(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err != "" {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}
