//go:build !integration

package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/rule"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"go.uber.org/mock/gomock"
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
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		createRule, err := raSvc.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: rule,
		})
		t.Logf("rule, createRule, err: %+v, %+v, %v", rule, createRule, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, rule, createRule)
	})

	t.Run("Create rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		createRule, err := raSvc.CreateRule(ctx, &api.CreateRuleRequest{})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		createRule, err := raSvc.CreateRule(ctx, &api.CreateRuleRequest{})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create invalid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		rule.Attr = random.String(41)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Create(gomock.Any(), rule).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		createRule, err := raSvc.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: rule,
		})
		t.Logf("rule, createRule, err: %+v, %+v, %v", rule, createRule, err)
		require.Nil(t, createRule)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"dao: invalid format"), err)
	})
}

func TestGetRule(t *testing.T) {
	t.Parallel()

	t.Run("Get rule by valid ID", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		retRule, _ := proto.Clone(rule).(*api.Rule)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), rule.GetId(), rule.GetOrgId()).Return(retRule,
			nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		getRule, err := raSvc.GetRule(ctx, &api.GetRuleRequest{Id: rule.GetId()})
		t.Logf("rule, getRule, err: %+v, %+v, %v", rule, getRule, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, rule, getRule)
	})

	t.Run("Get rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		getRule, err := raSvc.GetRule(ctx, &api.GetRuleRequest{})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Get rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		getRule, err := raSvc.GetRule(ctx, &api.GetRuleRequest{})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Get rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		getRule, err := raSvc.GetRule(ctx, &api.GetRuleRequest{
			Id: uuid.NewString(),
		})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.Equal(t, status.Error(codes.NotFound, "dao: object not found"),
			err)
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
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule,
		})
		t.Logf("rule, updateRule, err: %+v, %+v, %v", rule, updateRule, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, rule, updateRule)
	})

	t.Run("Partial update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		retRule, _ := proto.Clone(rule).(*api.Rule)
		part := &api.Rule{
			Id: rule.GetId(), Status: api.Status_ACTIVE, Expr: `true`,
		}
		merged := &api.Rule{
			Id: rule.GetId(), OrgId: rule.GetOrgId(), Name: rule.GetName(),
			Status: part.GetStatus(), DeviceTag: rule.GetDeviceTag(), Attr: rule.GetAttr(),
			Expr: part.GetExpr(),
		}
		retMerged, _ := proto.Clone(merged).(*api.Rule)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), rule.GetId(), rule.GetOrgId()).Return(retRule,
			nil).Times(1)
		ruler.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(merged)).
			Return(retMerged, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status", "expr"},
			},
		})
		t.Logf("merged, updateRule, err: %+v, %+v, %v", merged, updateRule, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, merged, updateRule)
	})

	t.Run("Update rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Update rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Update nil rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: nil,
		})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateRuleRequest.Rule: value is required"), err)
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"},
			},
		})
		t.Logf("rule, updateRule, err: %+v, %+v, %v", rule, updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid field mask"), err)
	})

	t.Run("Partial update rule by unknown rule", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		part := &api.Rule{Id: uuid.NewString(), Status: api.Status_ACTIVE}

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Read(gomock.Any(), part.GetId(), orgID).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"},
			},
		})
		t.Logf("part, updateRule, err: %+v, %+v, %v", part, updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.NotFound, "dao: object not found"),
			err)
	})

	t.Run("Update rule validation failure", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		rule.Attr = random.String(41)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule,
		})
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
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		updateRule, err := raSvc.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: rule,
		})
		t.Logf("rule, updateRule, err: %+v, %+v, %v", rule, updateRule, err)
		require.Nil(t, updateRule)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"dao: invalid format"), err)
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
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		_, err := raSvc.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		_, err := raSvc.DeleteRule(ctx, &api.DeleteRuleRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		_, err := raSvc.DeleteRule(ctx, &api.DeleteRuleRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		_, err := raSvc.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "dao: object not found"),
			err)
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
			t.Context(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		listRules, err := raSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listRules.GetTotalSize())
		require.EqualExportedValues(t,
			&api.ListRulesResponse{Rules: rules, TotalSize: 3}, listRules)
	})

	t.Run("List rules by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		rules := []*api.Rule{
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
			random.Rule("api-rule", uuid.NewString()),
		}

		next, err := session.GeneratePageToken(rules[1].GetCreatedAt().AsTime(),
			rules[1].GetId())
		require.NoError(t, err)

		ruler := NewMockRuler(gomock.NewController(t))
		ruler.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(rules, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		listRules, err := raSvc.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2,
		})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listRules.GetTotalSize())
		require.EqualExportedValues(t, &api.ListRulesResponse{
			Rules: rules[:2], NextPageToken: next, TotalSize: 3,
		}, listRules)
	})

	t.Run("List rules with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		listRules, err := raSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List rules with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		listRules, err := raSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List rules by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		listRules, err := raSvc.ListRules(ctx, &api.ListRulesRequest{
			PageToken: badUUID,
		})
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
			t.Context(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		listRules, err := raSvc.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"dao: invalid format"), err)
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
			t.Context(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(ruler, nil)
		listRules, err := raSvc.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2,
		})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listRules.GetTotalSize())
		require.EqualExportedValues(t,
			&api.ListRulesResponse{Rules: rules[:2], TotalSize: 3}, listRules)
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
			{
				&common.DataPoint{}, `true`, true, "",
			},
			{
				&common.DataPoint{}, `10 > 15`, false, "",
			},
			{
				&common.DataPoint{}, `point.Token == ""`, true, "",
			},
			{
				&common.DataPoint{
					Ts: timestamppb.New(time.Now().Add(-time.Second)),
				}, `pointTS < currTS`, true, "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_IntVal{IntVal: 40},
				}, `point.GetIntVal() == 40`, true, "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_IntVal{IntVal: 40},
				}, `pointVal > 32`, true, "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 37.7},
				}, `pointVal < 32`, false, "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"},
				}, `pointVal == battery`, false, "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_BoolVal{BoolVal: true},
				}, `pointVal`, true, "",
			},
			{
				&common.DataPoint{}, `1 + "aaa"`, false,
				"invalid operation: int + string",
			},
			{
				&common.DataPoint{}, `"aaa"`, false, rule.ErrNotBool.Error(),
			},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Can evaluate %+v", test), func(t *testing.T) {
				t.Parallel()

				test.inpPoint.Attr = "api-rule" + random.String(10)

				ctx, cancel := context.WithTimeout(session.NewContext(
					t.Context(), &session.Session{
						OrgID: uuid.NewString(), Role: api.Role_ADMIN,
					}),
					testTimeout)
				defer cancel()

				raSvc := NewRuleAlarm(nil, nil)
				testRes, err := raSvc.TestRule(ctx, &api.TestRuleRequest{
					Point: test.inpPoint, Rule: &api.Rule{
						Attr: test.inpPoint.GetAttr(), Expr: test.inpRuleExpr,
					},
				})
				t.Logf("testRes, err: %+v, %v", testRes, err)
				if test.err == "" {
					require.Equal(t, test.res, testRes.GetResult())
					require.NoError(t, err)
				} else {
					require.Nil(t, testRes)
					require.Equal(t, codes.InvalidArgument, status.Code(err))
					require.Contains(t, err.Error(), test.err)
				}
			})
		}
	})

	t.Run("Test rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		testRes, err := raSvc.TestRule(ctx, &api.TestRuleRequest{})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Test rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		testRes, err := raSvc.TestRule(ctx, &api.TestRuleRequest{})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Test invalid rule attribute", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{Attr: "api-rule" + random.String(10)}
		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: rule.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		testRes, err := raSvc.TestRule(ctx, &api.TestRuleRequest{
			Point: point, Rule: rule,
		})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"data point and rule attribute mismatch"), err)
	})
}
