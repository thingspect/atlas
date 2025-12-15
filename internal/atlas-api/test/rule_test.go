//go:build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/rule"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateRule(t *testing.T) {
	t.Parallel()

	t.Run("Create valid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx,
			&api.CreateRuleRequest{Rule: rule})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)
		require.NotEqual(t, rule.GetId(), createRule.GetId())
		require.WithinDuration(t, time.Now(), createRule.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createRule.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create invalid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		rule.Attr = "api-rule-" + random.String(40)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx,
			&api.CreateRuleRequest{Rule: rule})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateRuleRequest.Rule: embedded message failed "+
			"validation | caused by: invalid Rule.Attr: value length must be "+
			"at most 40 runes")
	})
}

func TestGetRule(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
	createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
		Rule: random.Rule("api-rule", uuid.NewString()),
	})
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Get rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getRule, err := raCli.GetRule(ctx,
			&api.GetRuleRequest{Id: createRule.GetId()})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, createRule, getRule)
	})

	t.Run("Get rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getRule, err := raCli.GetRule(ctx,
			&api.GetRuleRequest{Id: uuid.NewString()})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = "+
			"dao: object not found")
	})

	t.Run("Gets are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		getRule, err := secCli.GetRule(ctx,
			&api.GetRuleRequest{Id: createRule.GetId()})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = "+
			"dao: object not found")
	})
}

func TestUpdateRule(t *testing.T) {
	t.Parallel()

	t.Run("Update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Name = "api-rule-" + random.String(10)
		createRule.Status = api.Status_DISABLED

		updateRule, err := raCli.UpdateRule(ctx,
			&api.UpdateRuleRequest{Rule: createRule})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.NoError(t, err)
		require.Equal(t, createRule.GetCreatedAt().AsTime(),
			updateRule.GetCreatedAt().AsTime())
		require.True(t, updateRule.GetUpdatedAt().AsTime().After(
			updateRule.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createRule.GetCreatedAt().AsTime(),
			updateRule.GetUpdatedAt().AsTime(), 2*time.Second)

		getRule, err := raCli.GetRule(ctx,
			&api.GetRuleRequest{Id: createRule.GetId()})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateRule, getRule)
	})

	t.Run("Partial update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminKeyGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		part := &api.Rule{Id: createRule.GetId(), Name: "api-rule-" +
			random.String(10), Status: api.Status_DISABLED, Expr: `false`}

		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status", "expr"},
			},
		})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.NoError(t, err)
		require.Equal(t, createRule.GetCreatedAt().AsTime(),
			updateRule.GetCreatedAt().AsTime())
		require.True(t, updateRule.GetUpdatedAt().AsTime().After(
			updateRule.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createRule.GetCreatedAt().AsTime(),
			updateRule.GetUpdatedAt().AsTime(), 2*time.Second)

		getRule, err := raCli.GetRule(ctx,
			&api.GetRuleRequest{Id: createRule.GetId()})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateRule, getRule)
	})

	t.Run("Update rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Update nil rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx,
			&api.UpdateRuleRequest{Rule: nil})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateRuleRequest.Rule: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule:       random.Rule("api-rule", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"aaa"}},
		})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update rule by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status", "expr"},
			},
		})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = "+
			"dao: object not found")
	})

	t.Run("Update rule by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = "+
			"dao: object not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.OrgId = uuid.NewString()
		createRule.Name = "api-rule-" + random.String(10)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		updateRule, err := secCli.UpdateRule(ctx,
			&api.UpdateRuleRequest{Rule: createRule})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = "+
			"dao: object not found")
	})

	t.Run("Update rule validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Attr = "api-rule-" + random.String(40)

		updateRule, err := raCli.UpdateRule(ctx,
			&api.UpdateRuleRequest{Rule: createRule})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateRuleRequest.Rule: embedded message failed "+
			"validation | caused by: invalid Rule.Attr: value length must be "+
			"at most 40 runes")
	})
}

