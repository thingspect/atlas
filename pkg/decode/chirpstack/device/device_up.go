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

	if upMsg.GetData() != nil {
		msgs = append(msgs, &decode.Point{
			Attr: "raw_data", Value: hex.EncodeToString(upMsg.GetData()),
		})
	}

	// Parse UplinkRXInfos.
	upTime, rxMsgs := chirpstack.ParseRXInfos(upMsg.GetRxInfo())
	msgs = append(msgs, rxMsgs...)

	// Parse UplinkTXInfo.
	if upMsg.GetTxInfo().GetFrequency() != 0 {
		msgs = append(msgs, &decode.Point{
			//nolint:gosec // Safe conversion for limited values.
			Attr: "frequency", Value: int32(upMsg.GetTxInfo().GetFrequency()),
		})
	}

	// Parse UplinkEvent.
	msgs = append(msgs, &decode.Point{Attr: "adr", Value: upMsg.GetAdr()})
	msgs = append(msgs, &decode.Point{
		//nolint:gosec // Safe conversion for limited values.
		Attr: "data_rate", Value: int32(upMsg.GetDr()),
	})
	msgs = append(msgs, &decode.Point{
		Attr: "confirmed", Value: upMsg.GetConfirmed(),
	})
	if upMsg.GetRegionConfigId() != "" {
		msgs = append(msgs, &decode.Point{
			Attr: "region_config_id", Value: upMsg.GetRegionConfigId(),
		})
	}

	return msgs, upTime, upMsg.GetData(), nil
}
