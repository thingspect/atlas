// +build !integration

package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/rule"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateRule(t *testing.T) {
	t.Parallel()

	t.Run("Create valid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		retRule, _ := proto.Clone(rule).(*api.Rule)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Create(gomock.Any(), rule).Return(retRule, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		createRule, err := ruleSvc.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: rule})
		t.Logf("rule, createRule, err: %+v, %+v, %v", rule, createRule, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(rule, createRule) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", rule, createRule)
		}
	})

	t.Run("Create rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		createRule, err := ruleSvc.CreateRule(ctx, &api.CreateRuleRequest{})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Create rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_VIEWER}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		createRule, err := ruleSvc.CreateRule(ctx, &api.CreateRuleRequest{})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Create invalid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		rule.Attr = random.String(41)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Create(gomock.Any(), rule).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		createRule, err := ruleSvc.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: rule})
		t.Logf("rule, createRule, err: %+v, %+v, %v", rule, createRule, err)
		require.Nil(t, createRule)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestGetRule(t *testing.T) {
	t.Parallel()

	t.Run("Get rule by valid ID", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		retRule, _ := proto.Clone(rule).(*api.Rule)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), rule.Id, rule.OrgId).Return(retRule,
			nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		getRule, err := ruleSvc.GetRule(ctx, &api.GetRuleRequest{Id: rule.Id})
		t.Logf("rule, getRule, err: %+v, %+v, %v", rule, getRule, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(rule, getRule) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", rule, getRule)
		}
	})

	t.Run("Get rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		getRule, err := ruleSvc.GetRule(ctx, &api.GetRuleRequest{})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("Get rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_CONTACT}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		getRule, err := ruleSvc.GetRule(ctx, &api.GetRuleRequest{})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("Get rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		getRule, err := ruleSvc.GetRule(ctx, &api.GetRuleRequest{
			Id: uuid.NewString()})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestUpdateRule(t *testing.T) {
	t.Parallel()

	t.Run("Update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		retRule, _ := proto.Clone(rule).(*api.Rule)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Update(gomock.Any(), rule).Return(retRule, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule})
		t.Logf("rule, updateRule, err: %+v, %+v, %v", rule, updateRule, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(rule, updateRule) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", rule, updateRule)
		}
	})

	t.Run("Partial update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		retRule, _ := proto.Clone(rule).(*api.Rule)
		part := &api.Rule{Id: rule.Id, Status: common.Status_ACTIVE,
			Expr: `true`}
		merged := &api.Rule{Id: rule.Id, OrgId: rule.OrgId, Name: rule.Name,
			Status: part.Status, DeviceTag: rule.DeviceTag, Attr: rule.Attr,
			Expr: part.Expr}
		retMerged, _ := proto.Clone(merged).(*api.Rule)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), rule.Id, rule.OrgId).Return(retRule,
			nil).Times(1)
		ruler.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(merged)).
			Return(retMerged, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status", "expr"}}})
		t.Logf("merged, updateRule, err: %+v, %+v, %v", merged, updateRule, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(merged, updateRule) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", merged, updateRule)
		}
	})

	t.Run("Update rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Update rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_VIEWER}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Update nil rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: nil})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateRuleRequest.Rule: value is required"), err)
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"}}})
		t.Logf("rule, updateRule, err: %+v, %+v, %v", rule, updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid field mask"), err)
	})

	t.Run("Partial update rule by unknown rule", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		part := &api.Rule{Id: uuid.NewString(), Status: common.Status_ACTIVE}

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), part.Id, orgID).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"}}})
		t.Logf("part, updateRule, err: %+v, %+v, %v", part, updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Update rule validation failure", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		rule.Attr = random.String(41)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule})
		t.Logf("rule, updateRule, err: %+v, %+v, %v", rule, updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateRuleRequest.Rule: embedded message failed "+
				"validation | caused by: invalid Rule.Attr: value length must "+
				"be at most 40 runes"), err)
	})

	t.Run("Update rule by invalid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Update(gomock.Any(), rule).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		updateRule, err := ruleSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule})
		t.Logf("rule, updateRule, err: %+v, %+v, %v", rule, updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestDeleteRule(t *testing.T) {
	t.Parallel()

	t.Run("Delete rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		_, err := ruleSvc.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		_, err := ruleSvc.DeleteRule(ctx, &api.DeleteRuleRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Delete rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_VIEWER}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		_, err := ruleSvc.DeleteRule(ctx, &api.DeleteRuleRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Delete rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		_, err := ruleSvc.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestListRules(t *testing.T) {
	t.Parallel()

	t.Run("List rules by valid org ID", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		rules := []*api.Rule{
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
		}

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(51)).
			Return(rules, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		listRules, err := ruleSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listRules.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListRulesResponse{Rules: rules, TotalSize: 3},
			listRules) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListRulesResponse{Rules: rules, TotalSize: 3}, listRules)
		}
	})

	t.Run("List rules by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		rules := []*api.Rule{
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
		}

		next, err := session.GeneratePageToken(rules[1].CreatedAt.AsTime(),
			rules[1].Id)
		require.NoError(t, err)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(rules, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		listRules, err := ruleSvc.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listRules.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListRulesResponse{Rules: rules[:2],
			NextPageToken: next, TotalSize: 3}, listRules) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListRulesResponse{Rules: rules[:2], NextPageToken: next,
					TotalSize: 3}, listRules)
		}
	})

	t.Run("List rules with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		listRules, err := ruleSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("List rules with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_CONTACT}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		listRules, err := ruleSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("List rules by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		listRules, err := ruleSvc.ListRules(ctx, &api.ListRulesRequest{
			PageToken: badUUID})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid page token"), err)
	})

	t.Run("List rules by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil, int32(0), dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: "aaa",
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		listRules, err := ruleSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("List rules with generation failure", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		rules := []*api.Rule{
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
		}
		rules[1].Id = badUUID

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(rules, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(ruler)
		listRules, err := ruleSvc.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listRules.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListRulesResponse{Rules: rules[:2],
			TotalSize: 3}, listRules) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListRulesResponse{Rules: rules[:2], TotalSize: 3},
				listRules)
		}
	})
}

