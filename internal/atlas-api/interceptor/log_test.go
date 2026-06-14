//go:build !integration

package interceptor

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestLog(t *testing.T) {
	t.Parallel()

	skipPath := random.String(10)

	tests := []struct {
		inpMD         []string
		inpHandlerErr error
		inpReq        string
		inpInfo       *grpc.UnaryServerInfo
	}{
		{nil, nil, random.String(105), &grpc.UnaryServerInfo{
			FullMethod: random.String(10),
		}},
		{nil, nil, random.String(10), &grpc.UnaryServerInfo{
			FullMethod: skipPath,
		}},
		{
			[]string{random.String(10), random.String(10)},
			nil, random.String(10),
			&grpc.UnaryServerInfo{FullMethod: random.String(10)},
		},
		{nil, io.EOF, "", &grpc.UnaryServerInfo{FullMethod: random.String(10)}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can log %+v", test), func(t *testing.T) {
			t.Parallel()

			ctx := metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(test.inpMD...))

			handler := func(_ context.Context, req any) (any, error) {
				return req, test.inpHandlerErr
			}

			res, err := Log()(ctx, test.inpReq, test.inpInfo, handler)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, test.inpReq, res)
			require.Equal(t, test.inpHandlerErr, err)
		})
	}
}
