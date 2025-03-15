//go:build !integration

package interceptor

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const testTimeout = 2 * time.Second

func TestLog(t *testing.T) {
	t.Parallel()

	skipPath := random.String(10)

	tests := []struct {
		inpMD         []string
		inpHandlerErr error
		inpSkipPaths  map[string]struct{}
		inpReq        string
		inpInfo       *grpc.UnaryServerInfo
	}{
		{
			nil, nil, nil, random.String(105), &grpc.UnaryServerInfo{
				FullMethod: random.String(10),
			},
		},
		{
			nil, nil,
			map[string]struct{}{skipPath: {}},
			random.String(10),
			&grpc.UnaryServerInfo{FullMethod: skipPath},
		},
		{
			[]string{random.String(10), random.String(10)},
			nil, nil,
			random.String(10), &grpc.UnaryServerInfo{
				FullMethod: random.String(10),
			},
		},
		{
			nil, io.EOF, nil, "", &grpc.UnaryServerInfo{
				FullMethod: random.String(10),
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can log %+v", test), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()
			ctx = metadata.NewIncomingContext(ctx,
				metadata.Pairs(test.inpMD...))

			handler := func(_ context.Context, req interface{}) (
				interface{}, error,
			) {
				return req, test.inpHandlerErr
			}

			res, err := Log(test.inpSkipPaths)(ctx, test.inpReq,
				test.inpInfo, handler)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, test.inpReq, res)
			require.Equal(t, test.inpHandlerErr, err)
		})
	}
}
