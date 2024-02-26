//go:build !integration

package gateway

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestParse(t *testing.T) {
	t.Parallel()

	badEvent := random.String(10)

	// Trivial gateway payloads, see Parse() for format description. Parsers
	// are exercised more thoroughly in their respective tests.
	tests := []struct {
		inpEvent string
		inpBody  string
		res      []*decode.Point
		err      error
	}{
		// Gateway.
		{
			"up", "",
			[]*decode.Point{{Attr: "raw_gateway", Value: `{}`}},
			nil,
		},
		{
			"stats", "",
			[]*decode.Point{{Attr: "raw_gateway", Value: `{}`}},
			nil,
		},
		{
			"ack", "", []*decode.Point{{Attr: "raw_gateway", Value: `{}`}}, nil,
		},
		{
			"exec", "",
			[]*decode.Point{{Attr: "raw_gateway", Value: `{}`}},
			nil,
		},
		{
			"conn", "", []*decode.Point{
				{Attr: "raw_gateway", Value: `{}`},
				{Attr: "conn", Value: "OFFLINE"},
			}, nil,
		},
		// Gateway unknown event type.
		{
			badEvent, "", nil, fmt.Errorf("%w: %s, %x", decode.ErrUnknownEvent,
				badEvent, []byte{}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can parse %+v", test), func(t *testing.T) {
			t.Parallel()

			bInpBody, err := hex.DecodeString(test.inpBody)
			require.NoError(t, err)

			res, err := Parse(test.inpEvent, bInpBody)
			t.Logf("res, err: %#v, %v", res, err)
			require.Equal(t, test.res, res)
			require.Equal(t, test.err, err)
		})
	}
}
