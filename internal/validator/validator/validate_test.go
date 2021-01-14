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

	uniqID := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	token := uuid.New().String()
	traceID := uuid.New().String()
	orgID := uuid.New().String()
	devID := uuid.New().String()
	boolVal := &common.DataPoint_BoolVal{BoolVal: []bool{true,
		false}[random.Intn(2)]}

	tests := []struct {
		inpVIn *message.ValidatorIn
		res    *message.ValidatorOut
	}{
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: now, Token: token, TraceId: traceID}, OrgId: orgID},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: now, Token: token, TraceId: traceID}, OrgId: orgID,
				DevId: devID}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
			Ts: now, Token: token, TraceId: traceID}, OrgId: orgID},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Ts: now, Token: token, TraceId: traceID}, OrgId: orgID,
				DevId: devID}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: "power", ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
			Ts: now, Token: token, TraceId: traceID}, OrgId: orgID},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "power", ValOneof: &common.DataPoint_StrVal{
					StrVal: "batt"}, Ts: now, Token: token, TraceId: traceID},
				OrgId: orgID, DevId: devID}},
		{&message.ValidatorIn{Point: &common.DataPoint{UniqId: uniqID,
			Attr: "leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
			OrgId: orgID, SkipToken: true},
			&message.ValidatorOut{Point: &common.DataPoint{UniqId: uniqID,
				Attr: "leak", ValOneof: boolVal, Ts: now, TraceId: traceID},
				OrgId: orgID, DevId: devID}},
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

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			devicer := NewMockdevicer(ctrl)
			devicer.EXPECT().
				ReadByUniqID(gomock.Any(), lTest.inpVIn.Point.UniqId).
				Return(&api.Device{Id: devID, OrgId: orgID,
					Status: common.Status_ACTIVE, Token: token}, nil).Times(1)

			val := Validator{
				devDAO:       devicer,
				vInSub:       vInSub,
				vOutQueue:    vOutQueue,
				vOutPubTopic: vOutPubTopic,
			}
			go func() {
				val.validateMessages()
			}()

			bVIn, err := proto.Marshal(lTest.inpVIn)
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

	orgID := uuid.New().String()

	tests := []struct {
		inpVIn    *message.ValidatorIn
		inpStatus common.Status
		inpErr    error
		inpTimes  int
	}{
		// Bad payload.
		{nil, common.Status_ACTIVE, nil, 0},
		// Missing data point.
		{&message.ValidatorIn{}, common.Status_ACTIVE, nil, 0},
		// Device not found.
		{&message.ValidatorIn{Point: &common.DataPoint{}}, common.Status_ACTIVE,
			dao.ErrNotFound, 1},
		// Devicer error.
		{&message.ValidatorIn{Point: &common.DataPoint{}}, common.Status_ACTIVE,
			errTestProc, 1},
		// Missing value.
		{&message.ValidatorIn{Point: &common.DataPoint{}}, common.Status_ACTIVE,
			nil, 1},
		// Invalid org ID.
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: random.String(10),
			ValOneof: &common.DataPoint_IntVal{}}, OrgId: "val-aaa"},
			common.Status_ACTIVE, nil, 1},
		// Device status.
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: random.String(10),
			ValOneof: &common.DataPoint_IntVal{}}, OrgId: orgID},
			common.Status_DISABLED, nil, 1},
		// Invalid token.
		{&message.ValidatorIn{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: random.String(10),
			ValOneof: &common.DataPoint_IntVal{}, Token: "val-aaa"},
			OrgId: orgID}, common.Status_ACTIVE, nil, 1},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot validate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			vInQueue := queue.NewFake()
			vInSub, err := vInQueue.Subscribe("")
			require.NoError(t, err)

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			devicer := NewMockdevicer(ctrl)
			devicer.EXPECT().
				ReadByUniqID(gomock.Any(), gomock.Any()).
				Return(&api.Device{Id: uuid.New().String(), OrgId: orgID,
					Status: lTest.inpStatus, Token: uuid.New().String()},
					lTest.inpErr).Times(lTest.inpTimes)

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
