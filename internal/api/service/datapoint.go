package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DataPoint service contains functions to create and query data points.
type DataPoint struct {
	api.UnimplementedDataPointServiceServer

	pubQueue queue.Queuer
	pubTopic string
}

// NewDataPoint instantiates and returns a new DataPoint service.
func NewDataPoint(pubQueue queue.Queuer, pubTopic string) *DataPoint {
	return &DataPoint{
		pubQueue: pubQueue,
		pubTopic: pubTopic,
	}
}

// Publish publishes a data point.
func (d *DataPoint) Publish(ctx context.Context,
	req *api.PublishDataPointRequest) (*empty.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	logger.Logger = logger.WithStr("paylType", "api")

	// Build and publish ValidatorIn messages.
	for _, point := range req.Points {
		// Set up per-point logging fields.
		traceID := uuid.New().String()
		logger := logger.WithStr("traceID", traceID)

		vIn := &message.ValidatorIn{
			Point:     point,
			OrgId:     sess.OrgID,
			SkipToken: true,
		}
		vIn.Point.TraceId = traceID

		// Default to current timestamp if not provided.
		if vIn.Point.Ts == nil {
			vIn.Point.Ts = timestamppb.Now()
		}

		bVIn, err := proto.Marshal(vIn)
		if err != nil {
			logger.Errorf("DataPoint.Publish proto.Marshal: %v", err)
			return nil, status.Error(codes.Internal, "publish failure")
		}

		if err = d.pubQueue.Publish(d.pubTopic, bVIn); err != nil {
			logger.Errorf("DataPoint.Publish d.pubQueue.Publish: %v", err)
			return nil, status.Error(codes.Internal, "publish failure")
		}

		logger.Debugf("DataPoint.Publish published: %+v", vIn)
	}

	return &empty.Empty{}, nil
}
