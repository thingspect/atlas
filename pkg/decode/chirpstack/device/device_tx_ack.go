package device

import (
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// deviceTXAck parses a device TX ACK payload from a []byte according to the
// spec.
func deviceTXAck(body []byte) ([]*decode.Point, error) {
	txAckMsg := &integration.TxAckEvent{}
	if err := proto.Unmarshal(body, txAckMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_device", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(txAckMsg), " ", "")}}

	// Parse TxAckEvent.
	msgs = append(msgs, &decode.Point{Attr: "tx_queued", Value: true})
	if txAckMsg.GetGatewayId() != "" {
		msgs = append(msgs, &decode.Point{
			Attr: "tx_gateway_id", Value: txAckMsg.GetGatewayId(),
		})
	}

	// Parse DownlinkTXInfo.
	if txAckMsg.GetTxInfo().GetFrequency() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "tx_frequency", Value: int32(txAckMsg.GetTxInfo().GetFrequency()),
		})
	}

	return msgs, nil
}
