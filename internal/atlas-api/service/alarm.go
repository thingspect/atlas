package service

//go:generate mockgen -source alarm.go -destination mock_alarmer_test.go -package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/mennanov/fmutils"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/template"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Alarmer defines the methods provided by a alarm.DAO.
type Alarmer interface {
	Create(ctx context.Context, alarm *api.Alarm) (*api.Alarm, error)
	Read(ctx context.Context, alarmID, orgID, ruleID string) (*api.Alarm, error)
	Update(ctx context.Context, alarm *api.Alarm) (*api.Alarm, error)
	Delete(ctx context.Context, alarmID, orgID, ruleID string) error
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32, ruleID string) ([]*api.Alarm, int32, error)
}

// CreateAlarm creates an alarm.
func (ra *RuleAlarm) CreateAlarm(
	ctx context.Context, req *api.CreateAlarmRequest,
) (*api.Alarm, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	req.Alarm.OrgId = sess.OrgID

	alarm, err := ra.alarmDAO.Create(ctx, req.GetAlarm())
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusCreated))); err != nil {
		logger := alog.FromContext(ctx)
		logger.Errorf("CreateAlarm grpc.SetHeader: %v", err)
	}

	return alarm, nil
}

// GetAlarm retrieves an alarm by ID.
func (ra *RuleAlarm) GetAlarm(ctx context.Context, req *api.GetAlarmRequest) (
	*api.Alarm, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_VIEWER {
		return nil, errPerm(api.Role_VIEWER)
	}

	alarm, err := ra.alarmDAO.Read(ctx, req.GetId(), sess.OrgID, req.GetRuleId())
	if err != nil {
		return nil, errToStatus(err)
	}

	return alarm, nil
}

// UpdateAlarm updates an alarm. Update actions validate after merge to support
// partial updates.
func (ra *RuleAlarm) UpdateAlarm(
	ctx context.Context, req *api.UpdateAlarmRequest,
) (*api.Alarm, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	if req.GetAlarm() == nil {
		return nil, status.Error(codes.InvalidArgument,
			req.Validate().Error())
	}
	req.Alarm.OrgId = sess.OrgID

	// Perform partial update if directed.
	if len(req.GetUpdateMask().GetPaths()) > 0 {
		// Normalize and validate field mask.
		req.GetUpdateMask().Normalize()
		if !req.GetUpdateMask().IsValid(req.GetAlarm()) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		alarm, err := ra.alarmDAO.Read(ctx, req.GetAlarm().GetId(), sess.OrgID,
			req.GetAlarm().GetRuleId())
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.GetAlarm(), req.GetUpdateMask().GetPaths())
		if req.GetAlarm().GetUserTags() != nil {
			alarm.UserTags = nil
		}
		proto.Merge(alarm, req.GetAlarm())
		req.Alarm = alarm
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	alarm, err := ra.alarmDAO.Update(ctx, req.GetAlarm())
	if err != nil {
		return nil, errToStatus(err)
	}

	return alarm, nil
}

// DeleteAlarm deletes an alarm by ID.
func (ra *RuleAlarm) DeleteAlarm(
	ctx context.Context, req *api.DeleteAlarmRequest,
) (*emptypb.Empty, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	if err := ra.alarmDAO.Delete(ctx, req.GetId(), sess.OrgID,
		req.GetRuleId()); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusNoContent))); err != nil {
		logger := alog.FromContext(ctx)
		logger.Errorf("DeleteAlarm grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListAlarms retrieves alarms.
func (ra *RuleAlarm) ListAlarms(
	ctx context.Context, req *api.ListAlarmsRequest,
) (*api.ListAlarmsResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_VIEWER {
		return nil, errPerm(api.Role_VIEWER)
	}

	if req.GetPageSize() == 0 {
		req.PageSize = defaultPageSize
	}

	lBoundTS, prevID, err := session.ParsePageToken(req.GetPageToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	alarms, count, err := ra.alarmDAO.List(ctx, sess.OrgID, lBoundTS, prevID,
		req.GetPageSize()+1, req.GetRuleId())
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListAlarmsResponse{Alarms: alarms, TotalSize: count}

	// Populate next page token.
	if len(alarms) == int(req.GetPageSize()+1) {
		resp.Alarms = alarms[:len(alarms)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			alarms[len(alarms)-2].GetCreatedAt().AsTime(),
			alarms[len(alarms)-2].GetId()); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger := alog.FromContext(ctx)
			logger.Errorf("ListAlarms session.GeneratePageToken alarm, err: "+
				"%+v, %v", alarms[len(alarms)-2], err)
		}
	}

	return resp, nil
}

// TestAlarm tests an alarm.
func (ra *RuleAlarm) TestAlarm(ctx context.Context, req *api.TestAlarmRequest) (
	*api.TestAlarmResponse, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	subj, err := template.Generate(req.GetPoint(), req.GetRule(), req.GetDevice(),
		req.GetAlarm().GetSubjectTemplate())
	if err != nil {
		// Template does not provide sentinel errors, always consider errors to
		// be invalid input.
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	body, err := template.Generate(req.GetPoint(), req.GetRule(), req.GetDevice(),
		req.GetAlarm().GetBodyTemplate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &api.TestAlarmResponse{Result: subj + " - " + body}, nil
}
