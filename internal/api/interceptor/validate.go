package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate(all bool) error
}

// Validate performs request validation, and implements the
// grpc.UnaryServerInterceptor type signature.
func Validate(skipPaths map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{},
		error) {
		if _, ok := skipPaths[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		if v, ok := req.(validator); ok {
			if err := v.Validate(false); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		return handler(ctx, req)
	}
}
