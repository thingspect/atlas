package service

//go:generate mockgen -source rule.go -destination mock_ruler_alarmer_test.go -package service

import (
	"context"
	"time"

	"github.com/mennanov/fmutils"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alarm"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/rule"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Ruler defines the methods provided by a rule.DAO.
type Ruler interface {
	Create(ctx context.Context, rule *common.Rule) (*common.Rule, error)
	Read(ctx context.Context, ruleID, orgID string) (*common.Rule, error)
	Update(ctx context.Context, rule *common.Rule) (*common.Rule, error)
	Delete(ctx context.Context, ruleID, orgID string) error
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32) ([]*common.Rule, int32, error)
}

// Alarmer defines the methods provided by a alarm.DAO.
type Alarmer interface {
	Create(ctx context.Context, alarm *api.Alarm) (*api.Alarm, error)
	Read(ctx context.Context, alarmID, orgID, ruleID string) (*api.Alarm, error)
	Update(ctx context.Context, alarm *api.Alarm) (*api.Alarm, error)
	Delete(ctx context.Context, alarmID, orgID, ruleID string) error
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32, ruleID string) ([]*api.Alarm, int32, error)
}

// Rule service contains functions to query and modify rules.
type Rule struct {
	api.UnimplementedRuleServiceServer

	ruleDAO  Ruler
	alarmDAO Alarmer
}

// NewRule instantiates and returns a new Rule service.
func NewRule(ruleDAO Ruler, alarmDAO Alarmer) *Rule {
	return &Rule{
		ruleDAO:  ruleDAO,
		alarmDAO: alarmDAO,
	}
}

// CreateRule creates a rule.
func (r *Rule) CreateRule(ctx context.Context,
	req *api.CreateRuleRequest) (*common.Rule, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	req.Rule.OrgId = sess.OrgID

	rule, err := r.ruleDAO.Create(ctx, req.Rule)
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"201")); err != nil {
		logger.Errorf("CreateRule grpc.SetHeader: %v", err)
	}

	return rule, nil
}

// CreateAlarm creates an alarm.
func (r *Rule) CreateAlarm(ctx context.Context,
	req *api.CreateAlarmRequest) (*api.Alarm, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	req.Alarm.OrgId = sess.OrgID

	alarm, err := r.alarmDAO.Create(ctx, req.Alarm)
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"201")); err != nil {
		logger.Errorf("CreateAlarm grpc.SetHeader: %v", err)
	}

	return alarm, nil
}

// GetRule retrieves a rule by ID.
func (r *Rule) GetRule(ctx context.Context,
	req *api.GetRuleRequest) (*common.Rule, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_VIEWER {
		return nil, errPerm(common.Role_VIEWER)
	}

	rule, err := r.ruleDAO.Read(ctx, req.Id, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return rule, nil
}

// GetAlarm retrieves an alarm by ID.
func (r *Rule) GetAlarm(ctx context.Context,
	req *api.GetAlarmRequest) (*api.Alarm, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_VIEWER {
		return nil, errPerm(common.Role_VIEWER)
	}

	alarm, err := r.alarmDAO.Read(ctx, req.Id, sess.OrgID, req.RuleId)
	if err != nil {
		return nil, errToStatus(err)
	}

	return alarm, nil
}

// UpdateRule updates a rule. Update actions validate after merge to support
// partial updates.
func (r *Rule) UpdateRule(ctx context.Context,
	req *api.UpdateRuleRequest) (*common.Rule, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	if req.Rule == nil {
		return nil, status.Error(codes.InvalidArgument, req.Validate().Error())
	}
	req.Rule.OrgId = sess.OrgID

	// Perform partial update if directed.
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		// Normalize and validate field mask.
		req.UpdateMask.Normalize()
		if !req.UpdateMask.IsValid(req.Rule) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		rule, err := r.ruleDAO.Read(ctx, req.Rule.Id, sess.OrgID)
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.Rule, req.UpdateMask.Paths)
		proto.Merge(rule, req.Rule)
		req.Rule = rule
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rule, err := r.ruleDAO.Update(ctx, req.Rule)
	if err != nil {
		return nil, errToStatus(err)
	}

	return rule, nil
}

