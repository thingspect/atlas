package service

//go:generate mockgen -source datapoint.go -destination mock_datapointer_test.go -package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/metric"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DataPointer defines the methods provided by a datapoint.DAO.
type DataPointer interface {
	List(ctx context.Context, orgID, uniqID, devID, attr string, end,
		start time.Time) ([]*common.DataPoint, error)
	Latest(ctx context.Context, orgID, uniqID, devID string,
		start time.Time) ([]*common.DataPoint, error)
}

// DataPoint service contains functions to create and query data points.
type DataPoint struct {
	api.UnimplementedDataPointServiceServer

	dpQueue     queue.Queuer
	vInPubTopic string

	dpDAO DataPointer
}

// NewDataPoint instantiates and returns a new DataPoint service.
func NewDataPoint(
	pubQueue queue.Queuer, pubTopic string, dpDAO DataPointer,
) *DataPoint {
	return &DataPoint{
		dpQueue:     pubQueue,
		vInPubTopic: pubTopic,

		dpDAO: dpDAO,
	}
}

// PublishDataPoints publishes a data point.
func (d *DataPoint) PublishDataPoints(
	ctx context.Context, req *api.PublishDataPointsRequest,
) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_PUBLISHER {
		return nil, errPerm(api.Role_PUBLISHER)
	}

	logger.Logger = logger.WithField("paylType", "api")

	// Build and publish ValidatorIn messages.
	for _, point := range req.GetPoints() {
		vIn := &message.ValidatorIn{
			Point:     point,
			OrgId:     sess.OrgID,
			SkipToken: true,
		}
		vIn.Point.TraceId = sess.TraceID.String()

		// Default to current timestamp if not provided.
		if vIn.GetPoint().GetTs() == nil {
			vIn.Point.Ts = timestamppb.Now()
		}

		bVIn, err := proto.Marshal(vIn)
		if err != nil {
			logger.Errorf("PublishDataPoints proto.Marshal: %v", err)

			return nil, status.Error(codes.Internal, "encode failure")
		}

		if err = d.dpQueue.Publish(d.vInPubTopic, bVIn); err != nil {
			logger.Errorf("PublishDataPoints d.pubQueue.Publish: %v", err)

			return nil, status.Error(codes.Internal, "publish failure")
		}

		metric.Incr("published", nil)
		logger.Debugf("PublishDataPoints published: %+v", vIn)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusAccepted))); err != nil {
		logger.Errorf("PublishDataPoints grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListDataPoints retrieves all data points for a device in a [end, start) time
// range, in descending timestamp order.
func (d *DataPoint) ListDataPoints(
	ctx context.Context, req *api.ListDataPointsRequest,
) (*api.ListDataPointsResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_VIEWER {
		return nil, errPerm(api.Role_VIEWER)
	}

	var uniqID string
	var devID string

	switch v := req.GetIdOneof().(type) {
	case *api.ListDataPointsRequest_UniqId:
		uniqID = v.UniqId
	case *api.ListDataPointsRequest_DeviceId:
		devID = v.DeviceId
	}

	end := time.Now().UTC()
	if req.GetEndTime() != nil {
		end = req.GetEndTime().AsTime()
	}

	start := end.Add(-24 * time.Hour)
	if req.GetStartTime() != nil && req.GetStartTime().AsTime().Before(end) {
		start = req.GetStartTime().AsTime()
	}

	if end.Sub(start) > 90*24*time.Hour {
		return nil, status.Error(codes.InvalidArgument,
			"maximum time range exceeded")
	}

	points, err := d.dpDAO.List(ctx, sess.OrgID, uniqID, devID, req.GetAttr(), end,
		start)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.ListDataPointsResponse{Points: points}, nil
}

// LatestDataPoints retrieves the latest data point for each of a device's
// attributes.
func (d *DataPoint) LatestDataPoints(
	ctx context.Context, req *api.LatestDataPointsRequest,
) (*api.LatestDataPointsResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_VIEWER {
		return nil, errPerm(api.Role_VIEWER)
	}

	var uniqID string
	var devID string

	switch v := req.GetIdOneof().(type) {
	case *api.LatestDataPointsRequest_UniqId:
		uniqID = v.UniqId
	case *api.LatestDataPointsRequest_DeviceId:
		devID = v.DeviceId
	}

	now := time.Now().UTC()

	start := now.Add(30 * -24 * time.Hour)
	if req.GetStartTime() != nil && req.GetStartTime().AsTime().Before(now) {
		start = req.GetStartTime().AsTime()
	}

	if now.Sub(start) > 90*24*time.Hour {
		return nil, status.Error(codes.InvalidArgument,
			"maximum time range exceeded")
	}

	points, err := d.dpDAO.Latest(ctx, sess.OrgID, uniqID, devID, start)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.LatestDataPointsResponse{Points: points}, nil
}
