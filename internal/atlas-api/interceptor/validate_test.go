//go:build !integration

package interceptor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	skipPath := random.String(10)

	badOrg := random.Org("int-valid")
	badOrg.Email = random.String(10)

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
			nil, nil, random.Org("int-valid"), &grpc.UnaryServerInfo{
				FullMethod: random.String(10),
			},
		},
		{
			status.Error(codes.InvalidArgument, "validation error: email: "+
				"must be a valid email address"), nil, badOrg,
			&grpc.UnaryServerInfo{FullMethod: random.String(10)},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can log %+v", test), func(t *testing.T) {
			t.Parallel()

			handler := func(_ context.Context, _ any) (any, error) {
				return nil, test.err
			}

			res, err := Validate(test.inpSkipPaths)(t.Context(), test.inpReq,
				test.inpInfo, handler)
			t.Logf("res, err: %v, %v", res, err)
			require.Nil(t, res)
			require.Equal(t, test.err, err)
		})
	}
}
