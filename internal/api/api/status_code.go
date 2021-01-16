package api

//go:generate mockgen -destination mock_respwriter_test.go -package api net/http ResponseWriter

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

// statusCode modifies the HTTP response status code based on header.
func statusCode(ctx context.Context, w http.ResponseWriter,
	p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if mdCode := md.HeaderMD["atlas-status-code"]; len(mdCode) > 0 {
		code, err := strconv.Atoi(mdCode[0])
		if err != nil {
			return err
		}

		delete(md.HeaderMD, "atlas-status-code")
		delete(w.Header(), "Grpc-Metadata-Atlas-Status-Code")
		w.WriteHeader(code)
	}

	return nil
}
