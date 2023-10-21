package service

//go:generate mockgen -source alert.go -destination mock_alerter_test.go -package service

import (
	"context"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Alerter defines the methods provided by an alert.DAO.
type Alerter interface {
	List(ctx context.Context, orgID, uniqID, devID, alarmID, userID string, end,
		start time.Time) ([]*api.Alert, error)
}

// Alert service contains functions to query alerts.
type Alert struct {
	api.UnimplementedAlertServiceServer

	aleDAO Alerter
}

// NewAlert instantiates and returns a new Alert service.
func NewAlert(aleDAO Alerter) *Alert {
	return &Alert{
		aleDAO: aleDAO,
	}
}

// ListAlerts retrieves all alerts for a device, alarm, and/or user in a [end,
// start) time range, in descending timestamp order.
func (a *Alert) ListAlerts(ctx context.Context, req *api.ListAlertsRequest) (
	*api.ListAlertsResponse, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_VIEWER {
		return nil, errPerm(api.Role_VIEWER)
	}

	var uniqID string
	var devID string

	switch v := req.GetIdOneof().(type) {
	case *api.ListAlertsRequest_UniqId:
		uniqID = v.UniqId
	case *api.ListAlertsRequest_DeviceId:
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

	alerts, err := a.aleDAO.List(ctx, sess.OrgID, uniqID, devID, req.GetAlarmId(),
		req.GetUserId(), end, start)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.ListAlertsResponse{Alerts: alerts}, nil
}
