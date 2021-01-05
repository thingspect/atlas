// +build !integration

package service

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
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
		{status.Error(codes.Unauthenticated, "unauthenticated"),
			"rpc error: code = Unauthenticated desc = unauthenticated"},
		{io.EOF, "rpc error: code = Unknown desc = EOF"},
		{dao.ErrAlreadyExists, "rpc error: code = AlreadyExists desc = " +
			"object already exists"},
		{fmt.Errorf("%w: UUID", dao.ErrInvalidFormat),
			"rpc error: code = InvalidArgument desc = invalid format: UUID"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can map %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := errToStatus(lTest.inp)
			t.Logf("res: %#v", res)
			// Comparison of gRPC status errors does not play well with
			// require.Equal.
			if lTest.res == "" {
				require.NoError(t, res)
			} else {
				require.EqualError(t, res, lTest.res)
			}
		})
	}
}
