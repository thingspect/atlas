// +build !integration

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
		{nil, nil, nil, random.String(105), &grpc.UnaryServerInfo{
			FullMethod: random.String(10)}},
		{nil, nil, map[string]struct{}{skipPath: {}}, random.String(10),
			&grpc.UnaryServerInfo{FullMethod: skipPath}},
		{[]string{random.String(10), random.String(10)}, nil, nil,
			random.String(10), &grpc.UnaryServerInfo{
				FullMethod: random.String(10)}},
		{nil, io.EOF, nil, "", &grpc.UnaryServerInfo{
			FullMethod: random.String(10)}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can log %+v", lTest), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()
			ctx = metadata.NewIncomingContext(ctx,
				metadata.Pairs(lTest.inpMD...))

			handler := func(ctx context.Context, req interface{}) (interface{},
				error) {
				return req, lTest.inpHandlerErr
			}

			res, err := Log(lTest.inpSkipPaths)(ctx, lTest.inpReq,
				lTest.inpInfo, handler)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, lTest.inpReq, res)
			require.Equal(t, lTest.inpHandlerErr, err)
		})
	}
}
