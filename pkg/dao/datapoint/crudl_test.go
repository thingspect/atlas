//go:build !unit

package datapoint

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testTimeout = 6 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid data points", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		tests := []struct {
			inp *common.DataPoint
		}{
			{
				&common.DataPoint{
					UniqId: "dao-point-" + random.String(16), Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123},
					Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
				},
			},
			{
				&common.DataPoint{
					UniqId: "dao-point-" + random.String(16), Attr: "temp_c",
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
					Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
				},
			},
			{
				&common.DataPoint{
					UniqId: "dao-point-" + random.String(16), Attr: "power",
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"},
					Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
				},
			},
			{
				&common.DataPoint{
					UniqId: "dao-point-" + random.String(16), Attr: "leak",
					ValOneof: &common.DataPoint_BoolVal{BoolVal: []bool{
						true, false,
					}[random.Intn(2)]}, Ts: timestamppb.Now(),
					TraceId: uuid.NewString(),
				},
			},
			{
				&common.DataPoint{
					UniqId: "dao-point-" + random.String(16), Attr: "raw",
					ValOneof: &common.DataPoint_BytesVal{
						BytesVal: random.Bytes(10),
					}, Ts: timestamppb.Now(), TraceId: uuid.NewString(),
				},
			},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Can create %+v", test), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(t.Context(),
					testTimeout)
				defer cancel()

				err := globalDPDAO.Create(ctx, test.inp, createOrg.GetId())
				t.Logf("err: %v", err)
				require.NoError(t, err)
			})
		}
	})

	t.Run("Create invalid data point", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		tests := []struct {
			inp *common.DataPoint
			err error
		}{
			{
				&common.DataPoint{
					UniqId: "dao-point-" + random.String(40), Attr: "count",
					ValOneof: &common.DataPoint_IntVal{IntVal: 123},
					Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
				}, dao.ErrInvalidFormat,
			},
			{
				&common.DataPoint{
					UniqId: "dao-point-" + random.String(16), Attr: "raw",
					ValOneof: &common.DataPoint_BytesVal{
						BytesVal: random.Bytes(3000),
					}, Ts: timestamppb.Now(), TraceId: uuid.NewString(),
				}, dao.ErrInvalidFormat,
			},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Cannot create %+v", test), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(t.Context(),
					testTimeout)
				defer cancel()

				err := globalDPDAO.Create(ctx, test.inp, orgID)
				t.Logf("err: %#v", err)
				require.ErrorIs(t, err, test.err)
			})
		}
	})

	t.Run("Dedupe data point", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		point := &common.DataPoint{
			UniqId: "dao-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}

		err = globalDPDAO.Create(ctx, point, createOrg.GetId())
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		err = globalDPDAO.Create(ctx, point, createOrg.GetId())
		t.Logf("err: %#v", err)
		require.Equal(t, dao.ErrAlreadyExists, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("List data points by UniqID, dev ID, and attr", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-point",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		points := []*common.DataPoint{
			{
				UniqId: createDev.GetUniqId(), Attr: "count",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				TraceId:  uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "temp_c",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				TraceId:  uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "line"},
				TraceId:  uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "leak",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: []bool{
					true, false,
				}[random.Intn(2)]}, TraceId: uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "raw",
				ValOneof: &common.DataPoint_BytesVal{
					BytesVal: random.Bytes(10),
				}, TraceId: uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "count",
				ValOneof: &common.DataPoint_IntVal{IntVal: 321},
				TraceId:  uuid.NewString(),
			},
		}

		for _, point := range points {
			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()

			// Set a new in-place timestamp.
			point.Ts = timestamppb.New(time.Now().UTC().Truncate(
				time.Millisecond))

			err := globalDPDAO.Create(ctx, point, createOrg.GetId())
			t.Logf("err: %v", err)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		// Flip points to descending timestamp order.
		sort.Slice(points, func(i, j int) bool {
			return points[i].GetTs().AsTime().After(points[j].GetTs().AsTime())
		})

		ctx, cancel = context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		// Verify results by UniqID.
		listPointsUniqID, err := globalDPDAO.List(ctx, createOrg.GetId(),
			createDev.GetUniqId(), "", "", points[0].GetTs().AsTime(),
			points[len(points)-1].GetTs().AsTime().Add(-time.Millisecond))
		t.Logf("listPointsUniqID, err: %+v, %v", listPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listPointsUniqID, len(points))

		for i, point := range points {
			require.EqualExportedValues(t, point, listPointsUniqID[i])
		}

		// Verify results by dev ID without oldest point.
		listPointsDevID, err := globalDPDAO.List(ctx, createOrg.GetId(), "",
			createDev.GetId(), "", points[0].GetTs().AsTime(),
			points[len(points)-1].GetTs().AsTime())
		t.Logf("listPointsDevID, err: %+v, %v", listPointsDevID, err)
		require.NoError(t, err)
		require.Len(t, listPointsDevID, len(points)-1)

		for i, point := range points[:len(points)-1] {
			require.EqualExportedValues(t, point, listPointsDevID[i])
		}

		// Verify results by UniqID and attribute.
		listPointsUniqID, err = globalDPDAO.List(ctx, createOrg.GetId(),
			createDev.GetUniqId(), "", "count", points[0].GetTs().AsTime(),
			points[len(points)-1].GetTs().AsTime().Add(-time.Millisecond))
		t.Logf("listPointsUniqID, err: %+v, %v", listPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listPointsUniqID, 2)

		mcount := 0
		for _, point := range points {
			if point.GetAttr() == "count" {
				require.EqualExportedValues(t, point, listPointsUniqID[mcount])
				mcount++
			}
		}
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		point := &common.DataPoint{
			UniqId: "dao-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}

		err = globalDPDAO.Create(ctx, point, createOrg.GetId())
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		listPoints, err := globalDPDAO.List(ctx, uuid.NewString(),
			point.GetUniqId(), "", "", point.GetTs().AsTime(),
			point.GetTs().AsTime().Add(-time.Millisecond))
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.NoError(t, err)
		require.Empty(t, listPoints)
	})

	t.Run("List data points by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-point",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		start := time.Now().UTC().Add(-time.Millisecond)
		dpStart := time.Time{}

		// The first point intentionally sorts first by attribute.
		points := []*common.DataPoint{
			{
				UniqId: createDev.GetUniqId(), Attr: "count",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				TraceId:  uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "temp_c",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				TraceId:  uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "line"},
				TraceId:  uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "leak",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: []bool{
					true, false,
				}[random.Intn(2)]}, TraceId: uuid.NewString(),
			},
			{
				UniqId: createDev.GetUniqId(), Attr: "raw",
				ValOneof: &common.DataPoint_BytesVal{
					BytesVal: random.Bytes(10),
				}, TraceId: uuid.NewString(),
			},
		}

		for i, point := range points {
			for range random.Intn(6) + 3 {
				ctx, cancel := context.WithTimeout(t.Context(),
					testTimeout)
				defer cancel()

				// Set a new in-place timestamp each point group.
				point.Ts = timestamppb.New(time.Now().UTC().Truncate(
					time.Millisecond))

				// Track the first point's latest time.
				if i == 0 {
					dpStart = point.GetTs().AsTime()
				}

				err := globalDPDAO.Create(ctx, point, createOrg.GetId())
				t.Logf("err: %v", err)
				require.NoError(t, err)
				time.Sleep(time.Millisecond)
			}
		}

		sort.Slice(points, func(i, j int) bool {
			return points[i].GetAttr() < points[j].GetAttr()
		})

		ctx, cancel = context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		// Verify results by UniqID.
		latPointsUniqID, err := globalDPDAO.Latest(ctx, createOrg.GetId(),
			createDev.GetUniqId(), "", start)
		t.Logf("latPointsUniqID, err: %+v, %v", latPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, latPointsUniqID, len(points))

		for i, point := range points {
			require.EqualExportedValues(t, point, latPointsUniqID[i])
		}

		// Verify results by dev ID without oldest point's attribute.
		latPointsDevID, err := globalDPDAO.Latest(ctx, createOrg.GetId(), "",
			createDev.GetId(), dpStart)
		t.Logf("latPointsDevID, err: %+v, %v", latPointsDevID, err)
		require.NoError(t, err)
		require.Len(t, latPointsDevID, len(points)-1)

		for i, point := range points[1:] {
			require.EqualExportedValues(t, point, latPointsDevID[i])
		}
	})

	t.Run("Latest data points are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		point := &common.DataPoint{
			UniqId: "dao-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}

		err = globalDPDAO.Create(ctx, point, createOrg.GetId())
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		latPoints, err := globalDPDAO.Latest(ctx, uuid.NewString(),
			point.GetUniqId(), "", time.Now().UTC().Add(-15*time.Minute))
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.NoError(t, err)
		require.Empty(t, latPoints)
	})

	t.Run("Latest data points by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		latPoints, err := globalDPDAO.Latest(ctx, uuid.NewString(), "",
			random.String(10), time.Now().UTC().Add(-15*time.Minute))
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
