// +build !integration

package registry

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/decode"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
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
		{&decode.Point{Attr: "data_rate", Value: int32(3)},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "data_rate", ValOneof: &common.DataPoint_IntVal{
					IntVal: 3}, Ts: now, TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "data_rate", Value: 3}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "data_rate",
				ValOneof: &common.DataPoint_IntVal{IntVal: 3}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "data_rate", Value: int64(3)},
			&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "data_rate", ValOneof: &common.DataPoint_IntVal{
					IntVal: 3}, Ts: now, TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "snr", Value: 7.8}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "snr",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 7.8}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "snr", Value: float32(7.0)}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "snr",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 7.0}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "ack", Value: "OK"}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "ack",
				ValOneof: &common.DataPoint_StrVal{StrVal: "OK"}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "adr", Value: false}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "adr",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: false}, Ts: now,
				TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "ack", Value: []byte{0x00}}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "ack",
				ValOneof: &common.DataPoint_BytesVal{BytesVal: []byte{0x00}},
				Ts:       now, TraceId: traceID}, SkipToken: true}},
		{&decode.Point{Attr: "error", Value: io.EOF}, &message.ValidatorIn{
			Point: &common.DataPoint{UniqId: uniqID, Attr: "error", Ts: now,
				TraceId: traceID}, SkipToken: true}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can convert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := PointToVIn(traceID, uniqID, lTest.inp, now)
			t.Logf("res: %+v", res)

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.res, res) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, res)
			}
		})
	}
}
