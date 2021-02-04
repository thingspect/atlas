// +build !integration

package gateway

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/parse"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestGateway(t *testing.T) {
	t.Parallel()

	badEvent := random.String(10)

	// Trivial gateway payloads, see Gateway() for format description. Parsers
	// are exercised more thoroughly in their respective tests.
	tests := []struct {
		inpEvent string
		inpBody  string
		res      []*parse.Point
		err      error
	}{
		// Gateway.
		{"up", "", []*parse.Point{
			{Attr: "raw_gateway", Value: `{}`},
		}, nil},
		{"stats", "", []*parse.Point{
			{Attr: "raw_gateway", Value: `{}`},
		}, nil},
		{"ack", "", []*parse.Point{
			{Attr: "raw_gateway", Value: `{}`},
		}, nil},
		{"exec", "", []*parse.Point{
			{Attr: "raw_gateway", Value: `{}`},
		}, nil},
		// Gateway unknown event type.
		{badEvent, "", nil, fmt.Errorf("%w: %s, %x", parse.ErrUnknownEvent,
			badEvent, []byte{})},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can parse %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bInpBody, err := hex.DecodeString(lTest.inpBody)
			require.NoError(t, err)

			res, err := Gateway(lTest.inpEvent, bInpBody)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, lTest.res, res)
			require.Equal(t, lTest.err, err)
		})
	}
}
