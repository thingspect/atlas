// +build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateMessages(t *testing.T) {
	uniqID := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := device.Device{OrgID: createOrg.ID, UniqID: uniqID}
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	tests := []struct {
		inpVIn *message.ValidatorIn
		res    *message.ValidatorOut
	}{
		{&message.ValidatorIn{UniqId: uniqID, Attr: "motion",
			ValOneof: &message.ValidatorIn_IntVal{IntVal: 123}, Ts: now,
			Token: createDev.Token, OrgId: createOrg.ID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "motion",
				ValOneof: &message.ValidatorOut_IntVal{IntVal: 123}, Ts: now,
				DevId: createDev.ID, OrgId: createOrg.ID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "temp",
			ValOneof: &message.ValidatorIn_Fl64Val{Fl64Val: 20.3}, Ts: now,
			Token: createDev.Token, OrgId: createOrg.ID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "temp",
				ValOneof: &message.ValidatorOut_Fl64Val{Fl64Val: 20.3}, Ts: now,
				DevId: createDev.ID, OrgId: createOrg.ID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "power",
			ValOneof: &message.ValidatorIn_StrVal{StrVal: "batt"}, Ts: now,
			Token: createDev.Token, OrgId: createOrg.ID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "power",
				ValOneof: &message.ValidatorOut_StrVal{StrVal: "batt"}, Ts: now,
				DevId: createDev.ID, OrgId: createOrg.ID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "leak",
			ValOneof: &message.ValidatorIn_BoolVal{BoolVal: true}, Ts: now,
			Token: createDev.Token, OrgId: createOrg.ID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "leak",
				ValOneof: &message.ValidatorOut_BoolVal{BoolVal: true}, Ts: now,
				DevId: createDev.ID, OrgId: createOrg.ID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "raw",
			ValOneof: &message.ValidatorIn_BytesVal{BytesVal: []byte{0x00}},
			Ts:       now, Token: createDev.Token, OrgId: createOrg.ID,
			TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "raw",
				ValOneof: &message.ValidatorOut_BytesVal{
					BytesVal: []byte{0x00}}, Ts: now, DevId: createDev.ID,
				OrgId: createOrg.ID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "metadata",
			MapVal: map[string]string{"val-aaa": "val-bbb"}, Ts: now,
			Token: createDev.Token, OrgId: createOrg.ID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "metadata",
				MapVal: map[string]string{"val-aaa": "val-bbb"}, Ts: now,
				DevId: createDev.ID, OrgId: createOrg.ID, TraceId: traceID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can validate %+v", lTest), func(t *testing.T) {
			bVIn, err := proto.Marshal(lTest.inpVIn)
			require.NoError(t, err)
			t.Logf("bVIn: %s", bVIn)

			require.NoError(t, globalVInQueue.Publish(globalVInSubTopic, bVIn))

			select {
			case msg := <-globalVOutSub.C():
				msg.Ack()
				t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
					msg.Payload())
				require.Equal(t, globalVOutPubTopic, msg.Topic())

				vOut := &message.ValidatorOut{}
				require.NoError(t, proto.Unmarshal(msg.Payload(), vOut))
				t.Logf("vOut: %#v", vOut)

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, vOut) {
					t.Fatalf("Expected, actual: %#v, %#v", lTest.res, vOut)
				}
			case <-time.After(5 * time.Second):
				t.Error("Message timed out")
			}
		})
	}
}

func TestValidateMessagesError(t *testing.T) {
	uniqID := random.String(16)
	disabledUniqID := random.String(16)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	org := org.Org{Name: random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := device.Device{OrgID: createOrg.ID, UniqID: uniqID}
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	disabledDev := device.Device{OrgID: createOrg.ID, UniqID: disabledUniqID,
		Disabled: true}
	createDisabledDev, err := globalDevDAO.Create(ctx, disabledDev)
	t.Logf("createDisabledDev, err: %+v, %v", createDisabledDev, err)
	require.NoError(t, err)

	tests := []struct {
		inpVIn *message.ValidatorIn
	}{
		// Bad payload.
		{nil},
		// Device not found.
		{&message.ValidatorIn{UniqId: random.String(16)}},
		// Invalid org ID.
		{&message.ValidatorIn{UniqId: uniqID, OrgId: "val-aaa"}},
		// Device disabled.
		{&message.ValidatorIn{UniqId: disabledUniqID, OrgId: createOrg.ID}},
		// Invalid token.
		{&message.ValidatorIn{UniqId: uniqID, Token: "val-aaa",
			OrgId: createOrg.ID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot validate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bVIn := []byte("val-aaa")
			if lTest.inpVIn != nil {
				var err error
				bVIn, err = proto.Marshal(lTest.inpVIn)
				require.NoError(t, err)
				t.Logf("bVIn: %s", bVIn)
			}

			require.NoError(t, globalVInQueue.Publish(globalVInSubTopic, bVIn))

			select {
			case msg := <-globalVOutSub.C():
				t.Errorf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
