// +build !integration

package validator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errTestProc = errors.New("validator: test processor error")

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
	traceID := uuid.New().String()
	orgID := uuid.New().String()
	devID := uuid.New().String()

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
				Attr:     "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"}, Ts: now,
				Token: token, TraceId: traceID}, OrgId: orgID, DevId: devID}},
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

	devID := uuid.New().String()
	orgID := uuid.New().String()
	token := uuid.New().String()

	noRowsDevicer := newFakeDevicer(devID, orgID, false, token)
	noRowsDevicer.fReadByUniqID = func() (*device.Device, error) {
		return nil, sql.ErrNoRows
	}

	errDevicer := newFakeDevicer(devID, orgID, false, token)
	errDevicer.fReadByUniqID = func() (*device.Device, error) {
		return nil, errTestProc
	}

	tests := []struct {
		inpDevicer *fakeDevicer
		inpVIn     *message.ValidatorIn
	}{
		// Bad payload.
		{newFakeDevicer(devID, orgID, false, token), nil},
		// Device not found.
		{noRowsDevicer, &message.ValidatorIn{Point: &common.DataPoint{}}},
		// Devicer error.
		{errDevicer, &message.ValidatorIn{Point: &common.DataPoint{}}},
		// Missing value.
		{newFakeDevicer(devID, orgID, false, token),
			&message.ValidatorIn{Point: &common.DataPoint{}}},
		// Invalid org ID.
		{newFakeDevicer(devID, orgID, false, token),
			&message.ValidatorIn{Point: &common.DataPoint{
				ValOneof: &common.DataPoint_IntVal{}}, OrgId: "val-aaa"}},
		// Device disabled.
		{newFakeDevicer(devID, orgID, true, token),
			&message.ValidatorIn{Point: &common.DataPoint{
				ValOneof: &common.DataPoint_IntVal{}}, OrgId: orgID}},
		// Invalid token.
		{newFakeDevicer(devID, orgID, false, token),
			&message.ValidatorIn{Point: &common.DataPoint{
				ValOneof: &common.DataPoint_IntVal{}, Token: "val-aaa"},
				OrgId: orgID}},
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
				t.Fatalf("Received unexpected msg.Topic, msg.Payload: %v, %s",
					msg.Topic(), msg.Payload())
			case <-time.After(100 * time.Millisecond):
				// Successful timeout without publish (normally 0.02s).
			}
		})
	}
}
