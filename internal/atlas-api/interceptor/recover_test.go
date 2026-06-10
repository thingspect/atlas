//go:build !integration

package interceptor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRecover(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpHandler grpc.UnaryHandler
		err        error
	}{
		{func(_ context.Context, _ any) (any, error) { return true, nil }, nil},
		{func(_ context.Context, _ any) (any, error) {
			panic("panic")
		}, status.Error(codes.Internal, "internal server error")},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can recover %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := Recover()(t.Context(), nil, nil, test.inpHandler)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, test.err, err)
		})
	}
}
