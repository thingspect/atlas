// +build !integration

package validator

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errTestProc = errors.New("validator: test processor error")

func TestValidateMessages(t *testing.T) {
	t.Parallel()

	dev := random.Device("val", uuid.NewString())
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	dev.Status = api.Status_ACTIVE
	traceID := uuid.NewString()
	boolVal := &common.DataPoint_BoolVal{BoolVal: []bool{true,
		false}[random.Intn(2)]}

	tests := []struct {
		inp *message.ValidatorIn
		res *message.ValidatorOut
	}{
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: dev.UniqId,
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: now, Token: dev.Token, TraceId: traceID}, OrgId: dev.OrgId},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: dev.UniqId,
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: now, Token: dev.Token, TraceId: traceID}, OrgId: dev.OrgId,
				DevId: dev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: dev.UniqId,
			Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
			Ts: now, Token: dev.Token, TraceId: traceID}, OrgId: dev.OrgId},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: dev.UniqId,
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Ts: now, Token: dev.Token, TraceId: traceID}, OrgId: dev.OrgId,
				DevId: dev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: dev.UniqId,
			Attr: "power", ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
			Ts: now, Token: dev.Token, TraceId: traceID}, OrgId: dev.OrgId},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: dev.UniqId,
				Attr: "power", ValOneof: &common.DataPoint_StrVal{
					StrVal: "batt"}, Ts: now, Token: dev.Token,
				TraceId: traceID}, OrgId: dev.OrgId, DevId: dev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: dev.UniqId,
			Attr: "leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
			OrgId: dev.OrgId, SkipToken: true},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: dev.UniqId,
				Attr: "leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
				OrgId: dev.OrgId, DevId: dev.Id}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: dev.UniqId,
			Attr: "leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
			SkipToken: true},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: dev.UniqId,
				Attr: "leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
				OrgId: dev.OrgId, DevId: dev.Id}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can validate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			vInQueue := queue.NewFake()
			vInSub, err := vInQueue.Subscribe("")
			require.NoError(t, err)

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)
			vOutPubTopic := "topic-" + random.String(10)

			devicer := NewMockdevicer(gomock.NewController(t))
			devicer.EXPECT().ReadByUniqID(gomock.Any(), lTest.inp.Point.UniqId).
				Return(dev, nil).Times(1)

			val := Validator{
				devDAO:       devicer,
				vInSub:       vInSub,
				vOutQueue:    vOutQueue,
				vOutPubTopic: vOutPubTopic,
			}
			go func() {
				val.validateMessages()
			}()

			bVIn, err := proto.Marshal(lTest.inp)
			require.NoError(t, err)
			t.Logf("bVIn: %s", bVIn)

			require.NoError(t, vInQueue.Publish("", bVIn))

			select {
			case msg := <-vOutSub.C():
				msg.Ack()
				t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
					msg.Payload())
				require.Equal(t, vOutPubTopic, msg.Topic())

				vOut := &message.ValidatorOut{}
				require.NoError(t, proto.Unmarshal(msg.Payload(), vOut))
				t.Logf("vOut: %+v", vOut)

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, vOut) {
					t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.res, vOut)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("Message timed out")
			}
		})
	}
}

func TestValidateMessagesError(t *testing.T) {
	t.Parallel()

	dev := random.Device("val", uuid.NewString())

	tests := []struct {
		inpVIn    *message.ValidatorIn
		inpStatus api.Status
		inpErr    error
		inpTimes  int
	}{
		// Bad payload.
		{nil, api.Status_ACTIVE, nil, 0},
		// Missing data point.
		{&message.ValidatorIn{}, api.Status_ACTIVE, nil, 0},
		// Device not found.
		{&message.ValidatorIn{Point: &common.DataPoint{}}, api.Status_ACTIVE,
			dao.ErrNotFound, 1},
		// Devicer error.
		{&message.ValidatorIn{Point: &common.DataPoint{}}, api.Status_ACTIVE,
			errTestProc, 1},
		// Missing value.
		{&message.ValidatorIn{Point: &common.DataPoint{}}, api.Status_ACTIVE,
			nil, 1},
		// Invalid org ID.
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: random.String(10),
			ValOneof: &common.DataPoint_IntVal{}}, OrgId: "val-aaa"},
			api.Status_ACTIVE, nil, 1},
		// Device status.
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: random.String(10),
			ValOneof: &common.DataPoint_IntVal{}}, OrgId: dev.OrgId},
			api.Status_DISABLED, nil, 1},
		// Invalid token.
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: random.String(10),
			ValOneof: &common.DataPoint_IntVal{}, Token: "val-aaa"},
			OrgId: dev.OrgId}, api.Status_ACTIVE, nil, 1},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot validate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			lDev := proto.Clone(dev).(*api.Device)
			lDev.Status = lTest.inpStatus

			vInQueue := queue.NewFake()
			vInSub, err := vInQueue.Subscribe("")
			require.NoError(t, err)

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			devicer := NewMockdevicer(gomock.NewController(t))
			devicer.EXPECT().ReadByUniqID(gomock.Any(), gomock.Any()).
				Return(lDev, lTest.inpErr).Times(lTest.inpTimes)

			val := Validator{
				devDAO:       devicer,
				vInSub:       vInSub,
				vOutQueue:    vOutQueue,
				vOutPubTopic: "topic-" + random.String(10),
			}
			go func() {
				val.validateMessages()
			}()

			bVIn := []byte("val-aaa")
			if lTest.inpVIn != nil {
				var err error
				bVIn, err = proto.Marshal(lTest.inpVIn)
				require.NoError(t, err)
				t.Logf("bVIn: %s", bVIn)
			}

			require.NoError(t, vInQueue.Publish("", bVIn))

			select {
			case msg := <-vOutSub.C():
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(100 * time.Millisecond):
				// Successful timeout without publish (normally 0.02s).
			}
		})
	}
}
