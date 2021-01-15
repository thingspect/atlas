// +build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPublishDataPoint(t *testing.T) {
	t.Parallel()

	t.Run("Publish valid data point", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.New().String()
		point := &common.DataPoint{UniqId: random.String(16), Attr: "motion",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute))}

		pubQueue := queue.NewFake()
		pubSub, err := pubQueue.Subscribe("")
		require.NoError(t, err)
		pubTopic := "topic-" + random.String(10)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID}),
			2*time.Second)
		defer cancel()

		dpSvc := NewDataPoint(pubQueue, pubTopic, nil)
		_, err = dpSvc.Publish(ctx, &api.PublishDataPointRequest{
			Points: []*common.DataPoint{point}})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		select {
		case msg := <-pubSub.C():
			msg.Ack()
			t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(), msg.Payload())
			require.Equal(t, pubTopic, msg.Topic())

			vIn := &message.ValidatorIn{}
			require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
			t.Logf("vIn: %+v", vIn)

			// Normalize generated trace ID.
			point.TraceId = vIn.Point.TraceId

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(&message.ValidatorIn{Point: point, OrgId: orgID,
				SkipToken: true}, vIn) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", &message.ValidatorIn{
					Point: point, OrgId: orgID, SkipToken: true}, vIn)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish valid data point without timestamp", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.New().String()
		point := &common.DataPoint{UniqId: random.String(16), Attr: "motion",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123}}

		pubQueue := queue.NewFake()
		pubSub, err := pubQueue.Subscribe("")
		require.NoError(t, err)
		pubTopic := "topic-" + random.String(10)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID}),
			2*time.Second)
		defer cancel()

		dpSvc := NewDataPoint(pubQueue, pubTopic, nil)
		_, err = dpSvc.Publish(ctx, &api.PublishDataPointRequest{
			Points: []*common.DataPoint{point}})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		select {
		case msg := <-pubSub.C():
			msg.Ack()
			t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(), msg.Payload())
			require.Equal(t, pubTopic, msg.Topic())

			vIn := &message.ValidatorIn{}
			require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
			t.Logf("vIn: %+v", vIn)

			// Normalize generated trace ID.
			point.TraceId = vIn.Point.TraceId
			// Normalize timestamps.
			require.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
				2*time.Second)
			point.Ts = vIn.Point.Ts

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(&message.ValidatorIn{Point: point, OrgId: orgID,
				SkipToken: true}, vIn) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", &message.ValidatorIn{
					Point: point, OrgId: orgID, SkipToken: true}, vIn)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish data point with invalid session", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: random.String(16), Attr: "motion",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123}}

		pubQueue := queue.NewFake()
		pubTopic := "topic-" + random.String(10)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dpSvc := NewDataPoint(pubQueue, pubTopic, nil)
		_, err := dpSvc.Publish(ctx, &api.PublishDataPointRequest{
			Points: []*common.DataPoint{point}})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"), err)
	})
}

func TestLatestDataPoint(t *testing.T) {
	t.Parallel()

	t.Run("Latest data points by valid UniqID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "dao-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), TraceId: uuid.New().String()}
		orgID := uuid.New().String()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		datapointer := NewMockDataPointer(ctrl)
		datapointer.EXPECT().Latest(gomock.Any(), orgID, point.UniqId, "").
			Return([]*common.DataPoint{point}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID}),
			2*time.Second)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		latPoints, err := dpSvc.Latest(ctx, &api.LatestDataPointRequest{
			IdOneof: &api.LatestDataPointRequest_UniqId{UniqId: point.UniqId}})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.LatestDataPointResponse{
			Points: []*common.DataPoint{point}}, latPoints) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.LatestDataPointResponse{
					Points: []*common.DataPoint{point}}, latPoints)
		}
	})

	t.Run("Latest data points by valid dev ID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "dao-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.Now(), TraceId: uuid.New().String()}
		orgID := uuid.New().String()
		devID := uuid.New().String()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		datapointer := NewMockDataPointer(ctrl)
		datapointer.EXPECT().Latest(gomock.Any(), orgID, "", devID).
			Return([]*common.DataPoint{point}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID}),
			2*time.Second)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		latPoints, err := dpSvc.Latest(ctx, &api.LatestDataPointRequest{
			IdOneof: &api.LatestDataPointRequest_DevId{DevId: devID}})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.LatestDataPointResponse{
			Points: []*common.DataPoint{point}}, latPoints) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.LatestDataPointResponse{
					Points: []*common.DataPoint{point}}, latPoints)
		}
	})

	t.Run("Latest data points with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		datapointer := NewMockDataPointer(ctrl)
		datapointer.EXPECT().Latest(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any()).Times(0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		latPoints, err := dpSvc.Latest(ctx, &api.LatestDataPointRequest{
			IdOneof: &api.LatestDataPointRequest_UniqId{
				UniqId: random.String(16)}})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"), err)
	})

	t.Run("Latest data points by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		datapointer := NewMockDataPointer(ctrl)
		datapointer.EXPECT().Latest(gomock.Any(), "aaa", gomock.Any(),
			gomock.Any()).Return(nil, dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: "aaa"}),
			2*time.Second)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		latPoints, err := dpSvc.Latest(ctx, &api.LatestDataPointRequest{
			IdOneof: &api.LatestDataPointRequest_UniqId{
				UniqId: random.String(16)}})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}
