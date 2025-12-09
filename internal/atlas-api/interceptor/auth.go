package interceptor

import (
	"context"
	"errors"
	"strings"

	"github.com/thingspect/atlas/internal/atlas-api/key"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Auth performs authentication and authorization via web token, and implements
// the grpc.UnaryServerInterceptor type signature.
func Auth(
	skipPaths map[string]struct{}, pwtKey []byte, c cache.Cacher[string],
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if _, ok := skipPaths[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		// Retrieve token from 'Authorization: Bearer ...' header.
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		auth := md["authorization"]
		if len(auth) < 1 {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		if !strings.HasPrefix(auth[0], "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		// Validate token.
		token := strings.TrimPrefix(auth[0], "Bearer ")
		sess, err := session.ValidateWebToken(pwtKey, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		// Check for disabled API key. Disabled can return nil or non-nil except
		//  cache.ErrNotFound.
		if sess.KeyID != "" {
			if _, err := c.Get(ctx, key.Disabled(sess.OrgID,
				sess.KeyID)); !errors.Is(err, cache.ErrNotFound) {
				return nil, status.Error(codes.Unauthenticated, "unauthorized")
			}
		}

		// Add logging fields.
		logger := alog.FromContext(ctx)
		if sess.UserID != "" {
			logger.Logger = logger.WithField("userID", sess.UserID)
		} else {
			logger.Logger = logger.WithField("keyID", sess.KeyID)
		}
		logger.Logger = logger.WithField("orgID", sess.OrgID)
		logger.Logger = logger.WithField("traceID", sess.TraceID.String())

		ctx = session.NewContext(ctx, sess)

		return handler(ctx, req)
	}
}
