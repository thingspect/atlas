// +build !unit

package datapoint

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/org"
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
				Ts: timestamppb.Now(), TraceId: uuid.NewString()},
				uuid.NewString()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "temp", ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				Ts: timestamppb.Now(), TraceId: uuid.NewString()},
				uuid.NewString()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "power", ValOneof: &common.DataPoint_StrVal{
					StrVal: "batt"}, Ts: timestamppb.Now(),
				TraceId: uuid.NewString()}, uuid.NewString()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "leak", ValOneof: &common.DataPoint_BoolVal{
					BoolVal: []bool{true, false}[random.Intn(2)]},
				Ts: timestamppb.Now(), TraceId: uuid.NewString()},
				uuid.NewString()},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "raw", ValOneof: &common.DataPoint_BytesVal{
					BytesVal: []byte{0x00}}, Ts: timestamppb.Now(),
				TraceId: uuid.NewString()}, uuid.NewString()},
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

		orgID := uuid.NewString()

		tests := []struct {
			inpPoint *common.DataPoint
			err      error
		}{
			{&common.DataPoint{UniqId: "dao-point-" + random.String(40),
				Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				Ts: timestamppb.Now(), TraceId: uuid.NewString()},
				dao.ErrInvalidFormat},
			{&common.DataPoint{UniqId: "dao-point-" + random.String(16),
				Attr: "raw", ValOneof: &common.DataPoint_BytesVal{
					BytesVal: []byte(random.String(256))},
				Ts: timestamppb.Now(), TraceId: uuid.NewString()},
				dao.ErrInvalidFormat},
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
				require.ErrorIs(t, err, lTest.err)
			})
		}
	})

	t.Run("Dedupe data point", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "dao-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), TraceId: uuid.NewString()}
		orgID := uuid.NewString()

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

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("List data points by UniqID, dev ID, and attr", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := org.Org{Name: "dao-point-" + random.String(10)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-point-" +
			random.String(16), Status: []api.Status{api.Status_ACTIVE,
			api.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		points := []*common.DataPoint{
			{UniqId: createDev.UniqId, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				TraceId:  uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				TraceId:  uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
				TraceId:  uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "leak",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: []bool{true,
					false}[random.Intn(2)]}, TraceId: uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "raw",
				ValOneof: &common.DataPoint_BytesVal{BytesVal: []byte{0x00}},
				TraceId:  uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 321},
				TraceId:  uuid.NewString()},
		}

		for _, point := range points {
			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			// Set a new in-place timestamp.
			point.Ts = timestamppb.New(time.Now().UTC().Truncate(
				time.Millisecond))
			time.Sleep(time.Millisecond)

			err := globalDPDAO.Create(ctx, point, createOrg.ID)
			t.Logf("err: %v", err)
			require.NoError(t, err)
		}

		sort.Slice(points, func(i, j int) bool {
			return points[i].Ts.AsTime().After(points[j].Ts.AsTime())
		})

		ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		// Verify results by UniqID.
		listPointsUniqID, err := globalDPDAO.List(ctx, createOrg.ID,
			createDev.UniqId, "", "", points[0].Ts.AsTime(),
			points[len(points)-1].Ts.AsTime().Add(-time.Millisecond))
		t.Logf("listPointsUniqID, err: %+v, %v", listPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listPointsUniqID, len(points))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, point := range points {
			if !proto.Equal(point, listPointsUniqID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", point,
					listPointsUniqID[i])
			}
		}

		// Verify results by dev ID without oldest point.
		listPointsDevID, err := globalDPDAO.List(ctx, createOrg.ID, "",
			createDev.Id, "", points[0].Ts.AsTime(),
			points[len(points)-1].Ts.AsTime())
		t.Logf("listPointsDevID, err: %+v, %v", listPointsDevID, err)
		require.NoError(t, err)
		require.Len(t, listPointsDevID, len(points)-1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, point := range points[:len(points)-1] {
			if !proto.Equal(point, listPointsDevID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", point,
					listPointsDevID[i])
			}
		}

		// Verify results by UniqID and attribute.
		listPointsUniqID, err = globalDPDAO.List(ctx, createOrg.ID,
			createDev.UniqId, "", "motion", points[0].Ts.AsTime(),
			points[len(points)-1].Ts.AsTime().Add(-time.Millisecond))
		t.Logf("listPointsUniqID, err: %+v, %v", listPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listPointsUniqID, 2)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		mcount := 0
		for _, point := range points {
			if point.Attr == "motion" {
				if !proto.Equal(point, listPointsUniqID[mcount]) {
					t.Fatalf("\nExpect: %+v\nActual: %+v", point,
						listPointsUniqID[mcount])
				}
				mcount++
			}
		}
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "dao-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), TraceId: uuid.NewString()}
		orgID := uuid.NewString()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		err := globalDPDAO.Create(ctx, point, orgID)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		listPoints, err := globalDPDAO.List(ctx, uuid.NewString(),
			point.UniqId, "", "", point.Ts.AsTime(), point.Ts.AsTime())
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.NoError(t, err)
		require.Len(t, listPoints, 0)
	})

	t.Run("List data points by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listPoints, err := globalDPDAO.List(ctx, random.String(10),
			uuid.NewString(), "", "", time.Now(), time.Now())
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestLatest(t *testing.T) {
	t.Parallel()

	t.Run("Latest data points by valid UniqID and dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := org.Org{Name: "dao-point-" + random.String(10)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		dev := &api.Device{OrgId: createOrg.ID, UniqId: "dao-point-" +
			random.String(16), Status: []api.Status{api.Status_ACTIVE,
			api.Status_DISABLED}[random.Intn(2)]}
		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		points := []*common.DataPoint{
			{UniqId: createDev.UniqId, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				TraceId:  uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				TraceId:  uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
				TraceId:  uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "leak",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: []bool{true,
					false}[random.Intn(2)]}, TraceId: uuid.NewString()},
			{UniqId: createDev.UniqId, Attr: "raw",
				ValOneof: &common.DataPoint_BytesVal{BytesVal: []byte{0x00}},
				TraceId:  uuid.NewString()},
		}

		for _, point := range points {
			for i := 0; i < random.Intn(6)+3; i++ {
				ctx, cancel := context.WithTimeout(context.Background(),
					2*time.Second)
				defer cancel()

				// Set a new in-place timestamp each pass.
				point.Ts = timestamppb.New(time.Now().UTC().Truncate(
					time.Millisecond))
				time.Sleep(time.Millisecond)

				err := globalDPDAO.Create(ctx, point, createOrg.ID)
				t.Logf("err: %v", err)
				require.NoError(t, err)
			}
		}

		sort.Slice(points, func(i, j int) bool {
			return points[i].Attr < points[j].Attr
		})

		ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		// Verify results by UniqID.
		latPointsUniqID, err := globalDPDAO.Latest(ctx, createOrg.ID,
			createDev.UniqId, "")
		t.Logf("latPointsUniqID, err: %+v, %v", latPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, latPointsUniqID, len(points))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, point := range points {
			if !proto.Equal(point, latPointsUniqID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", point,
					latPointsUniqID[i])
			}
		}

		// Verify results by dev ID.
		latPointsDevID, err := globalDPDAO.Latest(ctx, createOrg.ID, "",
			createDev.Id)
		t.Logf("latPointsDevID, err: %+v, %v", latPointsDevID, err)
		require.NoError(t, err)
		require.Len(t, latPointsDevID, len(points))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, point := range points {
			if !proto.Equal(point, latPointsDevID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", point, latPointsDevID[i])
			}
		}
	})

	t.Run("Latest data points are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "dao-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), TraceId: uuid.NewString()}

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		err := globalDPDAO.Create(ctx, point, uuid.NewString())
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		latPoints, err := globalDPDAO.Latest(ctx, uuid.NewString(),
			point.UniqId, "")
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.NoError(t, err)
		require.Len(t, latPoints, 0)
	})

	t.Run("Latest data points by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		latPoints, err := globalDPDAO.Latest(ctx, uuid.NewString(), "",
			random.String(10))
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
