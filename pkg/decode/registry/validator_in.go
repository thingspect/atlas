package registry

import (
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PointToVIn converts a parsed Point to ValidatorIn.
func PointToVIn(
	traceID, uniqID string, point *decode.Point, ts *timestamppb.Timestamp,
) *message.ValidatorIn {
	vIn := &message.ValidatorIn{
		Point: &common.DataPoint{
			UniqId:  uniqID,
			Attr:    point.Attr,
			Ts:      ts,
			TraceId: traceID,
		},
		SkipToken: true,
	}

	switch v := point.Value.(type) {
	case int32:
		vIn.Point.ValOneof = &common.DataPoint_IntVal{IntVal: v}
	case int:
		vIn.Point.ValOneof = &common.DataPoint_IntVal{IntVal: int32(v)}
		alog.Errorf("PointToVIn casting from int: %v, %v,", point.Attr, v)
	case int64:
		vIn.Point.ValOneof = &common.DataPoint_IntVal{IntVal: int32(v)}
		alog.Errorf("PointToVIn casting from int64: %v, %v,", point.Attr, v)
	case float64:
		vIn.Point.ValOneof = &common.DataPoint_Fl64Val{Fl64Val: v}
	case float32:
		vIn.Point.ValOneof = &common.DataPoint_Fl64Val{Fl64Val: float64(v)}
		alog.Errorf("PointToVIn casting from float32: %v, %v,", point.Attr, v)
	case string:
		vIn.Point.ValOneof = &common.DataPoint_StrVal{StrVal: v}
	case bool:
		vIn.Point.ValOneof = &common.DataPoint_BoolVal{BoolVal: v}
	case []byte:
		vIn.Point.ValOneof = &common.DataPoint_BytesVal{BytesVal: v}
	default:
		alog.Errorf("PointToVIn unknown type: %v, %T, %v,", point.Attr,
			point.Value, point.Value)
	}

	return vIn
}
