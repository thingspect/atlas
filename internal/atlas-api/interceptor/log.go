// Package interceptor provides functions to intercept the execution of an RPC
// on the server and implement the grpc.UnaryServerInterceptor type signature.
package interceptor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Log logs requests, responses, and metadata, sends metrics, and implements the
// grpc.UnaryServerInterceptor type signature.
func Log() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Set up logging fields and context.
		logger := &alog.CtxLogger{Logger: alog.WithField(
			"path", info.FullMethod)}
		ctx = alog.NewContext(ctx, logger)

		// Populate additional fields with metadata.
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for k, v := range md {
				if !strings.Contains(k, "authorization") {
					logger.Logger = logger.WithField(k, strings.Join(v, ","))
				}
			}
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		dur := time.Since(start)

		// Send metrics.
		metric.Timing("durms", dur, map[string]string{"path": info.FullMethod})
		if err != nil {
			metric.Incr("error", map[string]string{"path": info.FullMethod})
		}

		// Add additional logging fields for final logging. Do not modify logger
		// once it has been passed downstream where it may be used by multiple
		// goroutines.
		flog := logger.
			WithField("durms", fmt.Sprintf("%d", dur/time.Millisecond)).
			WithField("code", status.Code(err).String())

		if err != nil {
			flog.Info(err)
		} else {
			respOut := fmt.Sprintf("%+v", resp)
			if len(respOut) > 80 {
				respOut = respOut[:77] + "..."
			}

			flog.Info(respOut)
		}

		return resp, err
	}
}
