// +build !integration

package eventer

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const errTestProc consterr.Error = "eventer: test processor error"

func TestEventMessages(t *testing.T) {
	t.Parallel()

	orgID := uuid.NewString()
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()
	ruleID := uuid.NewString()

	tests := []struct {
		inpVOut  *message.ValidatorOut
		inpRules []*common.Rule
		inpTimes int
		res      []*message.EventerOut
	}{
		{&message.ValidatorOut{Point: &common.DataPoint{Attr: "motion", Ts: now,
			TraceId: traceID}}, []*common.Rule{{Id: ruleID, Expr: `true`}}, 1,
			[]*message.EventerOut{{Point: &common.DataPoint{Attr: "motion",
				Ts: now, TraceId: traceID}, Rule: &common.Rule{Id: ruleID,
				Expr: `true`}}}},
		{&message.ValidatorOut{Point: &common.DataPoint{Attr: "temp", Ts: now,
			TraceId: traceID}}, []*common.Rule{{Id: ruleID, Expr: `true`},
			{Id: ruleID, Expr: `true`}}, 2, []*message.EventerOut{{
			Point: &common.DataPoint{Attr: "temp", Ts: now, TraceId: traceID},
			Rule:  &common.Rule{Id: ruleID, Expr: `true`}},
			{Point: &common.DataPoint{Attr: "temp", Ts: now, TraceId: traceID},
				Rule: &common.Rule{Id: ruleID, Expr: `true`}}}},
		{&message.ValidatorOut{Point: &common.DataPoint{Attr: "power", Ts: now,
			TraceId: traceID}}, nil, 0, nil},
		{&message.ValidatorOut{Point: &common.DataPoint{Attr: "leak", Ts: now,
			TraceId: traceID}}, []*common.Rule{{Id: ruleID, Expr: `false`}}, 0,
			nil},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can event %+v", lTest), func(t *testing.T) {
			t.Parallel()

			dev := random.Device("ev", uuid.NewString())
			dev.OrgId = orgID
			lTest.inpVOut.Device = dev

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			evQueue := queue.NewFake()
			vInSub, err := evQueue.Subscribe("")
			require.NoError(t, err)
			eOutPubTopic := "topic-" + random.String(10)

			ruler := NewMockruler(gomock.NewController(t))
			ruler.EXPECT().ListByTags(gomock.Any(), lTest.inpVOut.Device.OrgId,
				lTest.inpVOut.Point.Attr, lTest.inpVOut.Device.Tags).
				Return(lTest.inpRules, nil).Times(1)

			// Reuse ruleID for less branching in the mocking paths.
			event := &api.Event{OrgId: lTest.inpVOut.Device.OrgId,
				RuleId: ruleID, UniqId: dev.UniqId, CreatedAt: now,
				TraceId: traceID}
			eventer := NewMockeventer(gomock.NewController(t))
			eventer.EXPECT().Create(gomock.Any(),
				matcher.NewProtoMatcher(event)).Return(nil).
				Times(lTest.inpTimes)

			ev := Eventer{
				ruleDAO:  ruler,
				eventDAO: eventer,

				evQueue:      evQueue,
				vOutSub:      vOutSub,
				eOutPubTopic: eOutPubTopic,
			}
			go func() {
				ev.eventMessages()
			}()

			bVOut, err := proto.Marshal(lTest.inpVOut)
			require.NoError(t, err)
			t.Logf("bVOut: %s", bVOut)

			require.NoError(t, vOutQueue.Publish("", bVOut))

			for _, res := range lTest.res {
				select {
				case msg := <-vInSub.C():
					msg.Ack()
					t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
						msg.Payload())
					require.Equal(t, eOutPubTopic, msg.Topic())

					eOut := &message.EventerOut{}
					require.NoError(t, proto.Unmarshal(msg.Payload(), eOut))
					t.Logf("eOut: %+v", eOut)

					// Normalize device.
					res.Device = dev

					// Testify does not currently support protobuf equality:
					// https://github.com/stretchr/testify/issues/758
					if !proto.Equal(res, eOut) {
						t.Fatalf("\nExpect: %+v\nActual: %+v", res, eOut)
					}
				case <-time.After(2 * time.Second):
					t.Fatal("Message timed out")
				}
			}

			if len(lTest.res) == 0 {
				select {
				case msg := <-vInSub.C():
					t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, "+
						"%s", msg.Topic(), msg.Payload())
				case <-time.After(100 * time.Millisecond):
					// Successful timeout without publish (normally 0.02s).
				}
			}
		})
	}
}

func TestEventMessagesError(t *testing.T) {
	t.Parallel()

	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	tests := []struct {
		inpVOut         *message.ValidatorOut
		inpRulerErr     error
		inpRulerTimes   int
		inpRules        []*common.Rule
		inpEventerErr   error
		inpEventerTimes int
	}{
		// Bad payload.
		{nil, nil, 0, nil, nil, 0},
		// Missing data point.
		{&message.ValidatorOut{}, nil, 0, nil, nil, 0},
		// Missing device.
		{&message.ValidatorOut{Point: &common.DataPoint{}}, nil, 0, nil, nil,
			0},
		// Ruler error.
		{&message.ValidatorOut{Point: &common.DataPoint{},
			Device: &common.Device{}}, errTestProc, 1, nil, nil, 0},
		// Eval error.
		{&message.ValidatorOut{Point: &common.DataPoint{Ts: now},
			Device: &common.Device{}}, nil, 1,
			[]*common.Rule{{Expr: `1 + "aaa"`}}, nil, 0},
		// Eventer already exists.
		{&message.ValidatorOut{Point: &common.DataPoint{Ts: now},
			Device: &common.Device{}}, nil, 1, []*common.Rule{{Expr: `true`}},
			dao.ErrAlreadyExists, 1},
		// Eventer error.
		{&message.ValidatorOut{Point: &common.DataPoint{Ts: now},
			Device: &common.Device{}}, nil, 1, []*common.Rule{{Expr: `true`}},
			errTestProc, 1},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot event %+v", lTest), func(t *testing.T) {
			t.Parallel()

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			evQueue := queue.NewFake()
			vInSub, err := evQueue.Subscribe("")
			require.NoError(t, err)
			eOutPubTopic := "topic-" + random.String(10)

			ruler := NewMockruler(gomock.NewController(t))
			ruler.EXPECT().ListByTags(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(lTest.inpRules, lTest.inpRulerErr).
				Times(lTest.inpRulerTimes)

			eventer := NewMockeventer(gomock.NewController(t))
			eventer.EXPECT().Create(gomock.Any(), gomock.Any()).
				Return(lTest.inpEventerErr).Times(lTest.inpEventerTimes)

			ev := Eventer{
				ruleDAO:  ruler,
				eventDAO: eventer,

				evQueue:      evQueue,
				vOutSub:      vOutSub,
				eOutPubTopic: eOutPubTopic,
			}
			go func() {
				ev.eventMessages()
			}()

			bVOut := []byte("ev-aaa")
			if lTest.inpVOut != nil {
				bVOut, err = proto.Marshal(lTest.inpVOut)
				require.NoError(t, err)
				t.Logf("bVOut: %s", bVOut)
			}

			require.NoError(t, vOutQueue.Publish("", bVOut))

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
