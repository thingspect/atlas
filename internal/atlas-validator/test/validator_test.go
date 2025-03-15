//go:build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testTimeout = 6 * time.Second

func TestValidateMessages(t *testing.T) {
	now := timestamppb.New(time.Now().Add(-15 * time.Minute))
	traceID := uuid.NewString()
	boolVal := &common.DataPoint_BoolVal{
		BoolVal: []bool{true, false}[random.Intn(2)],
	}

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("val"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := random.Device("val", createOrg.GetId())
	dev.Status = api.Status_ACTIVE
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.ValidatorIn
		res *message.ValidatorOut
	}{
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-motion",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
					Token: createDev.GetToken(), TraceId: traceID,
				}, OrgId: createOrg.GetId(),
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-motion",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
					Token: createDev.GetToken(), TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-temp",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}, Ts: now,
					Token: createDev.GetToken(), TraceId: traceID,
				}, OrgId: createOrg.GetId(),
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-temp",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}, Ts: now,
					Token: createDev.GetToken(), TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"}, Ts: now,
					Token: createDev.GetToken(), TraceId: traceID,
				}, OrgId: createOrg.GetId(),
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"}, Ts: now,
					Token: createDev.GetToken(), TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-leak",
					ValOneof: boolVal, Ts: now, TraceId: traceID,
				}, OrgId: createOrg.GetId(), SkipToken: true,
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-leak",
					ValOneof: boolVal, Ts: now, TraceId: traceID,
				}, Device: dev,
			},
		},
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-leak",
					ValOneof: boolVal, Ts: now, TraceId: traceID,
				}, SkipToken: true,
			}, &message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: "val-leak",
					ValOneof: boolVal, Ts: now, TraceId: traceID,
				}, Device: dev,
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can validate %+v", test), func(t *testing.T) {
			bVIn, err := proto.Marshal(test.inp)
			require.NoError(t, err)
			t.Logf("bVIn: %s", bVIn)

			require.NoError(t, globalValQueue.Publish(globalVInSubTopic, bVIn))

			select {
			case msg := <-globalVOutSub.C():
				msg.Ack()
				t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(),
					msg.Payload())
				require.Equal(t, globalVOutPubTopic, msg.Topic())

				vOut := &message.ValidatorOut{}
				require.NoError(t, proto.Unmarshal(msg.Payload(), vOut))
				t.Logf("vOut: %+v", vOut)
				require.EqualExportedValues(t, test.res, vOut)
			case <-time.After(testTimeout):
				t.Fatal("Message timed out")
			}
		})
	}
}

func TestValidateMessagesError(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("val"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	dev := random.Device("val", createOrg.GetId())
	dev.Status = api.Status_ACTIVE
	createDev, err := globalDevDAO.Create(ctx, dev)
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	disDev := random.Device("val", createOrg.GetId())
	disDev.Status = api.Status_DISABLED
	createDisDev, err := globalDevDAO.Create(ctx, disDev)
	t.Logf("createDisDev, err: %+v, %v", createDisDev, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.ValidatorIn
	}{
		// Bad payload.
		{nil},
		// Device not found.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{UniqId: random.String(16)},
			},
		},
		// Missing value.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: random.String(10),
				},
			},
		},
		// Invalid org ID.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: random.String(10),
					ValOneof: &common.DataPoint_IntVal{},
				}, OrgId: "val-aaa",
			},
		},
		// Device disabled.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDisDev.GetUniqId(), Attr: random.String(10),
					ValOneof: &common.DataPoint_IntVal{},
				}, OrgId: createOrg.GetId(), SkipToken: true,
			},
		},
		// Invalid token.
		{
			&message.ValidatorIn{
				Point: &common.DataPoint{
					UniqId: createDev.GetUniqId(), Attr: random.String(10),
					ValOneof: &common.DataPoint_IntVal{}, Token: "val-aaa",
				}, OrgId: createOrg.GetId(),
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Cannot validate %+v", test), func(t *testing.T) {
			t.Parallel()

			bVIn := []byte("val-aaa")
			if test.inp != nil {
				var err error
				bVIn, err = proto.Marshal(test.inp)
				require.NoError(t, err)
				t.Logf("bVIn: %s", bVIn)
			}

			require.NoError(t, globalValQueue.Publish(globalVInSubTopic, bVIn))

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
