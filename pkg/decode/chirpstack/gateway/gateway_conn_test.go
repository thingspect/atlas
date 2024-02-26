//go:build !integration

package gateway

import (
	"fmt"
	"testing"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/proto"
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
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInp := []byte("aaa")
			if test.inp != nil {
				var err error
				bInp, err = proto.Marshal(test.inp)
				require.NoError(t, err)
			}

			res, err := gatewayConn(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, test.res, res)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