func TestTestRule(t *testing.T) {
	t.Parallel()

	t.Run("Test valid and invalid rules", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			inpPoint    *common.DataPoint
			inpRuleExpr string
			res         bool
			err         string
		}{
			{&common.DataPoint{}, `true`, true, ""},
			{&common.DataPoint{}, `10 > 15`, false, ""},
			{&common.DataPoint{}, `point.Token == ""`, true, ""},
			{&common.DataPoint{Ts: timestamppb.New(time.Now().
				Add(-time.Second))}, `pointTS < currTS`, true, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				`point.GetIntVal() == 40`, true, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				`pointVal > 32`, true, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_Fl64Val{
				Fl64Val: 37.7}}, `pointVal < 32`, false, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_StrVal{
				StrVal: "batt"}}, `pointVal == line`, false, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_BoolVal{
				BoolVal: true}}, `pointVal`, true, ""},
			{&common.DataPoint{}, `1 + "aaa"`, false,
				"invalid operation: int + string"},
			{&common.DataPoint{}, `"aaa"`, false, rule.ErrNotBool.Error()},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Can evaluate %+v", lTest), func(t *testing.T) {
				t.Parallel()

				lTest.inpPoint.Attr = "api-rule" + random.String(10)

				ctx, cancel := context.WithTimeout(session.NewContext(
					context.Background(), &session.Session{
						OrgID: uuid.NewString(), Role: common.Role_ADMIN}),
					testTimeout)
				defer cancel()

				ruleSvc := NewRule(nil)
				testRes, err := ruleSvc.TestRule(ctx, &api.TestRuleRequest{
					Point: lTest.inpPoint, Rule: &api.Rule{
						Attr: lTest.inpPoint.Attr, Expr: lTest.inpRuleExpr}})
				t.Logf("testRes, err: %+v, %v", testRes, err)
				if lTest.err == "" {
					require.Equal(t, lTest.res, testRes.Result)
					require.NoError(t, err)
				} else {
					require.Nil(t, testRes)
					require.Equal(t, codes.InvalidArgument, status.Code(err))
					require.Contains(t, err.Error(), lTest.err)
				}
			})
		}
	})

	t.Run("Test rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		testRes, err := ruleSvc.TestRule(ctx, &api.TestRuleRequest{})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Test rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_VIEWER}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		testRes, err := ruleSvc.TestRule(ctx, &api.TestRuleRequest{})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, errPerm(common.Role_BUILDER), err)
	})

	t.Run("Test invalid rule attribute", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{Attr: "api-rule" + random.String(10)}
		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: rule.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		ruleSvc := NewRule(nil)
		testRes, err := ruleSvc.TestRule(ctx, &api.TestRuleRequest{
			Point: point, Rule: rule})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"data point and rule attribute mismatch"), err)
	})
}
