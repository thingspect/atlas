// +build !unit

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

func TestAccumulateMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inp *message.ValidatorOut
	}{
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: "acc-" + random.String(16), Attr: "acc-motion",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString()},
			OrgId: uuid.NewString()}},
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: "acc-" + random.String(16), Attr: "acc-temp",
			ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString()},
			OrgId: uuid.NewString()}},
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: "acc-" + random.String(16), Attr: "acc-power",
			ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString()},
			OrgId: uuid.NewString()}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can accumulate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bVOut, err := proto.Marshal(lTest.inp)
			require.NoError(t, err)
			t.Logf("bVOut: %s", bVOut)

			require.NoError(t, globalVOutQueue.Publish(globalVOutSubTopic,
				bVOut))
			time.Sleep(2 * time.Second)

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			listPoints, err := globalDPDAO.List(ctx, lTest.inp.OrgId,
				lTest.inp.Point.UniqId, "", "", lTest.inp.Point.Ts.AsTime(),
				lTest.inp.Point.Ts.AsTime().Add(-time.Millisecond))
			t.Logf("listPoints, err: %+v, %v", listPoints, err)
			require.NoError(t, err)
			require.Len(t, listPoints, 1)

			// Normalize token.
			listPoints[0].Token = lTest.inp.Point.Token
			// Normalize timestamp.
			lTest.inp.Point.Ts = timestamppb.New(
				lTest.inp.Point.Ts.AsTime().Truncate(time.Millisecond))

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(lTest.inp.Point, listPoints[0]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.inp.Point,
					listPoints[0])
			}
		})
	}
}

func TestAccumulateMessagesDuplicate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	duplicateVOut := &message.ValidatorOut{Point: &common.DataPoint{
		UniqId: "acc-" + random.String(16), Attr: "acc-motion",
		ValOneof: &common.DataPoint_IntVal{IntVal: 123},
		Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
		Token:    uuid.NewString(), TraceId: uuid.NewString()},
		OrgId: uuid.NewString()}
	require.NoError(t, globalDPDAO.Create(ctx, duplicateVOut.Point,
		duplicateVOut.OrgId))

	bVOut, err := proto.Marshal(duplicateVOut)
	require.NoError(t, err)
	t.Logf("bVOut: %s", bVOut)

	require.NoError(t, globalVOutQueue.Publish(globalVOutSubTopic, bVOut))
	time.Sleep(2 * time.Second)

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	listPoints, err := globalDPDAO.List(ctx, duplicateVOut.OrgId,
		duplicateVOut.Point.UniqId, "", "", duplicateVOut.Point.Ts.AsTime(),
		duplicateVOut.Point.Ts.AsTime().Add(-time.Millisecond))
	t.Logf("listPoints, err: %+v, %v", listPoints, err)
	require.NoError(t, err)
	require.Len(t, listPoints, 1)

	// Normalize token.
	listPoints[0].Token = duplicateVOut.Point.Token
	// Normalize timestamp.
	duplicateVOut.Point.Ts = timestamppb.New(
		duplicateVOut.Point.Ts.AsTime().Truncate(time.Millisecond))

	// Testify does not currently support protobuf equality:
	// https://github.com/stretchr/testify/issues/758
	if !proto.Equal(duplicateVOut.Point, listPoints[0]) {
		t.Fatalf("\nExpect: %+v\nActual: %+v", duplicateVOut.Point,
			listPoints[0])
	}
}
