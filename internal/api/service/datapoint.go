package service

//go:generate mockgen -source datapoint.go -destination mock_datapointer_test.go -package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/queue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DataPointer defines the methods provided by a datapoint.DAO.
type DataPointer interface {
	Latest(ctx context.Context, orgID, uniqID,
		devID string) ([]*common.DataPoint, error)
}

// DataPoint service contains functions to create and query data points.
type DataPoint struct {
	api.UnimplementedDataPointServiceServer

	pubQueue queue.Queuer
	pubTopic string

	dpDAO DataPointer
}

// NewDataPoint instantiates and returns a new DataPoint service.
func NewDataPoint(pubQueue queue.Queuer, pubTopic string,
	dpDAO DataPointer) *DataPoint {
	return &DataPoint{
		pubQueue: pubQueue,
		pubTopic: pubTopic,

		dpDAO: dpDAO,
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

// Latest retrieves the latest data point for each of a device's attributes.
func (d *DataPoint) Latest(ctx context.Context,
	req *api.LatestDataPointRequest) (*api.LatestDataPointResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	var uniqID string
	var devID string

	switch v := req.IdOneof.(type) {
	case *api.LatestDataPointRequest_UniqId:
		uniqID = v.UniqId
	case *api.LatestDataPointRequest_DevId:
		devID = v.DevId
	}

	points, err := d.dpDAO.Latest(ctx, sess.OrgID, uniqID, devID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.LatestDataPointResponse{Points: points}, nil
}
