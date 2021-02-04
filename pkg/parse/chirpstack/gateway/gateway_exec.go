package gateway

import (
	"github.com/brocaar/chirpstack-api/go/v3/gw"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/parse"
)

// gatewayExec parses a gateway exec command response payload from a []byte
// according to the spec.
func gatewayExec(body []byte) ([]*parse.Point, error) {
	execMsg := &gw.GatewayCommandExecResponse{}
	if err := proto.Unmarshal(body, execMsg); err != nil {
		return nil, err
	}

	// Build raw gateway payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(execMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*parse.Point{{Attr: "raw_gateway", Value: gw}}

	// Parse GatewayCommandExecResponse.
	if len(execMsg.Stdout) != 0 {
		msgs = append(msgs, &parse.Point{Attr: "exec_stdout",
			Value: string(execMsg.Stdout)})
	}
	if len(execMsg.Stderr) != 0 {
		msgs = append(msgs, &parse.Point{Attr: "exec_stderr",
			Value: string(execMsg.Stderr)})
	}
	if execMsg.Error != "" {
		msgs = append(msgs, &parse.Point{Attr: "exec_error",
			Value: execMsg.Error})
	}

	return msgs, nil
}
