//go:build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testTimeout = 6 * time.Second

func TestDecodeMessages(t *testing.T) {
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dec"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	doorDev := random.Device("dec", createOrg.GetId())
	doorDev.Status = api.Status_ACTIVE
	doorDev.Decoder = api.Decoder_RADIO_BRIDGE_DOOR_V2
	createDoorDev, err := globalDevDAO.Create(ctx, doorDev)
	t.Logf("createDoorDev, err: %+v, %v", createDoorDev, err)
	require.NoError(t, err)

	co2Dev := random.Device("dec", createOrg.GetId())
	co2Dev.Status = api.Status_ACTIVE
	co2Dev.Decoder = api.Decoder_GLOBALSAT_CO2
	createCO2Dev, err := globalDevDAO.Create(ctx, co2Dev)
	t.Logf("createCO2Dev, err: %+v, %v", createCO2Dev, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.DecoderIn
		res []*message.ValidatorIn
	}{
		{
			&message.DecoderIn{
				UniqId: doorDev.GetUniqId(), Data: []byte{0x19, 0x03, 0x01}, Ts: now,
				TraceId: traceID[:],
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: doorDev.GetUniqId(), Attr: "count",
						ValOneof: &common.DataPoint_IntVal{IntVal: 9}, Ts: now,
						TraceId: traceID.String(),
					}, SkipToken: true,
				}, {
					Point: &common.DataPoint{
						UniqId: doorDev.GetUniqId(), Attr: "open",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: true},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
			},
		},
		{
			&message.DecoderIn{
				UniqId: doorDev.GetUniqId(), Data: []byte{0x1a, 0x03, 0x00}, Ts: now,
				TraceId: traceID[:],
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: doorDev.GetUniqId(), Attr: "count",
						ValOneof: &common.DataPoint_IntVal{IntVal: 10}, Ts: now,
						TraceId: traceID.String(),
					}, SkipToken: true,
				}, {
					Point: &common.DataPoint{
						UniqId: doorDev.GetUniqId(), Attr: "open",
						ValOneof: &common.DataPoint_BoolVal{BoolVal: false},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
			},
		},
		{
			&message.DecoderIn{
				UniqId: co2Dev.GetUniqId(), Data: []byte{
					0x01, 0x09, 0x61, 0x13, 0x95, 0x02, 0x92,
				}, Ts: now, TraceId: traceID[:],
			}, []*message.ValidatorIn{
				{
					Point: &common.DataPoint{
						UniqId: co2Dev.GetUniqId(), Attr: "temp_c",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 24},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: co2Dev.GetUniqId(), Attr: "temp_f",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 75.2},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: co2Dev.GetUniqId(), Attr: "humidity_pct",
						ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 50.13},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
				{
					Point: &common.DataPoint{
						UniqId: co2Dev.GetUniqId(), Attr: "co2_ppm",
						ValOneof: &common.DataPoint_IntVal{IntVal: 658},
						Ts:       now, TraceId: traceID.String(),
					}, SkipToken: true,
				},
			},
		},
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
	traceID := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dec"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := random.Device("dec", createOrg.GetId())
	dev.Status = api.Status_ACTIVE
	dev.Decoder = api.Decoder_RAW
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	invDev := random.Device("dec", createOrg.GetId())
	invDev.Status = api.Status_ACTIVE
	invDev.Decoder = api.Decoder(999)
	createInvDev, err := globalDevDAO.Create(ctx, invDev)
	t.Logf("createInvDev, err: %+v, %v", createInvDev, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.DecoderIn
	}{
		// Empty data.
		{
			&message.DecoderIn{
				UniqId: createDev.GetUniqId(), Ts: now, TraceId: traceID[:],
			},
		},
		// Bad payload.
		{nil},
		// Device not found.
		{
			&message.DecoderIn{
				UniqId: random.String(16), Ts: now, TraceId: traceID[:],
			},
		},
		// Decode error, defaults to Decoder zero value when not in registry.
		{
			&message.DecoderIn{
				UniqId: createInvDev.GetUniqId(), Ts: now, TraceId: traceID[:],
			},
		},
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