func TestDeleteRule(t *testing.T) {
	t.Parallel()

	t.Run("Delete rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		_, err = raCli.DeleteRule(ctx,
			&api.DeleteRuleRequest{Id: createRule.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read rule by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()

			raCli := api.NewRuleAlarmServiceClient(globalAdminKeyGRPCConn)
			getRule, err := raCli.GetRule(ctx,
				&api.GetRuleRequest{Id: createRule.GetId()})
			t.Logf("getRule, err: %+v, %v", getRule, err)
			require.Nil(t, getRule)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"dao: object not found")
		})
	})

	t.Run("Delete rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		_, err := raCli.DeleteRule(ctx,
			&api.DeleteRuleRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		_, err := raCli.DeleteRule(ctx,
			&api.DeleteRuleRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = "+
			"dao: object not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteRule(ctx,
			&api.DeleteRuleRequest{Id: createRule.GetId()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = "+
			"dao: object not found")
	})
}

func TestListRules(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	ruleIDs := []string{}
	ruleNames := []string{}
	ruleStatuses := []api.Status{}
	for range 3 {
		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		ruleIDs = append(ruleIDs, createRule.GetId())
		ruleNames = append(ruleNames, createRule.GetName())
		ruleStatuses = append(ruleStatuses, createRule.GetStatus())
	}

	t.Run("List rules by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listRules, err := raCli.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listRules.GetRules()), 3)
		require.GreaterOrEqual(t, listRules.GetTotalSize(), int32(3))

		var found bool
		for _, rule := range listRules.GetRules() {
			if rule.GetId() == ruleIDs[len(ruleIDs)-1] &&
				rule.GetName() == ruleNames[len(ruleNames)-1] &&
				rule.GetStatus() == ruleStatuses[len(ruleStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List rules by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminKeyGRPCConn)
		listRules, err := raCli.ListRules(ctx,
			&api.ListRulesRequest{PageSize: 2})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Len(t, listRules.GetRules(), 2)
		require.NotEmpty(t, listRules.GetNextPageToken())
		require.GreaterOrEqual(t, listRules.GetTotalSize(), int32(3))

		nextRules, err := raCli.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2, PageToken: listRules.GetNextPageToken(),
		})
		t.Logf("nextRules, err: %+v, %v", nextRules, err)
		require.NoError(t, err)
		require.NotEmpty(t, nextRules.GetRules())
		require.GreaterOrEqual(t, nextRules.GetTotalSize(), int32(3))
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		listRules, err := secCli.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Empty(t, listRules.GetRules())
		require.Equal(t, int32(0), listRules.GetTotalSize())
	})

	t.Run("List rules by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listRules, err := raCli.ListRules(ctx,
			&api.ListRulesRequest{PageToken: badUUID})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
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
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, `true`, true, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, `10 > 15`, false, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, `point.Token == ""`, true, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}, Ts: timestamppb.New(time.Now().Add(-time.Second))},
				`pointTS < currTS`, true, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, `point.GetIntVal() == 40`, true, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, `pointVal > 32`, true, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_Fl64Val{
					Fl64Val: 37.7,
				}}, `pointVal < 32`, false, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_StrVal{
					StrVal: "line",
				}}, `pointVal == battery`, false, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_BoolVal{
					BoolVal: true,
				}}, `pointVal`, true, "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, `1 + "aaa"`, false,
				"invalid operation: int + string",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, `"aaa"`, false, rule.ErrNotBool.Error(),
			},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Can evaluate %+v", test), func(t *testing.T) {
				t.Parallel()

				test.inpPoint.UniqId = "api-rule-" + random.String(16)
				test.inpPoint.Attr = "api-rule" + random.String(10)

				rule := random.Rule("api-rule", uuid.NewString())
				rule.Attr = test.inpPoint.GetAttr()
				rule.Expr = test.inpRuleExpr

				ctx, cancel := context.WithTimeout(t.Context(),
					testTimeout)
				defer cancel()

				raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
				testRes, err := raCli.TestRule(ctx, &api.TestRuleRequest{
					Point: test.inpPoint, Rule: rule,
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

	t.Run("Test rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-rule-" + random.String(16), Attr: "api-rule" +
				random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
		}
		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		testRes, err := raCli.TestRule(ctx,
			&api.TestRuleRequest{Point: point, Rule: rule})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Test invalid rule attribute", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-rule-" + random.String(16), Attr: "api-rule" +
				random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
		}
		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		testRes, err := raCli.TestRule(ctx,
			&api.TestRuleRequest{Point: point, Rule: rule})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"data point and rule attribute mismatch")
	})
}
