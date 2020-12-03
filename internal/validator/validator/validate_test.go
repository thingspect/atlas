// +build !integration

package validator

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fakeDevicer struct {
	fReadByUniqID func() (*device.Device, error)
}

func newFakeDevicer(id, orgID string, disabled bool,
	token string) *fakeDevicer {
	return &fakeDevicer{
		fReadByUniqID: func() (*device.Device, error) {
			return &device.Device{ID: id, OrgID: orgID, Disabled: disabled,
				Token: token}, nil
		},
	}
}

func (fp *fakeDevicer) ReadByUniqID(ctx context.Context,
	uniqID string) (*device.Device, error) {
	return fp.fReadByUniqID()
}

func TestValidateMessages(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	token := uuid.New().String()
	devID := uuid.New().String()
	orgID := uuid.New().String()
	traceID := uuid.New().String()

	tests := []struct {
		inpVIn *message.ValidatorIn
		res    *message.ValidatorOut
	}{
		{&message.ValidatorIn{UniqId: uniqID, Attr: "motion",
			ValOneof: &message.ValidatorIn_IntVal{IntVal: 123}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "motion",
				ValOneof: &message.ValidatorOut_IntVal{IntVal: 123}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "temp",
			ValOneof: &message.ValidatorIn_Fl64Val{Fl64Val: 20.3}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "temp",
				ValOneof: &message.ValidatorOut_Fl64Val{Fl64Val: 20.3}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "power",
			ValOneof: &message.ValidatorIn_StrVal{StrVal: "batt"}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "power",
				ValOneof: &message.ValidatorOut_StrVal{StrVal: "batt"}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "leak",
			ValOneof: &message.ValidatorIn_BoolVal{BoolVal: true}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "leak",
				ValOneof: &message.ValidatorOut_BoolVal{BoolVal: true}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "metadata",
			MapVal: map[string]string{"val-aaa": "val-bbb"}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "metadata",
				MapVal: map[string]string{"val-aaa": "val-bbb"}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
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

			val := Validator{
				devDAO:       newFakeDevicer(devID, orgID, false, token),
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
				t.Logf("vOut: %#v", vOut)

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.res, vOut) {
					t.Fatalf("Expected, actual: %#v, %#v", lTest.res, vOut)
				}
			case <-time.After(2 * time.Second):
				t.Error("Message timed out")
			}
		})
	}
}

func TestValidateMessagesError(t *testing.T) {
	t.Parallel()

	token := uuid.New().String()
	devID := uuid.New().String()
	orgID := uuid.New().String()

	noRowsDevicer := newFakeDevicer(devID, orgID, false, token)
	noRowsDevicer.fReadByUniqID = func() (*device.Device, error) {
		return nil, sql.ErrNoRows
	}

	tests := []struct {
		inpDevicer *fakeDevicer
		inpVIn     *message.ValidatorIn
	}{
		// Bad payload.
		{newFakeDevicer(devID, orgID, false, token), nil},
		// Device not found.
		{noRowsDevicer, &message.ValidatorIn{}},
		// Invalid org ID.
		{newFakeDevicer(devID, orgID, false, token),
			&message.ValidatorIn{OrgId: "val-aaa"}},
		// Device disabled.
		{newFakeDevicer(devID, orgID, true, token),
			&message.ValidatorIn{OrgId: orgID}},
		// Invalid token.
		{newFakeDevicer(devID, orgID, false, token),
			&message.ValidatorIn{Token: "val-aaa", OrgId: orgID}},
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

			val := Validator{
				devDAO:       lTest.inpDevicer,
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
				t.Errorf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(100 * time.Millisecond):
				// Successful timeout without publish (normally 0.02s).
			}
		})
	}
}

func TestVInToVOut(t *testing.T) {
	t.Parallel()

	uniqID := random.String(16)
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	token := uuid.New().String()
	orgID := uuid.New().String()
	traceID := uuid.New().String()
	devID := uuid.New().String()

	tests := []struct {
		inpVIn *message.ValidatorIn
		res    *message.ValidatorOut
	}{
		{&message.ValidatorIn{UniqId: uniqID, Attr: "motion",
			ValOneof: &message.ValidatorIn_IntVal{IntVal: 123}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "motion",
				ValOneof: &message.ValidatorOut_IntVal{IntVal: 123}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "temp",
			ValOneof: &message.ValidatorIn_Fl64Val{Fl64Val: 20.3}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "temp",
				ValOneof: &message.ValidatorOut_Fl64Val{Fl64Val: 20.3}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "power",
			ValOneof: &message.ValidatorIn_StrVal{StrVal: "batt"}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "power",
				ValOneof: &message.ValidatorOut_StrVal{StrVal: "batt"}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "leak",
			ValOneof: &message.ValidatorIn_BoolVal{BoolVal: true}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "leak",
				ValOneof: &message.ValidatorOut_BoolVal{BoolVal: true}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
		{&message.ValidatorIn{UniqId: uniqID, Attr: "metadata",
			MapVal: map[string]string{"val-aaa": "val-bbb"}, Ts: now,
			Token: token, OrgId: orgID, TraceId: traceID},
			&message.ValidatorOut{UniqId: uniqID, Attr: "metadata",
				MapVal: map[string]string{"val-aaa": "val-bbb"}, Ts: now,
				DevId: devID, OrgId: orgID, TraceId: traceID}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can convert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := vInToVOut(lTest.inpVIn, devID)
			t.Logf("res: %#v", res)

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.res, res) {
				t.Fatalf("Expected, actual: %#v, %#v", lTest.res, res)
			}
		})
	}
}
