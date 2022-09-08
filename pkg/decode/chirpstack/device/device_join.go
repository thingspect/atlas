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
	"github.com/thingspect/atlas/pkg/decode/chirpstack"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// deviceJoin parses a device Join payload from a []byte according to the spec.
func deviceJoin(body []byte) ([]*decode.Point, *timestamppb.Timestamp, error) {
	joinMsg := &as.JoinEvent{}
	if err := proto.Unmarshal(body, joinMsg); err != nil {
		return nil, nil, err
	}

	// Build raw device payload for debugging.
	marshaler := &jsonpb.Marshaler{}
	gw, err := marshaler.MarshalToString(joinMsg)
	if err != nil {
		return nil, nil, err
	}
	msgs := []*decode.Point{{Attr: "raw_device", Value: gw}}

	// Parse JoinEvent.
	msgs = append(msgs, &decode.Point{Attr: "join", Value: true})
	if len(joinMsg.DevEui) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "id", Value: hex.EncodeToString(joinMsg.DevEui),
		})
	}
	if len(joinMsg.DevAddr) != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "devaddr", Value: hex.EncodeToString(joinMsg.DevAddr),
		})
	}

	// Parse UplinkRXInfos.
	joinTime, rxMsgs := chirpstack.ParseRXInfos(joinMsg.RxInfo)
	msgs = append(msgs, rxMsgs...)

	// Parse UplinkTXInfo.
	if joinMsg.TxInfo != nil && joinMsg.TxInfo.Frequency != 0 {
		msgs = append(msgs, &decode.Point{
			Attr: "frequency", Value: int32(joinMsg.TxInfo.Frequency),
		})
	}

	// Parse JoinEvent data rate.
	msgs = append(msgs, &decode.Point{
		Attr: "data_rate", Value: int32(joinMsg.Dr),
	})

	return msgs, joinTime, nil
}