// UpdateAlarm updates an alarm. Update actions validate after merge to support
// partial updates.
func (r *Rule) UpdateAlarm(ctx context.Context,
	req *api.UpdateAlarmRequest) (*api.Alarm, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	if req.Alarm == nil {
		return nil, status.Error(codes.InvalidArgument, req.Validate().Error())
	}
	req.Alarm.OrgId = sess.OrgID

	// Perform partial update if directed.
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		// Normalize and validate field mask.
		req.UpdateMask.Normalize()
		if !req.UpdateMask.IsValid(req.Alarm) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		alarm, err := r.alarmDAO.Read(ctx, req.Alarm.Id, sess.OrgID,
			req.Alarm.RuleId)
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.Alarm, req.UpdateMask.Paths)
		proto.Merge(alarm, req.Alarm)
		req.Alarm = alarm
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	alarm, err := r.alarmDAO.Update(ctx, req.Alarm)
	if err != nil {
		return nil, errToStatus(err)
	}

	return alarm, nil
}

// DeleteRule deletes a rule by ID.
func (r *Rule) DeleteRule(ctx context.Context,
	req *api.DeleteRuleRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	if err := r.ruleDAO.Delete(ctx, req.Id, sess.OrgID); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"204")); err != nil {
		logger.Errorf("DeleteRule grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeleteAlarm deletes an alarm by ID.
func (r *Rule) DeleteAlarm(ctx context.Context,
	req *api.DeleteAlarmRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	if err := r.alarmDAO.Delete(ctx, req.Id, sess.OrgID,
		req.RuleId); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"204")); err != nil {
		logger.Errorf("DeleteAlarm grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListRules retrieves all rules.
func (r *Rule) ListRules(ctx context.Context,
	req *api.ListRulesRequest) (*api.ListRulesResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_VIEWER {
		return nil, errPerm(common.Role_VIEWER)
	}

	if req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}

	lBoundTS, prevID, err := session.ParsePageToken(req.PageToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	rules, count, err := r.ruleDAO.List(ctx, sess.OrgID, lBoundTS, prevID,
		req.PageSize+1)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListRulesResponse{Rules: rules, TotalSize: count}

	// Populate next page token.
	if len(rules) == int(req.PageSize+1) {
		resp.Rules = rules[:len(rules)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			rules[len(rules)-2].CreatedAt.AsTime(),
			rules[len(rules)-2].Id); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger.Errorf("ListRules session.GeneratePageToken rule, err: "+
				"%+v, %v", rules[len(rules)-2], err)
		}
	}

	return resp, nil
}

// ListAlarms retrieves alarms.
func (r *Rule) ListAlarms(ctx context.Context,
	req *api.ListAlarmsRequest) (*api.ListAlarmsResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_VIEWER {
		return nil, errPerm(common.Role_VIEWER)
	}

	if req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}

	lBoundTS, prevID, err := session.ParsePageToken(req.PageToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	alarms, count, err := r.alarmDAO.List(ctx, sess.OrgID, lBoundTS, prevID,
		req.PageSize+1, req.RuleId)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListAlarmsResponse{Alarms: alarms, TotalSize: count}

	// Populate next page token.
	if len(alarms) == int(req.PageSize+1) {
		resp.Alarms = alarms[:len(alarms)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			alarms[len(alarms)-2].CreatedAt.AsTime(),
			alarms[len(alarms)-2].Id); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger.Errorf("ListAlarms session.GeneratePageToken alarm, err: "+
				"%+v, %v", alarms[len(alarms)-2], err)
		}
	}

	return resp, nil
}

// TestRule tests a rule.
func (r *Rule) TestRule(ctx context.Context,
	req *api.TestRuleRequest) (*api.TestRuleResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	if req.Point.Attr != req.Rule.Attr {
		return nil, status.Error(codes.InvalidArgument,
			"data point and rule attribute mismatch")
	}

	// Default to current timestamp if not provided.
	if req.Point.Ts == nil {
		req.Point.Ts = timestamppb.Now()
	}

	res, err := rule.Eval(req.Point, req.Rule.Expr)
	if err != nil {
		// Expr does not provide sentinel errors, always consider errors to be
		// invalid input.
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &api.TestRuleResponse{Result: res}, nil
}

// TestAlarm tests an alarm.
func (r *Rule) TestAlarm(ctx context.Context,
	req *api.TestAlarmRequest) (*api.TestAlarmResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_BUILDER {
		return nil, errPerm(common.Role_BUILDER)
	}

	subj, err := alarm.Generate(req.Point, req.Rule, req.Device,
		req.Alarm.SubjectTemplate)
	if err != nil {
		// Template does not provide sentinel errors, always consider errors to
		// be invalid input.
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	body, err := alarm.Generate(req.Point, req.Rule, req.Device,
		req.Alarm.BodyTemplate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &api.TestAlarmResponse{Result: subj + " - " + body}, nil
}
