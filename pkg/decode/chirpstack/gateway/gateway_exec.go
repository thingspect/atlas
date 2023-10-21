package gateway

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// gatewayExec parses a gateway exec command response payload from a []byte
// according to the spec.
func gatewayExec(body []byte) ([]*decode.Point, error) {
	execMsg := &gw.GatewayCommandExecResponse{}
	if err := proto.Unmarshal(body, execMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(execMsg), " ", "")}}

	// Parse GatewayCommandExecResponse.
	if len(execMsg.GetStdout()) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "exec_stdout", Value: string(execMsg.GetStdout()),
		})
	}
	if len(execMsg.GetStderr()) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "exec_stderr", Value: string(execMsg.GetStderr()),
		})
	}
	if execMsg.GetError() != "" {
		msgs = append(msgs, &decode.Point{
			Attr: "exec_error", Value: execMsg.GetError(),
		})
	}

	return msgs, nil
}
