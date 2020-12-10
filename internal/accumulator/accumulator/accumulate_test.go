// +build !integration

package accumulator

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/dao/datapoint"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errTestProc = errors.New("accumulator: test processor error")

type fakeDataPointer struct {
	fCreate func() error
}

func newFakeDataPointer() *fakeDataPointer {
	return &fakeDataPointer{
		fCreate: func() error { return nil },
	}
}

func (fp *fakeDataPointer) Create(ctx context.Context, dp *common.DataPoint,
	orgID string) error {
	return fp.fCreate()
}

func TestAccumulateMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inp *message.ValidatorOut
	}{
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: "motion",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.New().String(), TraceId: uuid.New().String()},
			OrgId: uuid.New().String()}},
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: "temp",
			ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.New().String(), TraceId: uuid.New().String()},
			OrgId: uuid.New().String()}},
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: "power",
			ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.New().String(), TraceId: uuid.New().String()},
			OrgId: uuid.New().String()}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can accumulate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			acc := Accumulator{
				dpDAO:     newFakeDataPointer(),
				vOutQueue: vOutQueue,
				vOutSub:   vOutSub,
			}
			go func() {
				acc.accumulateMessages()
			}()

			bVOut, err := proto.Marshal(lTest.inp)
			require.NoError(t, err)
			t.Logf("bVOut: %s", bVOut)

			require.NoError(t, vOutQueue.Publish("", bVOut))
		})
	}
}

func TestAccumulateMessagesError(t *testing.T) {
	t.Parallel()

	duplicateDataPointer := newFakeDataPointer()
	duplicateDataPointer.fCreate = func() error {
		return datapoint.ErrDuplicate
	}

	errDataPointer := newFakeDataPointer()
	errDataPointer.fCreate = func() error {
		return errTestProc
	}

	tests := []struct {
		inpDataPointer *fakeDataPointer
		inpVOut        *message.ValidatorOut
	}{
		// Bad payload.
		{newFakeDataPointer(), nil},
		// Duplicate data point.
		{duplicateDataPointer,
			&message.ValidatorOut{Point: &common.DataPoint{}}},
		// DataPointer error.
		{errDataPointer, &message.ValidatorOut{Point: &common.DataPoint{}}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot accumulate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			acc := Accumulator{
				dpDAO:     lTest.inpDataPointer,
				vOutQueue: vOutQueue,
				vOutSub:   vOutSub,
			}
			go func() {
				acc.accumulateMessages()
			}()

			bVOut := []byte("val-aaa")
			if lTest.inpVOut != nil {
				var err error
				bVOut, err = proto.Marshal(lTest.inpVOut)
				require.NoError(t, err)
				t.Logf("bVOut: %s", bVOut)
			}

			require.NoError(t, vOutQueue.Publish("", bVOut))
		})
	}
}
