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

func TestGatewayExec(t *testing.T) {
	t.Parallel()

	// Gateway Exec payloads, see gatewayExec() for format description.
	tests := []struct {
		inp *gw.GatewayCommandExecResponse
		res []*decode.Point
		err string
	}{
		// Gateway Exec.
		{
			&gw.GatewayCommandExecResponse{}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{}`},
			}, "",
		},
		{
			&gw.GatewayCommandExecResponse{
				Stdout: []byte("STDOUT"), Stderr: []byte("STDERR"),
				Error: "TOO_LATE",
			}, []*decode.Point{
				{Attr: "raw_gateway", Value: `{"stdout":"U1RET1VU","stderr":` +
					`"U1RERVJS","error":"TOO_LATE"}`},
				{Attr: "exec_stdout", Value: "STDOUT"},
				{Attr: "exec_stderr", Value: "STDERR"},
				{Attr: "exec_error", Value: "TOO_LATE"},
			}, "",
		},
		// Gateway Exec bad length.
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

			res, err := gatewayExec(bInp)
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
