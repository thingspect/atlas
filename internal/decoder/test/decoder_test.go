// +build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testTimeout = 6 * time.Second

func TestDecodeMessages(t *testing.T) {
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dec"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := random.Device("dec", createOrg.Id)
	dev.Status = common.Status_ACTIVE
	dev.Decoder = common.Decoder_RADIO_BRIDGE_DOOR_V2
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.DecoderIn
		res []*message.ValidatorIn
	}{
		{&message.DecoderIn{UniqId: dev.UniqId, Data: []byte{0x19, 0x03, 0x01},
			Ts: now, TraceId: traceID}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: dev.UniqId, Attr: "count",
				ValOneof: &common.DataPoint_IntVal{IntVal: 9}, Ts: now,
				TraceId: traceID}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: dev.UniqId, Attr: "open",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: true}, Ts: now,
				TraceId: traceID}, SkipToken: true},
		}},
		{&message.DecoderIn{UniqId: dev.UniqId, Data: []byte{0x1a, 0x03, 0x00},
			Ts: now, TraceId: traceID}, []*message.ValidatorIn{
			{Point: &common.DataPoint{UniqId: dev.UniqId, Attr: "count",
				ValOneof: &common.DataPoint_IntVal{IntVal: 10}, Ts: now,
				TraceId: traceID}, SkipToken: true},
			{Point: &common.DataPoint{UniqId: dev.UniqId, Attr: "open",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
				Ts:       now, TraceId: traceID}, SkipToken: true},
		}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can decode %+v", lTest), func(t *testing.T) {
			bDIn, err := proto.Marshal(lTest.inp)
			require.NoError(t, err)
			t.Logf("bDIn: %s", bDIn)

			require.NoError(t, globalDecQueue.Publish(globalDInSubTopic, bDIn))

			// Don't stop the flow of execution (assert) to avoid leaving
			// messages orphaned in the queue.
			for _, res := range lTest.res {
				select {
				case msg := <-globalVInSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					assert.Equal(t, globalVInPubTopic, msg.Topic())

					vIn := &message.ValidatorIn{}
					assert.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
					t.Logf("vIn: %+v", vIn)

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, vIn) {
						t.Errorf("\nExpect: %+v\nActual: %+v", res, vIn)
					}
				case <-time.After(testTimeout):
					t.Error("Message timed out")
				}
			}
		})
	}
}

func TestDecodeMessagesError(t *testing.T) {
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dec"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := random.Device("dec", createOrg.Id)
	dev.Status = common.Status_ACTIVE
	dev.Decoder = common.Decoder_RAW
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	invDev := random.Device("dec", createOrg.Id)
	invDev.Status = common.Status_ACTIVE
	invDev.Decoder = common.Decoder(999)
	createInvDev, err := globalDevDAO.Create(ctx, invDev)
	t.Logf("createInvDev, err: %+v, %v", createInvDev, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.DecoderIn
	}{
		// Empty data.
		{&message.DecoderIn{UniqId: createDev.UniqId, Ts: now,
			TraceId: traceID}},
		// Bad payload.
		{nil},
		// Device not found.
		{&message.DecoderIn{UniqId: random.String(16), Ts: now,
			TraceId: traceID}},
		// Decode error, defaults to Decoder zero value when not in registry.
		{&message.DecoderIn{UniqId: createInvDev.UniqId, Ts: now,
			TraceId: traceID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot decode %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bDIn := []byte("dec-aaa")
			if lTest.inp != nil {
				var err error
				bDIn, err = proto.Marshal(lTest.inp)
				require.NoError(t, err)
				t.Logf("bDIn: %s", bDIn)
			}

			require.NoError(t, globalDecQueue.Publish(globalDInSubTopic, bDIn))

			select {
			case msg := <-globalVInSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
