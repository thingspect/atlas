//go:build !integration

package gateway

import (
	"fmt"
	"testing"

	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	//nolint:staticcheck // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
)

func TestGatewayConn(t *testing.T) {
	t.Parallel()

	// Gateway Connection payloads, see gatewayConn() for format description.
	tests := []struct {
		inp *gw.ConnState
		res []*decode.Point
		err string
	}{
		// Gateway ConnState.
		{
			&gw.ConnState{}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{}`},
				{Attr: "conn", Value: "OFFLINE"},
			}, "",
		},
		{
			&gw.ConnState{State: gw.ConnState_ONLINE}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"state":"ONLINE"}`},
				{Attr: "conn", Value: "ONLINE"},
			}, "",
		},
		// Gateway ConnState bad length.
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

			res, err := gatewayConn(bInp)
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
