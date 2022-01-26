package device

//nolint:staticcheck // third-party dependency
import (
	"encoding/hex"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/decode/chirpstack"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// deviceUp parses a device Uplink payload from a []byte according to the spec.
// Points, a timestamp, and a data []byte are built from successful parse
// results. If a fatal error is encountered, it is returned along with any valid
// points.
func deviceUp(body []byte) ([]*decode.Point, *timestamppb.Timestamp, []byte,
	error) {
	upMsg := &as.UplinkEvent{}
	if err := proto.Unmarshal(body, upMsg); err != nil {
		return nil, nil, nil, err
	}

	// Build raw device and data payloads for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(upMsg)
	if err != nil {
		return nil, nil, upMsg.Data, err
	}
	msgs := []*decode.Point{{Attr: "raw_device", Value: gw}}

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
		Attr: "confirmed", Value: upMsg.ConfirmedUplink,
	})

	return msgs, upTime, upMsg.Data, nil
}
