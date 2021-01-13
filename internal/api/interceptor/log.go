// Package interceptor provides functions to intercept the execution of an RPC
// on the server and implement the grpc.UnaryServerInterceptor type signature.
package interceptor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Log logs requests, responses, and metadata, and implements the
// grpc.UnaryServerInterceptor type signature.
func Log(skipPaths map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{},
		error) {
		// Set up logging fields and context.
		logger := &alog.CtxLogger{Logger: alog.WithStr("path", info.FullMethod)}
		ctx = alog.NewContext(ctx, logger)

		start := time.Now()
		resp, err := handler(ctx, req)

		if _, ok := skipPaths[info.FullMethod]; ok {
			return resp, err
		}

		// Add additional logging fields.
		logFields := map[string]interface{}{
			"durms": fmt.Sprintf("%d", time.Since(start)/time.Millisecond),
			"code":  status.Code(err).String(),
		}
		logger.Logger = logger.WithFields(logFields)

		// Populate additional fields with metadata.
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for k, v := range md {
				if !strings.Contains(k, "authorization") {
					logger.Logger = logger.WithStr(k, strings.Join(v, ","))
				}
			}
		}

		if err != nil {
			logger.Info(err)
		} else {
			respOut := fmt.Sprintf("%+v", resp)
			if len(respOut) > 80 {
				respOut = fmt.Sprintf("%v...", respOut[:77])
			}

			logger.Info(respOut)
		}

		return resp, err
	}
}
