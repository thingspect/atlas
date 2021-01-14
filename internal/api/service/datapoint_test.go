// +build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/api/session"
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

		dpSvc := NewDataPoint(pubQueue, pubTopic)
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

		dpSvc := NewDataPoint(pubQueue, pubTopic)
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

		dpSvc := NewDataPoint(pubQueue, pubTopic)
		_, err := dpSvc.Publish(ctx, &api.PublishDataPointRequest{
			Points: []*common.DataPoint{point}})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.PermissionDenied,
			"permission denied"), err)
	})
}
