// +build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/rule"
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: rule})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)
		require.NotEqual(t, rule.Id, createRule.Id)
		require.WithinDuration(t, time.Now(), createRule.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createRule.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(secondaryViewerGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create invalid rule", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-rule", uuid.NewString())
		rule.Attr = "api-rule-" + random.String(40)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: rule})
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

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
	createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
		Rule: random.Rule("api-rule", uuid.NewString())})
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Get rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		getRule, err := ruleCli.GetRule(ctx, &api.GetRuleRequest{
			Id: createRule.Id})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createRule, getRule) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createRule, getRule)
		}
	})

	t.Run("Get rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		getRule, err := ruleCli.GetRule(ctx, &api.GetRuleRequest{
			Id: uuid.NewString()})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Get are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewRuleServiceClient(secondaryAdminGRPCConn)
		getRule, err := secCli.GetRule(ctx, &api.GetRuleRequest{
			Id: createRule.Id})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestUpdateRule(t *testing.T) {
	t.Parallel()

	t.Run("Update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Name = "api-rule-" + random.String(10)
		createRule.Status = common.Status_DISABLED

		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: createRule})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.NoError(t, err)
		require.Equal(t, createRule.CreatedAt.AsTime(),
			updateRule.CreatedAt.AsTime())
		require.True(t, updateRule.UpdatedAt.AsTime().After(
			updateRule.CreatedAt.AsTime()))
		require.WithinDuration(t, createRule.CreatedAt.AsTime(),
			updateRule.UpdatedAt.AsTime(), 2*time.Second)

		getRule, err := ruleCli.GetRule(ctx, &api.GetRuleRequest{
			Id: createRule.Id})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateRule, getRule) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateRule, getRule)
		}
	})

	t.Run("Partial update rule by valid rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		part := &common.Rule{Id: createRule.Id, Name: "api-rule-" +
			random.String(10), Status: common.Status_DISABLED, Expr: `false`}

		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status", "expr"}}})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.NoError(t, err)
		require.Equal(t, createRule.CreatedAt.AsTime(),
			updateRule.CreatedAt.AsTime())
		require.True(t, updateRule.UpdatedAt.AsTime().After(
			updateRule.CreatedAt.AsTime()))
		require.WithinDuration(t, createRule.CreatedAt.AsTime(),
			updateRule.UpdatedAt.AsTime(), 2*time.Second)

		getRule, err := ruleCli.GetRule(ctx, &api.GetRuleRequest{
			Id: createRule.Id})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateRule, getRule) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateRule, getRule)
		}
	})

	t.Run("Update rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(secondaryViewerGRPCConn)
		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Update nil rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: nil})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateRuleRequest.Rule: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule:       random.Rule("api-rule", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"aaa"}}})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update rule by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString()),
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status", "expr"}}})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update rule by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.OrgId = uuid.NewString()
		createRule.Name = "api-rule-" + random.String(10)

		secCli := api.NewRuleServiceClient(secondaryAdminGRPCConn)
		updateRule, err := secCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: createRule})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update rule validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Attr = "api-rule-" + random.String(40)

		updateRule, err := ruleCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: createRule})
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		_, err = ruleCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: createRule.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read rule by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
			getRule, err := ruleCli.GetRule(ctx, &api.GetRuleRequest{
				Id: createRule.Id})
			t.Logf("getRule, err: %+v, %v", getRule, err)
			require.Nil(t, getRule)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(secondaryViewerGRPCConn)
		_, err := ruleCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		_, err := ruleCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		secCli := api.NewRuleServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: createRule.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListRules(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	ruleIDs := []string{}
	ruleNames := []string{}
	ruleStatuses := []common.Status{}
	for i := 0; i < 3; i++ {
		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		createRule, err := ruleCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		ruleIDs = append(ruleIDs, createRule.Id)
		ruleNames = append(ruleNames, createRule.Name)
		ruleStatuses = append(ruleStatuses, createRule.Status)
	}

	t.Run("List rules by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		listRules, err := ruleCli.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listRules.Rules), 3)
		require.GreaterOrEqual(t, listRules.TotalSize, int32(3))

		var found bool
		for _, rule := range listRules.Rules {
			if rule.Id == ruleIDs[len(ruleIDs)-1] &&
				rule.Name == ruleNames[len(ruleNames)-1] &&
				rule.Status == ruleStatuses[len(ruleStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List rules by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		listRules, err := ruleCli.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Len(t, listRules.Rules, 2)
		require.NotEmpty(t, listRules.NextPageToken)
		require.GreaterOrEqual(t, listRules.TotalSize, int32(3))

		nextRules, err := ruleCli.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2, PageToken: listRules.NextPageToken})
		t.Logf("nextRules, err: %+v, %v", nextRules, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(nextRules.Rules), 1)
		require.GreaterOrEqual(t, nextRules.TotalSize, int32(3))
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewRuleServiceClient(secondaryAdminGRPCConn)
		listRules, err := secCli.ListRules(ctx, &api.ListRulesRequest{})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Len(t, listRules.Rules, 0)
		require.Equal(t, int32(0), listRules.TotalSize)
	})

	t.Run("List rules by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		listRules, err := ruleCli.ListRules(ctx, &api.ListRulesRequest{
			PageToken: badUUID})
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
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				`true`, true, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				`10 > 15`, false, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				`point.Token == ""`, true, ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40},
				Ts: timestamppb.New(time.Now().Add(-time.Second))},
				`pointTS < currTS`, true, ""},
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
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				`1 + "aaa"`, false, "invalid operation: int + string"},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				`"aaa"`, false, rule.ErrNotBool.Error()},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Can evaluate %+v", lTest), func(t *testing.T) {
				t.Parallel()

				lTest.inpPoint.UniqId = "api-rule-" + random.String(16)
				lTest.inpPoint.Attr = "api-rule" + random.String(10)

				rule := random.Rule("api-rule", uuid.NewString())
				rule.Attr = lTest.inpPoint.Attr
				rule.Expr = lTest.inpRuleExpr

				ctx, cancel := context.WithTimeout(context.Background(),
					testTimeout)
				defer cancel()

				ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
				testRes, err := ruleCli.TestRule(ctx, &api.TestRuleRequest{
					Point: lTest.inpPoint, Rule: rule})
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

	t.Run("Test rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "api-rule-" + random.String(16),
			Attr:     "api-rule" + random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123}}
		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(secondaryViewerGRPCConn)
		testRes, err := ruleCli.TestRule(ctx, &api.TestRuleRequest{Point: point,
			Rule: rule})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Test invalid rule attribute", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "api-rule-" + random.String(16),
			Attr:     "api-rule" + random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123}}
		rule := random.Rule("api-rule", uuid.NewString())

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		ruleCli := api.NewRuleServiceClient(globalAdminGRPCConn)
		testRes, err := ruleCli.TestRule(ctx, &api.TestRuleRequest{
			Point: point, Rule: rule})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"data point and rule attribute mismatch")
	})
}
