// +build !unit

package datapoint

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid data points", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			inpPoint *common.DataPoint
			inpOrgID string
		}{
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "power", ValOneof: &common.DataPoint_StrVal{
					StrVal: "batt"}, Ts: timestamppb.Now(),
				Token: uuid.New().String(), TraceId: uuid.New().String()},
				uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "leak", ValOneof: &common.DataPoint_BoolVal{
					BoolVal: []bool{true, false}[random.Intn(2)]},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "raw", ValOneof: &common.DataPoint_BytesVal{
					BytesVal: []byte{0x00}}, Ts: timestamppb.Now(),
				Token: uuid.New().String(), TraceId: uuid.New().String()},
				uuid.New().String()},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Can create %+v", lTest), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(context.Background(),
					2*time.Second)
				defer cancel()

				err := globalDPDAO.Create(ctx, lTest.inpPoint, lTest.inpOrgID)
				t.Logf("err: %v", err)
				require.NoError(t, err)
			})
		}
	})

	t.Run("Create invalid data point", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.New().String()

		tests := []struct {
			inpPoint *common.DataPoint
			err      error
		}{
			{&common.DataPoint{UniqId: "dao-point-" + random.String(40),
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, dao.ErrInvalidFormat},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "raw", ValOneof: &common.DataPoint_BytesVal{
					BytesVal: []byte(random.String(256))},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, dao.ErrInvalidFormat},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Cannot create %+v", lTest), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(context.Background(),
					2*time.Second)
				defer cancel()

				err := globalDPDAO.Create(ctx, lTest.inpPoint, orgID)
				t.Logf("err: %#v", err)
				require.Equal(t, lTest.err, err)
			})
		}
	})

	t.Run("Dedupe data point", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "dao-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), Token: uuid.New().String(),
			TraceId: uuid.New().String()}
		orgID := uuid.New().String()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		err := globalDPDAO.Create(ctx, point, orgID)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		err = globalDPDAO.Create(ctx, point, orgID)
		t.Logf("err: %#v", err)
		require.Equal(t, dao.ErrAlreadyExists, err)
	})
}

func TestListDevices(t *testing.T) {
	t.Parallel()

	t.Run("List data points by valid org ID", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			inpPoint *common.DataPoint
			inpOrgID string
		}{
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "power", ValOneof: &common.DataPoint_StrVal{
					StrVal: "batt"}, Ts: timestamppb.Now(),
				Token: uuid.New().String(), TraceId: uuid.New().String()},
				uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "leak", ValOneof: &common.DataPoint_BoolVal{
					BoolVal: []bool{true, false}[random.Intn(2)]},
				Ts: timestamppb.Now(), Token: uuid.New().String(),
				TraceId: uuid.New().String()}, uuid.New().String()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "raw", ValOneof: &common.DataPoint_BytesVal{
					BytesVal: []byte{0x00}}, Ts: timestamppb.Now(),
				Token: uuid.New().String(), TraceId: uuid.New().String()},
				uuid.New().String()},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Can create %+v", lTest), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(context.Background(),
					4*time.Second)
				defer cancel()

				err := globalDPDAO.Create(ctx, lTest.inpPoint, lTest.inpOrgID)
				t.Logf("err: %v", err)
				require.NoError(t, err)

				listPoints, err := globalDPDAO.List(ctx, lTest.inpOrgID,
					lTest.inpPoint.UniqId, lTest.inpPoint.Ts.AsTime(),
					lTest.inpPoint.Ts.AsTime())
				t.Logf("listPoints, err: %+v, %v", listPoints, err)
				require.NoError(t, err)
				require.Len(t, listPoints, 1)

				// Normalize token.
				listPoints[0].Token = lTest.inpPoint.Token
				// Normalize timestamp.
				lTest.inpPoint.Ts = timestamppb.New(
					lTest.inpPoint.Ts.AsTime().Truncate(time.Millisecond))

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(lTest.inpPoint, listPoints[0]) {
					t.Fatalf("\nExpect: %+v\nActual: %+v", lTest.inpPoint,
						listPoints[0])
				}
			})
		}
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "dao-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), Token: uuid.New().String(),
			TraceId: uuid.New().String()}
		orgID := uuid.New().String()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		err := globalDPDAO.Create(ctx, point, orgID)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		listPoints, err := globalDPDAO.List(ctx, uuid.New().String(),
			point.UniqId, point.Ts.AsTime(), point.Ts.AsTime())
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.NoError(t, err)
		require.Len(t, listPoints, 0)
	})

	t.Run("List data points by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listPoints, err := globalDPDAO.List(ctx, random.String(10),
			uuid.New().String(), time.Now(), time.Now())
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.Equal(t, dao.ErrInvalidFormat, err)
	})
}
