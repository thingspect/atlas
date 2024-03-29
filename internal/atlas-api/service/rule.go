package service

//go:generate mockgen -source rule.go -destination mock_ruler_test.go -package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/mennanov/fmutils"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/rule"
	"github.com/thingspect/proto/go/api"
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
	Create(ctx context.Context, rule *api.Rule) (*api.Rule, error)
	Read(ctx context.Context, ruleID, orgID string) (*api.Rule, error)
	Update(ctx context.Context, rule *api.Rule) (*api.Rule, error)
	Delete(ctx context.Context, ruleID, orgID string) error
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32) ([]*api.Rule, int32, error)
}

// RuleAlarm service contains functions to query and modify rules and alarms.
type RuleAlarm struct {
	api.UnimplementedRuleAlarmServiceServer

	ruleDAO  Ruler
	alarmDAO Alarmer
}

// NewRuleAlarm instantiates and returns a new RuleAlarm service.
func NewRuleAlarm(ruleDAO Ruler, alarmDAO Alarmer) *RuleAlarm {
	return &RuleAlarm{
		ruleDAO:  ruleDAO,
		alarmDAO: alarmDAO,
	}
}

// CreateRule creates a rule.
func (ra *RuleAlarm) CreateRule(
	ctx context.Context, req *api.CreateRuleRequest,
) (*api.Rule, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	req.Rule.OrgId = sess.OrgID

	rule, err := ra.ruleDAO.Create(ctx, req.GetRule())
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusCreated))); err != nil {
		logger := alog.FromContext(ctx)
		logger.Errorf("CreateRule grpc.SetHeader: %v", err)
	}

	return rule, nil
}

// GetRule retrieves a rule by ID.
func (ra *RuleAlarm) GetRule(ctx context.Context, req *api.GetRuleRequest) (
	*api.Rule, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_VIEWER {
		return nil, errPerm(api.Role_VIEWER)
	}

	rule, err := ra.ruleDAO.Read(ctx, req.GetId(), sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return rule, nil
}

// UpdateRule updates a rule. Update actions validate after merge to support
// partial updates.
func (ra *RuleAlarm) UpdateRule(
	ctx context.Context, req *api.UpdateRuleRequest,
) (*api.Rule, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	if req.GetRule() == nil {
		return nil, status.Error(codes.InvalidArgument,
			req.Validate().Error())
	}
	req.Rule.OrgId = sess.OrgID

	// Perform partial update if directed.
	if len(req.GetUpdateMask().GetPaths()) > 0 {
		// Normalize and validate field mask.
		req.GetUpdateMask().Normalize()
		if !req.GetUpdateMask().IsValid(req.GetRule()) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		rule, err := ra.ruleDAO.Read(ctx, req.GetRule().GetId(), sess.OrgID)
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.GetRule(), req.GetUpdateMask().GetPaths())
		proto.Merge(rule, req.GetRule())
		req.Rule = rule
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rule, err := ra.ruleDAO.Update(ctx, req.GetRule())
	if err != nil {
		return nil, errToStatus(err)
	}

	return rule, nil
}

// DeleteRule deletes a rule by ID.
func (ra *RuleAlarm) DeleteRule(
	ctx context.Context, req *api.DeleteRuleRequest,
) (*emptypb.Empty, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	if err := ra.ruleDAO.Delete(ctx, req.GetId(), sess.OrgID); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusNoContent))); err != nil {
		logger := alog.FromContext(ctx)
		logger.Errorf("DeleteRule grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListRules retrieves all rules.
func (ra *RuleAlarm) ListRules(ctx context.Context, req *api.ListRulesRequest) (
	*api.ListRulesResponse, error,
) {
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
	rules, count, err := ra.ruleDAO.List(ctx, sess.OrgID, lBoundTS, prevID,
		req.GetPageSize()+1)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListRulesResponse{Rules: rules, TotalSize: count}

	// Populate next page token.
	if len(rules) == int(req.GetPageSize()+1) {
		resp.Rules = rules[:len(rules)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			rules[len(rules)-2].GetCreatedAt().AsTime(),
			rules[len(rules)-2].GetId()); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger := alog.FromContext(ctx)
			logger.Errorf("ListRules session.GeneratePageToken rule, err: "+
				"%+v, %v", rules[len(rules)-2], err)
		}
	}

	return resp, nil
}

// TestRule tests a rule.
func (ra *RuleAlarm) TestRule(ctx context.Context, req *api.TestRuleRequest) (
	*api.TestRuleResponse, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_BUILDER {
		return nil, errPerm(api.Role_BUILDER)
	}

	if req.GetPoint().GetAttr() != req.GetRule().GetAttr() {
		return nil, status.Error(codes.InvalidArgument,
			"data point and rule attribute mismatch")
	}

	// Default to current timestamp if not provided.
	if req.GetPoint().GetTs() == nil {
		req.Point.Ts = timestamppb.Now()
	}

	res, err := rule.Eval(req.GetPoint(), req.GetRule().GetExpr())
	if err != nil {
		// Expr does not provide sentinel errors, always consider errors to be
		// invalid input.
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &api.TestRuleResponse{Result: res}, nil
}
