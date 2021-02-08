// +build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateMessages(t *testing.T) {
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()
	boolVal := &common.DataPoint_BoolVal{BoolVal: []bool{true,
		false}[random.Intn(2)]}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("val"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := random.Device("val", createOrg.Id)
	dev.Status = api.Status_ACTIVE
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	tests := []struct {
		inpVIn *message.ValidatorIn
		res    *message.ValidatorOut
	}{
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: createDev.UniqId,
			Attr: "val-motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: now, Token: createDev.Token, TraceId: traceID},
			OrgId: createOrg.Id},
			&message.ValidatorOut{Point: &common.DataPoint{
				UniqId: createDev.UniqId, Attr: "val-motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
				Token: createDev.Token, TraceId: traceID}, OrgId: createOrg.Id,
				DevId: createDev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: createDev.UniqId,
			Attr: "val-temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
			Ts: now, Token: createDev.Token, TraceId: traceID},
			OrgId: createOrg.Id},
			&message.ValidatorOut{Point: &common.DataPoint{
				UniqId: createDev.UniqId, Attr: "val-temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}, Ts: now,
				Token: createDev.Token, TraceId: traceID}, OrgId: createOrg.Id,
				DevId: createDev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: createDev.UniqId,
			Attr: "val-power", ValOneof: &common.DataPoint_StrVal{
				StrVal: "batt"}, Ts: now, Token: createDev.Token,
			TraceId: traceID}, OrgId: createOrg.Id},
			&message.ValidatorOut{Point: &common.DataPoint{
				UniqId: createDev.UniqId, Attr: "val-power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"}, Ts: now,
				Token: createDev.Token, TraceId: traceID}, OrgId: createOrg.Id,
				DevId: createDev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: createDev.UniqId,
			Attr: "val-leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
			OrgId: createOrg.Id, SkipToken: true},
			&message.ValidatorOut{Point: &common.DataPoint{
				UniqId: createDev.UniqId, Attr: "val-leak",
				ValOneof: boolVal, Ts: now, TraceId: traceID},
				OrgId: createOrg.Id, DevId: createDev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: createDev.UniqId,
			Attr: "val-leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
			SkipToken: true},
			&message.ValidatorOut{Point: &common.DataPoint{
				UniqId: createDev.UniqId, Attr: "val-leak",
				ValOneof: boolVal, Ts: now, TraceId: traceID},
				OrgId: createOrg.Id, DevId: createDev.Id}},
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
				t.Logf("vOut: %+v", vOut)

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, vOut) {
					t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, vOut)
				}
			case <-time.After(testTimeout):
				t.Fatal("Message timed out")
			}
		})
	}
}

func TestValidateMessagesError(t *testing.T) {
	uniqID := "val-" + random.String(16)
	disUniqID := "val-" + random.String(16)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("val"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := random.Device("val", createOrg.Id)
	dev.Status = api.Status_ACTIVE
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	disDev := random.Device("val", createOrg.Id)
	dev.Status = api.Status_DISABLED
	createDisDev, err := globalDevDAO.Create(ctx, disDev)
	t.Logf("createDisDev, err: %+v, %v", createDisDev, err)
	require.NoError(t, err)

	tests := []struct {
		inpVIn *message.ValidatorIn
	}{
		// Bad payload.
		{nil},
		// Device not found.
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16)}}},
		// Missing value.
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: random.String(10)}}},
		// Invalid org ID.
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: random.String(10), ValOneof: &common.DataPoint_IntVal{}},
			OrgId: "val-aaa"}},
		// Device disabled.
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: disUniqID,
			Attr: random.String(10), ValOneof: &common.DataPoint_IntVal{}},
			OrgId: createOrg.Id}},
		// Invalid token.
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: random.String(10), ValOneof: &common.DataPoint_IntVal{},
			Token: "val-aaa"}, OrgId: createOrg.Id}},
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
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
