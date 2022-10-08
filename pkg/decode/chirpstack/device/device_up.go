package device

import (
	"encoding/hex"
	"strings"

	"github.com/chirpstack/chirpstack/api/go/v4/integration"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/decode/chirpstack"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// deviceUp parses a device Uplink payload from a []byte according to the spec.
// Points, a timestamp, and a data []byte are built from successful parse
// results.
func deviceUp(body []byte) (
	[]*decode.Point, *timestamppb.Timestamp, []byte, error,
) {
	upMsg := &integration.UplinkEvent{}
	if err := proto.Unmarshal(body, upMsg); err != nil {
		return nil, nil, nil, err
	}

	// Build raw device and data payloads for debugging, with consistent output.
	msgs := []*decode.Point{{Attr: "raw_device", Value: strings.ReplaceAll(
		protojson.MarshalOptions{}.Format(upMsg), " ", "")}}

	if upMsg.Data != nil {
		msgs = append(msgs, &decode.Point{
			Attr: "raw_data", Value: hex.EncodeToString(upMsg.Data),
		})
	}

	// Parse UplinkRXInfos.
	upTime, rxMsgs := chirpstack.ParseRXInfos(upMsg.RxInfo)
	msgs = append(msgs, rxMsgs...)

	// Parse UplinkTXInfo.
	if upMsg.TxInfo != nil && upMsg.TxInfo.Frequency != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "frequency", Value: int32(upMsg.TxInfo.Frequency),
		})
	}

	// Parse UplinkEvent.
	msgs = append(msgs, &decode.Point{Attr: "adr", Value: upMsg.Adr})
	msgs = append(msgs, &decode.Point{
		Attr: "data_rate", Value: int32(upMsg.Dr),
	})
	msgs = append(msgs, &decode.Point{
		Attr: "confirmed", Value: upMsg.Confirmed,
	})

	return msgs, upTime, upMsg.Data, nil
}
