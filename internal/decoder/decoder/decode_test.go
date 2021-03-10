// +build !integration

package decoder

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/decode/registry"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errTestProc = errors.New("decoder: test processor error")

func TestDecodeMessages(t *testing.T) {
	t.Parallel()

	uniqID := "dec-" + random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()

	tests := []struct {
		inpDIn     *message.DecoderIn
		inpDecoder common.Decoder
		res        []*message.ValidatorIn
	}{
		{&message.DecoderIn{UniqId: uniqID, Data: []byte{0x19, 0x03, 0x01},
			Ts: now, TraceId: traceID}, common.Decoder_RADIO_BRIDGE_DOOR_V1,
			[]*message.ValidatorIn{
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 9}, Ts: now,
					TraceId: traceID}, SkipToken: true},
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "open",
					ValOneof: &common.DataPoint_BoolVal{BoolVal: true}, Ts: now,
					TraceId: traceID}, SkipToken: true},
			}},
		{&message.DecoderIn{UniqId: uniqID, Data: []byte{0x1a, 0x03, 0x00},
			Ts: now, TraceId: traceID}, common.Decoder_RADIO_BRIDGE_DOOR_V2,
			[]*message.ValidatorIn{
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 10}, Ts: now,
					TraceId: traceID}, SkipToken: true},
				{Point: &common.DataPoint{UniqId: uniqID, Attr: "open",
					ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
					Ts:       now, TraceId: traceID}, SkipToken: true},
			}},
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
			devicer.EXPECT().ReadByUniqID(gomock.Any(), lTest.inpDIn.UniqId).
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
		inpDecoder common.Decoder
		inpErr     error
		inpTimes   int
	}{
		// Empty data.
		{&message.DecoderIn{}, common.Decoder_RAW, nil, 1},
		// Bad payload.
		{nil, common.Decoder_RAW, nil, 0},
		// Device not found.
		{&message.DecoderIn{}, common.Decoder_RAW, dao.ErrNotFound, 1},
		// Devicer error.
		{&message.DecoderIn{}, common.Decoder_RAW, errTestProc, 1},
		// Decode error.
		{&message.DecoderIn{}, common.Decoder(999), nil, 1},
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
