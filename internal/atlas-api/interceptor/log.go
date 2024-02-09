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
func Log(skipPaths map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Set up logging fields and context.
		logger := &alog.CtxLogger{Logger: alog.WithField("path", info.FullMethod)}
		ctx = alog.NewContext(ctx, logger)

		start := time.Now()
		resp, err := handler(ctx, req)
		dur := time.Since(start)

		// Send metrics.
		metric.Timing("durms", dur, nil)
		if err != nil {
			metric.Incr("error", nil)
		}

		if _, ok := skipPaths[info.FullMethod]; ok {
			return resp, err
		}

		// Add additional logging fields.
		logger.Logger = logger.
			WithField("durms", fmt.Sprintf("%d", dur/time.Millisecond)).
			WithField("code", status.Code(err).String())

		// Populate additional fields with metadata.
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for k, v := range md {
				if !strings.Contains(k, "authorization") {
					logger.Logger = logger.WithField(k, strings.Join(v, ","))
				}
			}
		}

		if err != nil {
			logger.Info(err)
		} else {
			respOut := fmt.Sprintf("%+v", resp)
			if len(respOut) > 80 {
				respOut = respOut[:77] + "..."
			}

			logger.Info(respOut)
		}

		return resp, err
	}
}
