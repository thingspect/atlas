//go:build !integration

package validator

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const errTestProc consterr.Error = "validator: test processor error"

func TestValidateMessages(t *testing.T) {
	t.Parallel()

	dev := random.Device("val", uuid.NewString())
	dev.Status = api.Status_ACTIVE
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()
	boolVal := &common.DataPoint_BoolVal{
		BoolVal: []bool{true, false}[random.Intn(2)],
	}

	tests := []struct {
		inp *message.ValidatorIn
		res *message.ValidatorOut
	}{
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
					Token: dev.GetToken(), TraceId: traceID,
				}, OrgId: dev.GetOrgId(),
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
					Token: dev.GetToken(), TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "temp_c",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}, Ts: now,
					Token: dev.GetToken(), TraceId: traceID,
				}, OrgId: dev.GetOrgId(),
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "temp_c",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}, Ts: now,
					Token: dev.GetToken(), TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"}, Ts: now,
					Token: dev.GetToken(), TraceId: traceID,
				}, OrgId: dev.GetOrgId(),
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"}, Ts: now,
					Token: dev.GetToken(), TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "leak", ValOneof: boolVal,
					Ts: now, TraceId: traceID,
				}, OrgId: dev.GetOrgId(), SkipToken: true,
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "leak", ValOneof: boolVal,
					Ts: now, TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "leak", ValOneof: boolVal,
					Ts: now, TraceId: traceID,
				}, SkipToken: true,
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: dev.GetUniqId(), Attr: "leak", ValOneof: boolVal,
					Ts: now, TraceId: traceID,
				}, Device: dev,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can validate %+v", test), func(t *testing.T) {
			t.Parallel()

			vInQueue := queue.NewFake()
			vInSub, err := vInQueue.Subscribe("")
			require.NoError(t, err)

			valQueue := queue.NewFake()
			vOutSub, err := valQueue.Subscribe("")
			require.NoError(t, err)
			vOutPubTopic := "topic-" + random.String(10)

			devicer := NewMockdevicer(gomock.NewController(t))
			devicer.EXPECT().ReadByUniqID(gomock.Any(),
				test.inp.GetPoint().GetUniqId()).Return(dev, nil).Times(1)

			val := Validator{
				devDAO: devicer,

				valQueue:     valQueue,
				vInSub:       vInSub,
				vOutPubTopic: vOutPubTopic,
			}
			go func() {
				val.validateMessages()
			}()

			bVIn, err := proto.Marshal(test.inp)
			require.NoError(t, err)
			t.Logf("bVIn: %s", bVIn)

			// Result must be cloned due to its shared pointers across tests.
			res, _ := proto.Clone(test.res).(*message.ValidatorOut)

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
				require.EqualExportedValues(t, res, vOut)
			case <-time.After(2 * time.Second):
				t.Fatal("Message timed out")
			}
		})
	}
}

func TestValidateMessagesError(t *testing.T) {
	t.Parallel()

	orgID := uuid.NewString()

	tests := []struct {
		inpVIn    *message.ValidatorIn
		inpStatus api.Status
		inpErr    error
		inpTimes  int
	}{
		// Bad payload.
		{
			nil, api.Status_ACTIVE, nil, 0,
		},
		// Missing data point.
		{
			&message.ValidatorIn{}, api.Status_ACTIVE, nil, 0,
		},
		// Device not found.
		{
			&message.ValidatorIn{Point: &common.DataPoint{}},
			api.Status_ACTIVE, dao.ErrNotFound, 1,
		},
		// Devicer error.
		{
			&message.ValidatorIn{Point: &common.DataPoint{}},
			api.Status_ACTIVE, errTestProc, 1,
		},
		// Missing value.
		{
			&message.ValidatorIn{Point: &common.DataPoint{}},
			api.Status_ACTIVE, nil, 1,
		},
		// Invalid org ID.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: random.String(16), Attr: random.String(10),
					ValOneof: &common.DataPoint_IntVal{},
				}, OrgId: "val-aaa",
			}, api.Status_ACTIVE, nil, 1,
		},
		// Device status.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: random.String(16), Attr: random.String(10),
					ValOneof: &common.DataPoint_IntVal{},
				}, OrgId: orgID,
			}, api.Status_DISABLED, nil, 1,
		},
		// Invalid token.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: random.String(16), Attr: random.String(10),
					ValOneof: &common.DataPoint_IntVal{}, Token: "val-aaa",
				}, OrgId: orgID,
			}, api.Status_ACTIVE, nil, 1,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Cannot validate %+v", test), func(t *testing.T) {
			t.Parallel()

			dev := random.Device("val", orgID)
			dev.Status = test.inpStatus

			vInQueue := queue.NewFake()
			vInSub, err := vInQueue.Subscribe("")
			require.NoError(t, err)

			valQueue := queue.NewFake()
			vOutSub, err := valQueue.Subscribe("")
			require.NoError(t, err)

			devicer := NewMockdevicer(gomock.NewController(t))
			devicer.EXPECT().ReadByUniqID(gomock.Any(), gomock.Any()).
				Return(dev, test.inpErr).Times(test.inpTimes)

			val := Validator{
				devDAO: devicer,

				valQueue:     valQueue,
				vInSub:       vInSub,
				vOutPubTopic: "topic-" + random.String(10),
			}
			go func() {
				val.validateMessages()
			}()

			bVIn := []byte("val-aaa")
			if test.inpVIn != nil {
				var err error
				bVIn, err = proto.Marshal(test.inpVIn)
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
