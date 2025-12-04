package gateway

import (
	"strconv"
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/gw"
	"github.com/thingspect/atlas/pkg/decode"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// gatewayStats parses a gateway Stats payload from a []byte according to the
// spec.
func gatewayStats(body []byte) ([]*decode.Point, error) {
	statsMsg := &gw.GatewayStats{}
	if err := proto.Unmarshal(body, statsMsg); err != nil {
		return nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_gateway", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(statsMsg), " ", "")}}

	// Parse GatewayStats.
	if len(statsMsg.GetGatewayId()) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "id", Value: statsMsg.GetGatewayId(),
		})
	}
	if statsMsg.GetTime() != nil {
		msgs = append(msgs, &decode.Point{
			Attr:  "gateway_time",
			Value: strconv.FormatInt(statsMsg.GetTime().GetSeconds(), 10),
		})
	}
	if statsMsg.GetRxPacketsReceived() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "rx_received", Value: int32(statsMsg.GetRxPacketsReceived()),
		})
	}
	if statsMsg.GetRxPacketsReceivedOk() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr:  "rx_received_valid",
			Value: int32(statsMsg.GetRxPacketsReceivedOk()),
		})
	}
	if statsMsg.GetTxPacketsReceived() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "tx_received", Value: int32(statsMsg.GetTxPacketsReceived()),
		})
	}
	if statsMsg.GetTxPacketsEmitted() != 0 {
		msgs = append(msgs, &decode.Point{
			Attr:  "tx_transmitted",
			Value: int32(statsMsg.GetTxPacketsEmitted()),
		})
	}
	for k, v := range statsMsg.GetMetadata() {
		msgs = append(msgs, &decode.Point{Attr: k, Value: v})
	}

	return msgs, nil
}
