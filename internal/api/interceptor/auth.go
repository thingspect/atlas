package interceptor

import (
	"context"
	"strings"

	"github.com/thingspect/atlas/internal/api/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Auth performs authentication and authorization via web token, and implements
// the grpc.UnaryServerInterceptor type signature.
func Auth(skipPaths map[string]struct{},
	key []byte) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{},
		error) {
		if _, ok := skipPaths[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		// Retrieve token from 'Authorization: Bearer ...' header.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		auth, ok := md["authorization"]
		if !ok || len(auth) == 0 {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		if !strings.HasPrefix(auth[0], "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		// Validate token.
		token := strings.TrimPrefix(auth[0], "Bearer ")
		sess, err := session.ValidateToken(key, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		ctx = session.NewContext(ctx, sess)
		return handler(ctx, req)
	}
}
