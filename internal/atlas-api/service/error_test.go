//go:build !integration

package service

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrToStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inp error
		res string
	}{
		{nil, ""},
		{
			status.Error(codes.Unauthenticated, "unauthenticated"),
			"rpc error: code = Unauthenticated desc = unauthenticated",
		},
		{
			io.EOF, "rpc error: code = Unknown desc = EOF",
		},
		{
			dao.ErrAlreadyExists,
			"rpc error: code = AlreadyExists desc = object already exists",
		},
		{
			fmt.Errorf("%w: UUID", dao.ErrInvalidFormat),
			"rpc error: code = InvalidArgument desc = invalid format: UUID",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can map %+v", test), func(t *testing.T) {
			t.Parallel()

			res := errToStatus(test.inp)
			t.Logf("res: %#v", res)
			// Comparison of gRPC status errors does not play well with
			// require.Equal.
			if test.res == "" {
				require.NoError(t, res)
			} else {
				require.EqualError(t, res, test.res)
			}
		})
	}
}

func TestErrPerm(t *testing.T) {
	t.Parallel()

	for i := range 5 {
		t.Run(fmt.Sprintf("Can generate %v", i), func(t *testing.T) {
			t.Parallel()

			role := []api.Role{
				api.Role_CONTACT, api.Role_VIEWER, api.Role_BUILDER,
				api.Role_ADMIN, api.Role_SYS_ADMIN,
			}[random.Intn(5)]

			err := errPerm(role)
			t.Logf("err: %v", err)

			require.Equal(t, status.Error(codes.PermissionDenied,
				fmt.Sprintf("permission denied, %s role required",
					role.String())), err)
		})
	}
}
