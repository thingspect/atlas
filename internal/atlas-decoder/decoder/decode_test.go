//go:build !integration

package decoder

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/decode/registry"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const errTestProc consterr.Error = "decoder: test processor error"

func TestDecodeMessages(t *testing.T) {
	t.Parallel()

	uniqID := "dec-" + random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.New()

	tests := []struct {
		inpDIn     *message.DecoderIn
		inpDecoder api.Decoder
		res        []*message.ValidatorIn
	}{
		{
			&message.DecoderIn{
				UniqId: uniqID, Data: []byte{0x19, 0x03, 0x01}, Ts: now,
				TraceId: traceID[:],
			}, api.Decoder_RADIO_BRIDGE_DOOR_V1, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "count",
						ValOneof: &common.DataPoint_IntVal{IntVal: 9}, Ts: now,
						TraceId: traceID.String(),
					}, SkipToken: true,
				}, {
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "open",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: true},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
			},
		},
		{
			&message.DecoderIn{
				UniqId: uniqID, Data: []byte{0x1a, 0x03, 0x00}, Ts: now,
				TraceId: traceID[:],
			}, api.Decoder_RADIO_BRIDGE_DOOR_V2, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "count",
						ValOneof: &common.DataPoint_IntVal{IntVal: 10}, Ts: now,
						TraceId: traceID.String(),
					}, SkipToken: true,
				}, {
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "open",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
			},
		},
		{
			&message.DecoderIn{
				UniqId: uniqID, Data: []byte{
					0x01, 0x09, 0x61, 0x13, 0x95, 0x02, 0x92,
				}, Ts: now, TraceId: traceID[:],
			}, api.Decoder_GLOBALSAT_CO2, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "temp_c",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 24},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "temp_f",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 75.2},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "humidity_pct",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 50.13},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "co2_ppm",
						ValOneof: &common.DataPoint_IntVal{IntVal: 658},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
			},
		},
		{
			&message.DecoderIn{
				UniqId: uniqID, Data: []byte{
					0x03, 0x67, 0x00, 0xc4, 0x04, 0x68, 0x7f, 0x00, 0xff, 0x01,
					0x38,
				}, Ts: now, TraceId: traceID[:],
			}, api.Decoder_TEKTELIC_HOME, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "temp_c",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 19.6},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "temp_f",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 67.3},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "humidity_pct",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 63.5},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: uniqID, Attr: "battery_v",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 3.12},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
			},
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			dev := random.Device("dec", uuid.NewString())
			dev.UniqId = uniqID
			dev.Decoder = lTest.inpDecoder

			dInQueue := queue.NewFake()
			dInSub, err := dInQueue.Subscribe("")
			require.NoError(t, err)

			decQueue := queue.NewFake()
			vInSub, err := decQueue.Subscribe("")
			require.NoError(t, err)
			vInPubTopic := "topic-" + random.String(10)

			devicer := NewMockdevicer(gomock.NewController(t))
			devicer.EXPECT().ReadByUniqID(gomock.Any(), lTest.inpDIn.GetUniqId()).
				Return(dev, nil).Times(1)

			dec := Decoder{
				devDAO: devicer,
				reg:    registry.New(),

				decQueue:    decQueue,
				dInSub:      dInSub,
				vInPubTopic: vInPubTopic,
			}
			go func() {
				dec.decodeMessages()
			}()

			bDIn, err := proto.Marshal(lTest.inpDIn)
			require.NoError(t, err)
			t.Logf("bDIn: %s", bDIn)

			require.NoError(t, dInQueue.Publish("", bDIn))

			for _, res := range lTest.res {
				select {
				case msg := <-vInSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, vInPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Fatalf("\nExpect: %+v\nActual: %+v", res, vIn)
					}
				case <-time.After(2 * time.Second):
					t.Fatal("Message timed out")
				}
			}
		})
	}
}

func TestDecodeMessagesError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpDIn     *message.DecoderIn
		inpDecoder api.Decoder
		inpErr     error
		inpTimes   int
	}{
		// Empty data.
		{&message.DecoderIn{}, api.Decoder_RAW, nil, 1},
		// Bad payload.
		{nil, api.Decoder_RAW, nil, 0},
		// Device not found.
		{&message.DecoderIn{}, api.Decoder_RAW, dao.ErrNotFound, 1},
		// Devicer error.
		{&message.DecoderIn{}, api.Decoder_RAW, errTestProc, 1},
		// Decode error.
		{&message.DecoderIn{}, api.Decoder(999), nil, 1},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			dev := random.Device("dec", uuid.NewString())
			dev.Decoder = lTest.inpDecoder

			dInQueue := queue.NewFake()
			dInSub, err := dInQueue.Subscribe("")
			require.NoError(t, err)

			decQueue := queue.NewFake()
			vInSub, err := decQueue.Subscribe("")
			require.NoError(t, err)

			devicer := NewMockdevicer(gomock.NewController(t))
			devicer.EXPECT().ReadByUniqID(gomock.Any(), gomock.Any()).
				Return(dev, lTest.inpErr).Times(lTest.inpTimes)

			dec := Decoder{
				devDAO: devicer,
				reg:    registry.New(),

				decQueue:    decQueue,
				dInSub:      dInSub,
				vInPubTopic: "topic-" + random.String(10),
			}
			go func() {
				dec.decodeMessages()
			}()

			bDIn := []byte("dec-aaa")
			if lTest.inpDIn != nil {
				var err error
				bDIn, err = proto.Marshal(lTest.inpDIn)
				require.NoError(t, err)
				t.Logf("bDIn: %s", bDIn)
			}

			require.NoError(t, dInQueue.Publish("", bDIn))

			select {
			case msg := <-vInSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(100 * time.Millisecond):
				// Successful timeout without publish (normally 0.02s).
			}
		})
	}
}
