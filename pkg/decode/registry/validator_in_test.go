//go:build !integration

package registry

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPointToVIn(t *testing.T) {
	t.Parallel()

	traceID := uuid.NewString()
	uniqID := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	tests := []struct {
		inp *decode.Point
		res *message.ValidatorIn
	}{
		{
			&decode.Point{Attr: "data_rate", Value: int32(3)},
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: uniqID, Attr: "data_rate",
					ValOneof: &common.DataPoint_IntVal{IntVal: 3}, Ts: now,
					TraceId: traceID,
				}, SkipToken: true,
			},
		},
		{
			&decode.Point{Attr: "data_rate", Value: 3}, &message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: uniqID, Attr: "data_rate",
					ValOneof: &common.DataPoint_IntVal{IntVal: 3}, Ts: now,
					TraceId: traceID,
				}, SkipToken: true,
			},
		},
		{
			&decode.Point{Attr: "data_rate", Value: int64(3)},
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: uniqID, Attr: "data_rate",
					ValOneof: &common.DataPoint_IntVal{IntVal: 3}, Ts: now,
					TraceId: traceID,
				}, SkipToken: true,
			},
		},
		{
			&decode.Point{Attr: "lora_snr", Value: 7.8}, &message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: uniqID, Attr: "lora_snr",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 7.8}, Ts: now,
					TraceId: traceID,
				}, SkipToken: true,
			},
		},
		{
			&decode.Point{Attr: "lora_snr", Value: float32(7.0)},
			&message.ValidatorIn{Point: &common.DataPoint{
				UniqId: uniqID, Attr: "lora_snr",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 7.0}, Ts: now,
				TraceId: traceID,
			}, SkipToken: true},
		},
		{
			&decode.Point{Attr: "ack", Value: "OK"}, &message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: uniqID, Attr: "ack",
					ValOneof: &common.DataPoint_StrVal{StrVal: "OK"}, Ts: now,
					TraceId: traceID,
				}, SkipToken: true,
			},
		},
		{
			&decode.Point{Attr: "adr", Value: false}, &message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: uniqID, Attr: "adr",
					ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
					Ts:       now, TraceId: traceID,
				}, SkipToken: true,
			},
		},
		{
			&decode.Point{Attr: "ack", Value: []byte{0x00}},
			&message.ValidatorIn{Point: &common.DataPoint{
				UniqId: uniqID, Attr: "ack",
				ValOneof: &common.DataPoint_BytesVal{BytesVal: []byte{0x00}},
				Ts:       now, TraceId: traceID,
			}, SkipToken: true},
		},
		{
			&decode.Point{Attr: "error", Value: io.EOF}, &message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: uniqID, Attr: "error", Ts: now, TraceId: traceID,
				}, SkipToken: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can convert %+v", test), func(t *testing.T) {
			t.Parallel()

			res := PointToVIn(traceID, uniqID, test.inp, now)
			t.Logf("res: %+v", res)
			require.EqualExportedValues(t, test.res, res)
		})
	}
}
