//go:build !unit

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
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testTimeout = 6 * time.Second

func TestAccumulateMessages(t *testing.T) {
	t.Parallel()

	now := timestamppb.New(time.Now().Add(-15 * time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("acc"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createDev, err := globalDevDAO.Create(ctx, random.Device("acc",
		createOrg.GetId()))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	tests := []struct {
		inp *message.ValidatorOut
	}{
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: "acc-" + random.String(16), Attr: "acc-motion",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123}, Ts: now,
					Token: uuid.NewString(), TraceId: uuid.NewString(),
				}, Device: createDev,
			},
		},
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: "acc-" + random.String(16), Attr: "acc-temp",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3}, Ts: now,
					Token: uuid.NewString(), TraceId: uuid.NewString(),
				}, Device: createDev,
			},
		},
		{
			&message.ValidatorOut{
				Point: &common.DataPoint{
					UniqId: "acc-" + random.String(16), Attr: "acc-power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"}, Ts: now,
					Token: uuid.NewString(), TraceId: uuid.NewString(),
				}, Device: createDev,
			},
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can accumulate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bVOut, err := proto.Marshal(lTest.inp)
			require.NoError(t, err)
			t.Logf("bVOut: %s", bVOut)

			require.NoError(t, globalAccQueue.Publish(globalVOutSubTopic,
				bVOut))
			time.Sleep(2 * time.Second)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			listPoints, err := globalDPDAO.List(ctx, lTest.inp.GetDevice().GetOrgId(),
				lTest.inp.GetPoint().GetUniqId(), "", "", lTest.inp.GetPoint().GetTs().AsTime(),
				lTest.inp.GetPoint().GetTs().AsTime().Add(-time.Millisecond))
			t.Logf("listPoints, err: %+v, %v", listPoints, err)
			require.NoError(t, err)
			require.Len(t, listPoints, 1)

			// Normalize token.
			listPoints[0].Token = lTest.inp.GetPoint().GetToken()
			// Normalize timestamp.
			lTest.inp.Point.Ts = timestamppb.New(
				lTest.inp.GetPoint().GetTs().AsTime().Truncate(time.Millisecond))

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.inp.GetPoint(), listPoints[0]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.inp.GetPoint(),
					listPoints[0])
			}
		})
	}
}

func TestAccumulateMessagesDuplicate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("acc"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createDev, err := globalDevDAO.Create(ctx, random.Device("acc",
		createOrg.GetId()))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	duplicateVOut := &message.ValidatorOut{
		Point: &common.DataPoint{
			UniqId: "acc-" + random.String(16), Attr: "acc-motion",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString(),
		}, Device: createDev,
	}
	require.NoError(t, globalDPDAO.Create(ctx, duplicateVOut.GetPoint(),
		duplicateVOut.GetDevice().GetOrgId()))

	bVOut, err := proto.Marshal(duplicateVOut)
	require.NoError(t, err)
	t.Logf("bVOut: %s", bVOut)

	require.NoError(t, globalAccQueue.Publish(globalVOutSubTopic, bVOut))
	time.Sleep(2 * time.Second)

	ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	listPoints, err := globalDPDAO.List(ctx, duplicateVOut.GetDevice().GetOrgId(),
		duplicateVOut.GetPoint().GetUniqId(), "", "", duplicateVOut.GetPoint().GetTs().AsTime(),
		duplicateVOut.GetPoint().GetTs().AsTime().Add(-time.Millisecond))
	t.Logf("listPoints, err: %+v, %v", listPoints, err)
	require.NoError(t, err)
	require.Len(t, listPoints, 1)

	// Normalize token.
	listPoints[0].Token = duplicateVOut.GetPoint().GetToken()
	// Normalize timestamp.
	duplicateVOut.Point.Ts = timestamppb.New(
		duplicateVOut.GetPoint().GetTs().AsTime().Truncate(time.Millisecond))

	// Testify does not currently support protobuf equality:
	// https://github.com/stretchr/testify/issues/758
	if !proto.Equal(duplicateVOut.GetPoint(), listPoints[0]) {
		t.Fatalf("\nExpect: %+v\nActual: %+v", duplicateVOut.GetPoint(),
			listPoints[0])
	}
}

func TestAccumulateMessagesError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("acc"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createDev, err := globalDevDAO.Create(ctx, random.Device("acc",
		createOrg.GetId()))
	t.Logf("createDev, err: %+v, %v", createDev, err)
	require.NoError(t, err)

	invalidVOut := &message.ValidatorOut{
		Point: &common.DataPoint{
			UniqId: "acc-" + random.String(16), Attr: "acc-raw",
			ValOneof: &common.DataPoint_BytesVal{BytesVal: random.Bytes(3000)},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString(),
		}, Device: createDev,
	}

	bVOut, err := proto.Marshal(invalidVOut)
	require.NoError(t, err)
	t.Logf("bVOut: %s", bVOut)

	require.NoError(t, globalAccQueue.Publish(globalVOutSubTopic, bVOut))
	time.Sleep(2 * time.Second)

	ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	listPoints, err := globalDPDAO.List(ctx, invalidVOut.GetDevice().GetOrgId(),
		invalidVOut.GetPoint().GetUniqId(), "", "", invalidVOut.GetPoint().GetTs().AsTime(),
		invalidVOut.GetPoint().GetTs().AsTime().Add(-time.Millisecond))
	t.Logf("listPoints, err: %+v, %v", listPoints, err)
	require.NoError(t, err)
	require.Len(t, listPoints, 0)
}
