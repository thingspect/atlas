// +build !unit

package test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPublishDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("Publish valid data point", func(t *testing.T) {
		point := &common.DataPoint{UniqId: "api-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.New(time.Now().Add(-15 * time.Minute))}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAuthGRPCConn)
		_, err := dpCli.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point}})
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
			point.TraceId = vIn.Point.TraceId

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(&message.ValidatorIn{Point: point,
				OrgId: globalAuthOrgID, SkipToken: true}, vIn) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", &message.ValidatorIn{
					Point: point, OrgId: globalAuthOrgID, SkipToken: true}, vIn)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish valid data point without timestamp", func(t *testing.T) {
		point := &common.DataPoint{UniqId: "api-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123}}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAuthGRPCConn)
		_, err := dpCli.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point}})
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
			point.TraceId = vIn.Point.TraceId
			// Normalize timestamps.
			require.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
				5*time.Second)
			point.Ts = vIn.Point.Ts

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(&message.ValidatorIn{Point: point,
				OrgId: globalAuthOrgID, SkipToken: true}, vIn) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", &message.ValidatorIn{
					Point: point, OrgId: globalAuthOrgID, SkipToken: true}, vIn)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish invalid data point", func(t *testing.T) {
		point := &common.DataPoint{UniqId: "api-point-" + random.String(40),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.New(time.Now().Add(-15 * time.Minute))}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAuthGRPCConn)
		_, err := dpCli.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point}})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid PublishDataPointsRequest.Points[0]: embedded message "+
			"failed validation | caused by: invalid DataPoint.UniqId: value "+
			"length must be between 5 and 40 runes, inclusive")
	})
}

func TestLatestDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("Latest data points by valid UniqID and dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dev := &api.Device{UniqId: "api-point-" + random.String(16),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		devCli := api.NewDeviceServiceClient(globalAuthGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: dev})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		dpCli := api.NewDataPointServiceClient(globalAuthGRPCConn)

		points := []*common.DataPoint{
			{UniqId: createDev.UniqId, Attr: "motion",
				ValOneof: &common.DataPoint_IntVal{IntVal: 123},
				TraceId:  uuid.New().String()},
			{UniqId: createDev.UniqId, Attr: "temp",
				ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 9.3},
				TraceId:  uuid.New().String()},
			{UniqId: createDev.UniqId, Attr: "power",
				ValOneof: &common.DataPoint_StrVal{StrVal: "batt"},
				TraceId:  uuid.New().String()},
			{UniqId: createDev.UniqId, Attr: "leak",
				ValOneof: &common.DataPoint_BoolVal{BoolVal: []bool{true,
					false}[random.Intn(2)]}, TraceId: uuid.New().String()},
			{UniqId: createDev.UniqId, Attr: "raw",
				ValOneof: &common.DataPoint_BytesVal{BytesVal: []byte{0x00}},
				TraceId:  uuid.New().String()},
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

				err := globalDPDAO.Create(ctx, point, globalAuthOrgID)
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
		latPointsUniqID, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_UniqId{
					UniqId: createDev.UniqId}})
		t.Logf("latPointsUniqID, err: %+v, %v", latPointsUniqID, err)
		require.NoError(t, err)
		require.Len(t, latPointsUniqID.Points, len(points))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.LatestDataPointsResponse{Points: points},
			latPointsUniqID) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.LatestDataPointsResponse{Points: points}, latPointsUniqID)
		}

		// Verify results by dev ID.
		latPointsDevID, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_DevId{
					DevId: createDev.Id}})
		t.Logf("latPointsDevID, err: %+v, %v", latPointsDevID, err)
		require.NoError(t, err)
		require.Len(t, latPointsDevID.Points, len(points))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.LatestDataPointsResponse{Points: points},
			latPointsDevID) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.LatestDataPointsResponse{Points: points}, latPointsDevID)
		}
	})

	t.Run("Latest data points are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "api-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), TraceId: uuid.New().String()}

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		err := globalDPDAO.Create(ctx, point, uuid.New().String())
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		dpCli := api.NewDataPointServiceClient(globalAuthGRPCConn)
		latPoints, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_UniqId{
					UniqId: uuid.New().String()}})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.NoError(t, err)
		require.Len(t, latPoints.Points, 0)
	})

	t.Run("Latest data points by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAuthGRPCConn)
		latPoints, err := dpCli.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_DevId{
					DevId: random.String(10)}})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: UUID")
	})
}
