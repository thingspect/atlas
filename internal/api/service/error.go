// Package service provides functions that implement gRPC service interfaces.
package service

import (
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errToCode maps DAO errors to gRPC error codes.
var errToCode = map[error]codes.Code{
	dao.ErrAlreadyExists: codes.AlreadyExists,
	dao.ErrInvalidFormat: codes.InvalidArgument,
	dao.ErrNotFound:      codes.NotFound,
}

// errToStatus maps DAO errors to gRPC status errors. This function is
// idempotent and is safe to call on the same error multiple times.
func errToStatus(err error) error {
	// If err is nil or is already a gRPC status, return it.
	if code := status.Code(err); code != codes.Unknown {
		return err
	}

	code, ok := errToCode[err]
	if !ok {
		code = codes.Unknown
		alog.Errorf("errToStatus unmatched error: %#v", err)
	}

	return status.Error(code, err.Error())
}
