// +build !integration

package gateway

import (
	"fmt"
	"testing"

	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/parse"
)

func TestGatewayExec(t *testing.T) {
	t.Parallel()

	// Gateway Exec payloads, see gatewayExec() for format description.
	tests := []struct {
		inp *gw.GatewayCommandExecResponse
		res []*parse.Point
		err string
	}{
		// Gateway Exec.
		{&gw.GatewayCommandExecResponse{}, []*parse.Point{
			{Attr: "raw_gateway", Value: `{}`},
		}, ""},
		{&gw.GatewayCommandExecResponse{Stdout: []byte("STDOUT"),
			Stderr: []byte("STDERR"), Error: "TOO_LATE"}, []*parse.Point{
			{Attr: "raw_gateway", Value: `{"stdout":"U1RET1VU","stderr":` +
				`"U1RERVJS","error":"TOO_LATE"}`},
			{Attr: "exec_stdout", Value: "STDOUT"},
			{Attr: "exec_stderr", Value: "STDERR"},
			{Attr: "exec_error", Value: "TOO_LATE"},
		}, ""},
		// Gateway Exec bad length.
		{nil, nil, "unexpected EOF"},
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

			res, err := gatewayExec(bInp)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			if lTest.err != "" {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
