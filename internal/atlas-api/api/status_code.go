package api

//go:generate mockgen -destination mock_respwriter_test.go -package api -build_flags=-mod=mod net/http ResponseWriter

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/thingspect/atlas/internal/atlas-api/service"
	"google.golang.org/protobuf/proto"
)

const grpcStatusCodeKey = "Grpc-Metadata-Atlas-Status-Code"

// statusCode modifies the HTTP response status code based on header.
func statusCode(
	ctx context.Context, w http.ResponseWriter, p proto.Message,
) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if mdCode := md.HeaderMD[service.StatusCodeKey]; len(mdCode) > 0 {
		code, err := strconv.Atoi(mdCode[0])
		if err != nil {
			return err
		}

		delete(md.HeaderMD, service.StatusCodeKey)
		delete(w.Header(), grpcStatusCodeKey)
		w.WriteHeader(code)
	}

	return nil
}
