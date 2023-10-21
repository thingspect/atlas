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

const testTimeout = 12 * time.Second

func TestEventMessages(t *testing.T) {
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("ev"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	singleDev, err := globalDevDAO.Create(ctx, random.Device("ev",
		createOrg.GetId()))
	t.Logf("singleDev, err: %+v, %v", singleDev, err)
	require.NoError(t, err)

	sRule := random.Rule("ev", createOrg.GetId())
	sRule.Status = api.Status_ACTIVE
	sRule.DeviceTag = singleDev.GetTags()[0]
	sRule.Attr = "ev-motion"
	sRule.Expr = `true`
	singleRule, err := globalRuleDAO.Create(ctx, sRule)
	t.Logf("singleRule, err: %+v, %v", singleRule, err)
	require.NoError(t, err)

	dDev := random.Device("ev", createOrg.GetId())
	dDev.Tags = random.Tags("ev", 2)
	doubleDev, err := globalDevDAO.Create(ctx, dDev)
	t.Logf("doubleDev, err: %+v, %v", doubleDev, err)
	require.NoError(t, err)

	dRule1 := random.Rule("ev", createOrg.GetId())
	dRule1.Status = api.Status_ACTIVE
	dRule1.DeviceTag = doubleDev.GetTags()[0]
	dRule1.Attr = "ev-temp"
	dRule1.Expr = `true`
	doubleRule1, err := globalRuleDAO.Create(ctx, dRule1)
	t.Logf("doubleRule1, err: %+v, %v", doubleRule1, err)
	require.NoError(t, err)

	dRule2, _ := proto.Clone(dRule1).(*api.Rule)
	dRule2.DeviceTag = doubleDev.GetTags()[1]
	doubleRule2, err := globalRuleDAO.Create(ctx, dRule2)
	t.Logf("doubleRule2, err: %+v, %v", doubleRule2, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.ValidatorOut
		res []*message.EventerOut
	}{
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{
					Attr: "ev-motion", Ts: now, TraceId: traceID,
				}, Device: singleDev,
			}, []*message.EventerOut{
				{Point: &common.DataPoint{
					Attr: "ev-motion", Ts: now, TraceId: traceID,
				}, Device: singleDev, Rule: singleRule},
			},
		},
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{
					Attr: "ev-temp", Ts: now, TraceId: traceID,
				}, Device: doubleDev,
			}, []*message.EventerOut{
				{Point: &common.DataPoint{
					Attr: "ev-temp", Ts: now, TraceId: traceID,
				}, Device: doubleDev, Rule: doubleRule1}, {
					Point: &common.DataPoint{
						Attr: "ev-temp", Ts: now, TraceId: traceID,
					}, Device: doubleDev, Rule: doubleRule2,
				},
			},
		},
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{
					Attr: "ev-power", Ts: now, TraceId: traceID,
				}, Device: singleDev,
			}, nil,
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can event %+v", lTest), func(t *testing.T) {
			bVOut, err := proto.Marshal(lTest.inp)
			require.NoError(t, err)
			t.Logf("bVOut: %s", bVOut)

			require.NoError(t, globalEvQueue.Publish(globalVOutSubTopic, bVOut))

			// Don't stop the flow of execution (assert) to avoid leaving
			// messages orphaned in the queue.
			for _, res := range lTest.res {
				t.Logf("DEBUG res: %+v", res)
				select {
				case msg := <-globalEOutSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					assert.Equal(t, globalEOutPubTopic, msg.Topic())

					eOut := &message.EventerOut{}
					assert.NoError(t, proto.Unmarshal(msg.Payload(), eOut))
					t.Logf("eOut: %+v", eOut)

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, eOut) {
						t.Errorf("\nExpect: %+v\nActual: %+v", res, eOut)
					}
				case <-time.After(testTimeout):
					t.Error("Message timed out")
				}

				// Verify events by rule ID.
				event := &api.Event{
					OrgId: createOrg.GetId(), RuleId: res.GetRule().GetId(),
					UniqId: lTest.inp.GetDevice().GetUniqId(), CreatedAt: timestamppb.New(
						now.AsTime().Truncate(time.Millisecond)),
					TraceId: traceID,
				}

				ctx, cancel := context.WithTimeout(context.Background(),
					testTimeout)
				defer cancel()

				listEvents, err := globalEvDAO.List(ctx, createOrg.GetId(),
					lTest.inp.GetDevice().GetUniqId(), "", res.GetRule().GetId(), now.AsTime(),
					now.AsTime().Add(-time.Millisecond))
				t.Logf("listEvents, err: %+v, %v", listEvents, err)
				assert.NoError(t, err)
				assert.Len(t, listEvents, 1)

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(event, listEvents[0]) {
					t.Errorf("\nExpect: %+v\nActual: %+v", event, listEvents[0])
				}
			}

			if len(lTest.res) == 0 {
				select {
				case msg := <-globalEOutSub.C():
					t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, "+
						"%s", msg.Topic(), msg.Payload())
				case <-time.After(500 * time.Millisecond):
					// Successful timeout without publish (normally 0.25s).
				}
			}
		})
	}
}

func TestEventMessagesError(t *testing.T) {
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("ev"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createDev, err := globalDevDAO.Create(ctx, random.Device("ev", createOrg.GetId()))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	rule := random.Rule("ev", createOrg.GetId())
	rule.Status = api.Status_ACTIVE
	rule.DeviceTag = createDev.GetTags()[0]
	rule.Attr = "ev-motion"
	rule.Expr = `1 + "aaa"`
	createRule, err := globalRuleDAO.Create(ctx, rule)
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.ValidatorOut
	}{
		// Bad payload.
		{nil},
		// Missing data point.
		{
			&message.ValidatorOut{Device: &api.Device{Id: createDev.GetId()}},
		},
		// Missing device.
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{UniqId: createDev.GetUniqId()},
			},
		},
		// Eval error.
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{
					Attr: "ev-motion", Ts: now, TraceId: uuid.NewString(),
				}, Device: createDev,
			},
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot event %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bVOut := []byte("ev-aaa")
			if lTest.inp != nil {
				var err error
				bVOut, err = proto.Marshal(lTest.inp)
				require.NoError(t, err)
				t.Logf("bVOut: %s", bVOut)
			}

			require.NoError(t, globalEvQueue.Publish(globalVOutSubTopic, bVOut))

			select {
			case msg := <-globalEOutSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(500 * time.Millisecond):
				// Successful timeout without publish (normally 0.25s).
			}
		})
	}
}
