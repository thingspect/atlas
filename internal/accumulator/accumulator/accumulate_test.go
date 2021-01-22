// +build !integration

package accumulator

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var errTestProc = errors.New("accumulator: test processor error")

func TestAccumulateMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inp *message.ValidatorOut
	}{
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: "motion",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString()},
			OrgId: uuid.NewString()}},
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: "temp",
			ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString()},
			OrgId: uuid.NewString()}},
		{&message.ValidatorOut{Point: &common.DataPoint{
			UniqId: random.String(16), Attr: "power",
			ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
			Token:    uuid.NewString(), TraceId: uuid.NewString()},
			OrgId: uuid.NewString()}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can accumulate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var wg sync.WaitGroup
			wg.Add(1)

			datapointer := NewMockdatapointer(ctrl)
			datapointer.EXPECT().
				Create(gomock.Any(), matcher.NewProtoMatcher(lTest.inp.Point),
					lTest.inp.OrgId).
				DoAndReturn(func(_ ...interface{}) error {
					defer wg.Done()
					return nil
				}).Times(1)

			acc := Accumulator{
				dpDAO:     datapointer,
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
			wg.Wait()
		})
	}
}

func TestAccumulateMessagesError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpVOut  *message.ValidatorOut
		inpErr   error
		inpTimes int
	}{
		// Bad payload.
		{nil, nil, 0},
		// Missing data point.
		{&message.ValidatorOut{}, nil, 0},
		// Duplicate data point.
		{&message.ValidatorOut{Point: &common.DataPoint{}},
			dao.ErrAlreadyExists, 1},
		// DataPointer error.
		{&message.ValidatorOut{Point: &common.DataPoint{}}, errTestProc, 1},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot accumulate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			vOutQueue := queue.NewFake()
			vOutSub, err := vOutQueue.Subscribe("")
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var wg sync.WaitGroup
			wg.Add(1)

			datapointer := NewMockdatapointer(ctrl)
			datapointer.EXPECT().
				Create(gomock.Any(), gomock.Any(), "").
				DoAndReturn(func(_ ...interface{}) error {
					defer wg.Done()
					return lTest.inpErr
				}).Times(lTest.inpTimes)

			acc := Accumulator{
				dpDAO:     datapointer,
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
			if lTest.inpTimes > 0 {
				wg.Wait()
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		})
	}
}
