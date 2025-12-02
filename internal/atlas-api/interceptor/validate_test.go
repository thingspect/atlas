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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type valPass struct{}

func (v *valPass) Validate() error { return nil }

type valFail struct{}

func (v *valFail) Validate() error { return io.EOF }

func TestValidate(t *testing.T) {
	t.Parallel()

	skipPath := random.String(10)

	tests := []struct {
		err          error
		inpSkipPaths map[string]struct{}
		inpReq       any
		inpInfo      *grpc.UnaryServerInfo
	}{
		{
			nil, nil, nil, &grpc.UnaryServerInfo{FullMethod: random.String(10)},
		},
		{
			nil, map[string]struct{}{skipPath: {}}, nil, &grpc.UnaryServerInfo{
				FullMethod: skipPath,
			},
		},
		{
			nil, nil, &valPass{}, &grpc.UnaryServerInfo{
				FullMethod: random.String(10),
			},
		},
		{
			status.Error(codes.InvalidArgument, io.EOF.Error()), nil,
			&valFail{}, &grpc.UnaryServerInfo{FullMethod: random.String(10)},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can log %+v", test), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()

			handler := func(_ context.Context, _ any) (any, error) {
				return nil, test.err
			}

			res, err := Validate(test.inpSkipPaths)(ctx, test.inpReq,
				test.inpInfo, handler)
			t.Logf("res, err: %v, %v", res, err)
			require.Nil(t, res)
			require.Equal(t, test.err, err)
		})
	}
}
