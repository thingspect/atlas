package device

import (
	"encoding/hex"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	//nolint:staticcheck // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	//nolint:staticcheck // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
)

// deviceTXAck parses a device TX ACK payload from a []byte according to the
// spec.
func deviceTXAck(body []byte) ([]*decode.Point, error) {
	txAckMsg := &as.TxAckEvent{}
	if err := proto.Unmarshal(body, txAckMsg); err != nil {
		return nil, err
	}

	// Build raw device payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(txAckMsg)
	if err != nil {
		return nil, err
	}
	msgs := []*decode.Point{{Attr: "raw_device", Value: gw}}

	// Parse TxAckEvent.
	msgs = append(msgs, &decode.Point{Attr: "ack_gateway_tx", Value: true})
	if len(txAckMsg.GatewayId) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "gateway_id", Value: hex.EncodeToString(txAckMsg.GatewayId),
		})
	}

	// Parse DownlinkTXInfo.
	if txAckMsg.TxInfo != nil && txAckMsg.TxInfo.Frequency != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "frequency", Value: int32(txAckMsg.TxInfo.Frequency),
		})
	}

	return msgs, nil
}
