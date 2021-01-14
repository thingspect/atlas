// +build !unit

package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPublishDataPoint(t *testing.T) {
	t.Parallel()

	t.Run("Publish valid data point", func(t *testing.T) {
		point := &common.DataPoint{UniqId: "api-point-" + random.String(16),
			Attr: "motion", ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts: timestamppb.New(time.Now().Add(-15 * time.Minute))}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		dpCli := api.NewDataPointServiceClient(globalAuthGRPCConn)
		_, err := dpCli.Publish(ctx, &api.PublishDataPointRequest{
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
		_, err := dpCli.Publish(ctx, &api.PublishDataPointRequest{
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
		_, err := dpCli.Publish(ctx, &api.PublishDataPointRequest{
			Points: []*common.DataPoint{point}})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid PublishDataPointRequest.Points[0]: embedded message "+
			"failed validation | caused by: invalid DataPoint.UniqId: value "+
			"length must be between 5 and 40 runes, inclusive")
	})
}
