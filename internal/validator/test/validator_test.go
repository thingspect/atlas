// +build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
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
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: now, Token: createDev.Token}, OrgId: createOrg.ID,
			TraceId: traceID},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: now, Token: createDev.Token}, OrgId: createOrg.ID,
				TraceId: traceID, DevId: createDev.ID}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
			Ts: now, Token: createDev.Token}, OrgId: createOrg.ID,
			TraceId: traceID},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Ts: now, Token: createDev.Token}, OrgId: createOrg.ID,
				TraceId: traceID, DevId: createDev.ID}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: "metadata", MapVal: map[string]string{"val-aaa": "val-bbb"},
			Ts: now, Token: createDev.Token}, OrgId: createOrg.ID,
			TraceId: traceID},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: uniqID,
				Attr:   "metadata",
				MapVal: map[string]string{"val-aaa": "val-bbb"}, Ts: now,
				Token: createDev.Token}, OrgId: createOrg.ID, TraceId: traceID,
				DevId: createDev.ID}},
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
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16)}}},
		// Invalid org ID.
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID},
			OrgId: "val-aaa"}},
		// Device disabled.
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: disabledUniqID},
			OrgId: createOrg.ID}},
		// Invalid token.
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Token: "val-aaa"}, OrgId: createOrg.ID}},
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
