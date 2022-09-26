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
	if len(statsMsg.GatewayId) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "id", Value: statsMsg.GatewayId,
		})
	}
	if statsMsg.Time != nil {
		msgs = append(msgs, &decode.Point{
			Attr: "time", Value: strconv.FormatInt(statsMsg.Time.Seconds, 10),
		})
	}
	if statsMsg.RxPacketsReceived != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "rx_received", Value: int32(statsMsg.RxPacketsReceived),
		})
	}
	if statsMsg.RxPacketsReceivedOk != 0 {
		msgs = append(msgs, &decode.Point{
			Attr:  "rx_received_valid",
			Value: int32(statsMsg.RxPacketsReceivedOk),
		})
	}
	if statsMsg.TxPacketsReceived != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "tx_received", Value: int32(statsMsg.TxPacketsReceived),
		})
	}
	if statsMsg.TxPacketsEmitted != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "tx_transmitted", Value: int32(statsMsg.TxPacketsEmitted),
		})
	}
	for k, v := range statsMsg.MetaData {
		msgs = append(msgs, &decode.Point{Attr: k, Value: v})
	}

	return msgs, nil
}
