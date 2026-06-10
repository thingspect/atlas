package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Recover recovers from panics and replaces them with internal server errors,
// and implements the grpc.UnaryServerInterceptor type signature.
func Recover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				alog.FromContext(ctx).Errorf("Recover: %v - %s", r,
					debug.Stack())
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
