package device

//nolint:staticcheck // third-party dependency
import (
	"encoding/hex"
	"sort"
	"strconv"
	"time"

	as "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/jsonpb"

	//lint:ignore SA1019 // third-party dependency
	"github.com/golang/protobuf/proto"
	"github.com/thingspect/atlas/pkg/decode"
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
	joinTime := timestamppb.Now()
	if len(joinMsg.RxInfo) > 0 {
		// Sort joinMsg.RxInfo by strongest RSSI.
		sort.Slice(joinMsg.RxInfo, func(i, j int) bool {
			return joinMsg.RxInfo[i].Rssi > joinMsg.RxInfo[j].Rssi
		})

		if len(joinMsg.RxInfo[0].GatewayId) != 0 {
			msgs = append(msgs, &decode.Point{
				Attr:  "gateway_id",
				Value: hex.EncodeToString(joinMsg.RxInfo[0].GatewayId),
			})
		}

		// Populate time channel if it is provided by the gateway. Use it as
		// joinTime if it is accurate.
		if joinMsg.RxInfo[0].Time != nil {
			msgs = append(msgs, &decode.Point{
				Attr:  "time",
				Value: strconv.FormatInt(joinMsg.RxInfo[0].Time.Seconds, 10),
			})

			ts := joinMsg.RxInfo[0].Time.AsTime()
			if ts.Before(joinTime.AsTime()) &&
				time.Since(ts) < decode.ValidWindow {
				joinTime = joinMsg.RxInfo[0].Time
			}
		}

		if joinMsg.RxInfo[0].Rssi != 0 {
			msgs = append(msgs, &decode.Point{
				Attr: "lora_rssi", Value: joinMsg.RxInfo[0].Rssi,
			})
		}
		if joinMsg.RxInfo[0].LoraSnr != 0 {
			msgs = append(msgs, &decode.Point{
				Attr: "snr", Value: joinMsg.RxInfo[0].LoraSnr,
			})
		}
		msgs = append(msgs, &decode.Point{
			Attr: "channel", Value: int32(joinMsg.RxInfo[0].Channel),
		})
	}

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
