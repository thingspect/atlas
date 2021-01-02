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
		start := time.Now()
		resp, err := handler(ctx, req)
		dur := fmt.Sprintf("%d", time.Since(start)/time.Millisecond)

		if _, ok := skipPaths[info.FullMethod]; ok {
			return resp, err
		}

		// Set up logging fields.
		logFields := map[string]interface{}{
			"path":  info.FullMethod,
			"durms": dur,
			"code":  status.Code(err).String(),
		}
		logEntry := alog.WithFields(logFields)

		// Populate additional fields with metadata.
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for k, v := range md {
				if len(v) > 0 {
					logEntry = logEntry.WithStr(k, strings.Join(v, ","))
				}
			}
		}

		if err != nil {
			logEntry.Info(err)
		} else {
			respOut := fmt.Sprintf("%+v", resp)
			if len(respOut) > 100 {
				respOut = fmt.Sprintf("%v...", respOut[:97])
			}

			logEntry.Info(respOut)
		}
		return resp, err
	}
}
