//go:build !unit

package test

import (
	"context"
	"sort"
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

func TestPublishDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("Publish valid data point", func(t *testing.T) {
		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		_, err := dpCli.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point},
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		select {
		case msg := <-globalPubSub.C():
			msg.Ack()
			t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(), msg.Payload())
			require.Equal(t, globalPubTopic, msg.Topic())

			vIn := &message.ValidatorIn{}
			require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
			t.Logf("vIn: %+v", vIn)

			// Normalize generated trace ID.
			point.TraceId = vIn.GetPoint().GetTraceId()

			require.EqualExportedValues(t, &message.ValidatorIn{
				Point: point, OrgId: globalAdminOrgID, SkipToken: true,
			}, vIn)
		case <-time.After(testTimeout):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish valid data point without timestamp", func(t *testing.T) {
		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAdminKeyGRPCConn)
		_, err := dpCli.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point},
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		select {
		case msg := <-globalPubSub.C():
			msg.Ack()
			t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(), msg.Payload())
			require.Equal(t, globalPubTopic, msg.Topic())

			vIn := &message.ValidatorIn{}
			require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
			t.Logf("vIn: %+v", vIn)

			// Normalize generated trace ID.
			point.TraceId = vIn.GetPoint().GetTraceId()
			// Normalize timestamp.
			require.WithinDuration(t, time.Now(),
				vIn.GetPoint().GetTs().AsTime(), 5*time.Second)
			point.Ts = vIn.GetPoint().GetTs()

			require.EqualExportedValues(t, &message.ValidatorIn{
				Point: point, OrgId: globalAdminOrgID, SkipToken: true,
			}, vIn)
		case <-time.After(testTimeout):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish data point with insufficient role", func(t *testing.T) {
		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(secondaryViewerGRPCConn)
		_, err := dpCli.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point},
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, PUBLISHER role required")
	})

	t.Run("Publish invalid data point", func(t *testing.T) {
		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(40), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		_, err := dpCli.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point},
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid PublishDataPointsRequest.Points[0]: embedded message "+
			"failed validation | caused by: invalid DataPoint.UniqId: value "+
			"length must be between 5 and 40 runes, inclusive")
	})
}

func TestListDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("List data points by UniqID, dev ID, and attr", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-point", uuid.NewString()),
		})
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

			err := globalDPDAO.Create(ctx, point, globalAdminOrgID)
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
		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		listPointsUniqID, err := dpCli.ListDataPoints(ctx,
			&api.ListDataPointsRequest{
				IdOneof: &api.ListDataPointsRequest_UniqId{
					UniqId: createDev.GetUniqId(),
				}, EndTime: points[0].GetTs(),
				StartTime: timestamppb.New(points[len(points)-1].GetTs().
					AsTime().Add(-time.Millisecond)),
			})
		t.Logf("listPointsUniqID, err: %+v, %v", listPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listPointsUniqID.GetPoints(), len(points))
		require.EqualExportedValues(t,
			&api.ListDataPointsResponse{Points: points}, listPointsUniqID)

		// Verify results by dev ID without oldest point.
		listPointsDevID, err := dpCli.ListDataPoints(ctx,
			&api.ListDataPointsRequest{
				IdOneof: &api.ListDataPointsRequest_DeviceId{
					DeviceId: createDev.GetId(),
				}, StartTime: points[len(points)-1].GetTs(),
			})
		t.Logf("listPointsDevID, err: %+v, %v", listPointsDevID, err)
		require.NoError(t, err)
		require.Len(t, listPointsDevID.GetPoints(), len(points)-1)
		require.EqualExportedValues(t, &api.ListDataPointsResponse{
			Points: points[:len(points)-1],
		}, listPointsDevID)

		// Verify results by UniqID and attribute.
		listPointsUniqID, err = dpCli.ListDataPoints(ctx,
			&api.ListDataPointsRequest{
				IdOneof: &api.ListDataPointsRequest_UniqId{
					UniqId: createDev.GetUniqId(),
				}, Attr: "count",
				StartTime: timestamppb.New(points[len(points)-1].GetTs().
					AsTime().Add(-time.Millisecond)),
			})
		t.Logf("listPointsUniqID, err: %+v, %v", listPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listPointsUniqID.GetPoints(), 2)

		mcount := 0
		for _, point := range points {
			if point.GetAttr() == "count" {
				require.EqualExportedValues(t, point,
					listPointsUniqID.GetPoints()[mcount])
				mcount++
			}
		}
	})

	t.Run("List data points are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		err = globalDPDAO.Create(ctx, point, createOrg.GetId())
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		listPoints, err := dpCli.ListDataPoints(ctx, &api.ListDataPointsRequest{
			IdOneof: &api.ListDataPointsRequest_UniqId{UniqId: point.GetUniqId()},
		})
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.NoError(t, err)
		require.Empty(t, listPoints.GetPoints())
	})

	t.Run("List data points by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		listPoints, err := dpCli.ListDataPoints(ctx, &api.ListDataPointsRequest{
			IdOneof: &api.ListDataPointsRequest_DeviceId{
				DeviceId: uuid.NewString(),
			}, EndTime: timestamppb.Now(),
			StartTime: timestamppb.New(time.Now().Add(-91 * 24 * time.Hour)),
		})
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"maximum time range exceeded")
	})

	t.Run("List data points by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		listPoints, err := dpCli.ListDataPoints(ctx, &api.ListDataPointsRequest{
			IdOneof: &api.ListDataPointsRequest_DeviceId{
				DeviceId: random.String(10),
			},
		})
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid ListDataPointsRequest.DeviceId: value must be a valid "+
			"UUID | caused by: invalid uuid format")
	})
}

func TestLatestDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("Latest data points by valid UniqID and dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-point", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		dpStart := timestamppb.Now()

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

				// Set a new in-place timestamp each pass.
				point.Ts = timestamppb.New(time.Now().UTC().Truncate(
					time.Millisecond))

				// Track the first point's latest time.
				if i == 0 {
					dpStart = timestamppb.Now()
				}

				err := globalDPDAO.Create(ctx, point, globalAdminOrgID)
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
		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		latPointsUniqID, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_UniqId{
					UniqId: createDev.GetUniqId(),
				},
			})
		t.Logf("latPointsUniqID, err: %+v, %v", latPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, latPointsUniqID.GetPoints(), len(points))
		require.EqualExportedValues(t,
			&api.LatestDataPointsResponse{Points: points}, latPointsUniqID)

		// Verify results by dev ID without oldest point's attribute.
		latPointsDevID, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_DeviceId{
					DeviceId: createDev.GetId(),
				}, StartTime: dpStart,
			})
		t.Logf("latPointsDevID, err: %+v, %v", latPointsDevID, err)
		require.NoError(t, err)
		require.Len(t, latPointsDevID.GetPoints(), len(points)-1)
		require.EqualExportedValues(t,
			&api.LatestDataPointsResponse{Points: points[1:]}, latPointsDevID)
	})

	t.Run("Latest data points are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-point"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		err = globalDPDAO.Create(ctx, point, createOrg.GetId())
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		latPoints, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_UniqId{
					UniqId: point.GetUniqId(),
				},
			})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.NoError(t, err)
		require.Empty(t, latPoints.GetPoints())
	})

	t.Run("Latest data points by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		latPoints, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_DeviceId{
					DeviceId: uuid.NewString(),
				}, StartTime: timestamppb.New(
					time.Now().Add(-91 * 24 * time.Hour)),
			})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"maximum time range exceeded")
	})

	t.Run("Latest data points by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAdminGRPCConn)
		latPoints, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_DeviceId{
					DeviceId: random.String(10),
				},
			})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid LatestDataPointsRequest.DeviceId: value must be a valid "+
			"UUID | caused by: invalid uuid format")
	})
}
