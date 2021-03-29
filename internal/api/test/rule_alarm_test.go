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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: rule})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.Nil(t, createRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateRuleRequest.Rule: embedded message failed "+
			"validation | caused by: invalid Rule.Attr: value length must be "+
			"at most 40 runes")
	})
}

func TestCreateAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Create valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		alarm := random.Alarm("api-alarm", uuid.NewString(), createRule.Id)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: alarm})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)
		require.NotEqual(t, alarm.Id, createAlarm.Id)
		require.WithinDuration(t, time.Now(), createAlarm.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createAlarm.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(),
				uuid.NewString())})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create invalid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		alarm.Name = "api-alarm-" + random.String(80)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: alarm})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateAlarmRequest.Alarm: embedded message failed "+
			"validation | caused by: invalid Alarm.Name: value length must be "+
			"between 5 and 80 runes, inclusive")
	})

	t.Run("Create valid alarm with unknown rule", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: alarm})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: alarms_rule_id_fkey")
	})
}

func TestGetRule(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
	createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
		Rule: random.Rule("api-rule", uuid.NewString())})
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Get rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getRule, err := raCli.GetRule(ctx, &api.GetRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getRule, err := raCli.GetRule(ctx, &api.GetRuleRequest{
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

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		getRule, err := secCli.GetRule(ctx, &api.GetRuleRequest{
			Id: createRule.Id})
		t.Logf("getRule, err: %+v, %v", getRule, err)
		require.Nil(t, getRule)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestGetAlarm(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
	createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
		Rule: random.Rule("api-alarm", uuid.NewString())})
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
		Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
	t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
	require.NoError(t, err)

	t.Run("Get alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.Id, RuleId: createRule.Id})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createAlarm, getAlarm) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createAlarm, getAlarm)
		}
	})

	t.Run("Get alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: uuid.NewString(), RuleId: createRule.Id})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Get alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.Id, RuleId: uuid.NewString()})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Get are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		getAlarm, err := secCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.Id, RuleId: createRule.Id})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Name = "api-rule-" + random.String(10)
		createRule.Status = common.Status_DISABLED

		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: createRule})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.NoError(t, err)
		require.Equal(t, createRule.CreatedAt.AsTime(),
			updateRule.CreatedAt.AsTime())
		require.True(t, updateRule.UpdatedAt.AsTime().After(
			updateRule.CreatedAt.AsTime()))
		require.WithinDuration(t, createRule.CreatedAt.AsTime(),
			updateRule.UpdatedAt.AsTime(), 2*time.Second)

		getRule, err := raCli.GetRule(ctx, &api.GetRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		part := &common.Rule{Id: createRule.Id, Name: "api-rule-" +
			random.String(10), Status: common.Status_DISABLED, Expr: `false`}

		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
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

		getRule, err := raCli.GetRule(ctx, &api.GetRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Update nil rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.OrgId = uuid.NewString()
		createRule.Name = "api-rule-" + random.String(10)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		// Update rule fields.
		createRule.Attr = "api-rule-" + random.String(40)

		updateRule, err := raCli.UpdateRule(ctx, &api.UpdateRuleRequest{
			Rule: createRule})
		t.Logf("updateRule, err: %+v, %v", updateRule, err)
		require.Nil(t, updateRule)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateRuleRequest.Rule: embedded message failed "+
			"validation | caused by: invalid Rule.Attr: value length must be "+
			"at most 40 runes")
	})
}

func TestUpdateAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.Name = "api-alarm-" + random.String(10)
		createAlarm.Status = common.Status_DISABLED

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.NoError(t, err)
		require.Equal(t, createAlarm.CreatedAt.AsTime(),
			updateAlarm.CreatedAt.AsTime())
		require.True(t, updateAlarm.UpdatedAt.AsTime().After(
			updateAlarm.CreatedAt.AsTime()))
		require.WithinDuration(t, createAlarm.CreatedAt.AsTime(),
			updateAlarm.UpdatedAt.AsTime(), 2*time.Second)

		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.Id, RuleId: createRule.Id})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateAlarm, getAlarm) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateAlarm, getAlarm)
		}
	})

	t.Run("Partial update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		part := &api.Alarm{Id: createAlarm.Id, RuleId: createRule.Id,
			Name:   "api-alarm-" + random.String(10),
			Status: common.Status_DISABLED}

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status"}}})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.NoError(t, err)
		require.Equal(t, createAlarm.CreatedAt.AsTime(),
			updateAlarm.CreatedAt.AsTime())
		require.True(t, updateAlarm.UpdatedAt.AsTime().After(
			updateAlarm.CreatedAt.AsTime()))
		require.WithinDuration(t, createAlarm.CreatedAt.AsTime(),
			updateAlarm.UpdatedAt.AsTime(), 2*time.Second)

		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.Id, RuleId: createRule.Id})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(updateAlarm, getAlarm) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", updateAlarm, getAlarm)
		}
	})

	t.Run("Update alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Update nil alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: nil})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateAlarmRequest.Alarm: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(),
				uuid.NewString()), UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"}}})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update alarm by unknown alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(),
				createRule.Id), UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status"}}})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Partial update alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		part := &api.Alarm{Id: createAlarm.Id, RuleId: uuid.NewString(),
			Name:   "api-alarm-" + random.String(10),
			Status: common.Status_DISABLED}

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status"}}})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update alarm by unknown alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(),
				createRule.Id)})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.RuleId = uuid.NewString()
		createAlarm.Name = "api-alarm-" + random.String(10)
		createAlarm.Status = common.Status_DISABLED

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.OrgId = uuid.NewString()
		createAlarm.Name = "api-alarm-" + random.String(10)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		updateAlarm, err := secCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update alarm validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.Name = "api-alarm-" + random.String(80)

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateAlarmRequest.Alarm: embedded message failed "+
			"validation | caused by: invalid Alarm.Name: value length must be "+
			"between 5 and 80 runes, inclusive")
	})
}

func TestDeleteRule(t *testing.T) {
	t.Parallel()

	t.Run("Delete rule by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		_, err = raCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: createRule.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read rule by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
			getRule, err := raCli.GetRule(ctx, &api.GetRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		_, err := raCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		_, err := raCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-rule", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteRule(ctx, &api.DeleteRuleRequest{
			Id: createRule.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestDeleteAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Delete alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		_, err = raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: createAlarm.Id, RuleId: createRule.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read alarm by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
			getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
				Id: createAlarm.Id, RuleId: createRule.Id})
			t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
			require.Nil(t, getAlarm)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		_, err := raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: uuid.NewString(), RuleId: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		_, err := raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: uuid.NewString(), RuleId: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Delete alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		_, err = raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: createAlarm.Id, RuleId: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString())})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: createAlarm.Id, RuleId: createRule.Id})
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
		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listRules, err := raCli.ListRules(ctx, &api.ListRulesRequest{})
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listRules, err := raCli.ListRules(ctx, &api.ListRulesRequest{
			PageSize: 2})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.NoError(t, err)
		require.Len(t, listRules.Rules, 2)
		require.NotEmpty(t, listRules.NextPageToken)
		require.GreaterOrEqual(t, listRules.TotalSize, int32(3))

		nextRules, err := raCli.ListRules(ctx, &api.ListRulesRequest{
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

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listRules, err := raCli.ListRules(ctx, &api.ListRulesRequest{
			PageToken: badUUID})
		t.Logf("listRules, err: %+v, %v", listRules, err)
		require.Nil(t, listRules)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}

func TestListAlarms(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
	createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
		Rule: random.Rule("api-alarm", uuid.NewString())})
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	alarmIDs := []string{}
	alarmNames := []string{}
	alarmStatuses := []common.Status{}
	for i := 0; i < 3; i++ {
		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.Id)})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		alarmIDs = append(alarmIDs, createAlarm.Id)
		alarmNames = append(alarmNames, createAlarm.Name)
		alarmStatuses = append(alarmStatuses, createAlarm.Status)
	}

	t.Run("List alarms by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listAlarms.Alarms), 3)
		require.GreaterOrEqual(t, listAlarms.TotalSize, int32(3))

		var found bool
		for _, alarm := range listAlarms.Alarms {
			if alarm.Id == alarmIDs[len(alarmIDs)-1] &&
				alarm.Name == alarmNames[len(alarmNames)-1] &&
				alarm.Status == alarmStatuses[len(alarmStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx, &api.ListAlarmsRequest{
			PageSize: 2})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.Len(t, listAlarms.Alarms, 2)
		require.NotEmpty(t, listAlarms.NextPageToken)
		require.GreaterOrEqual(t, listAlarms.TotalSize, int32(3))

		nextAlarms, err := raCli.ListAlarms(ctx, &api.ListAlarmsRequest{
			PageSize: 2, PageToken: listAlarms.NextPageToken})
		t.Logf("nextAlarms, err: %+v, %v", nextAlarms, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(nextAlarms.Alarms), 1)
		require.GreaterOrEqual(t, nextAlarms.TotalSize, int32(3))
	})

	t.Run("List alarms with rule filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx, &api.ListAlarmsRequest{
			RuleId: createRule.Id})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listAlarms.Alarms), 3)
		require.GreaterOrEqual(t, listAlarms.TotalSize, int32(3))

		var found bool
		for _, alarm := range listAlarms.Alarms {
			if alarm.Id == alarmIDs[len(alarmIDs)-1] &&
				alarm.Name == alarmNames[len(alarmNames)-1] &&
				alarm.Status == alarmStatuses[len(alarmStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		listAlarms, err := secCli.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.Len(t, listAlarms.Alarms, 0)
		require.Equal(t, int32(0), listAlarms.TotalSize)
	})

	t.Run("List alarms by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx, &api.ListAlarmsRequest{
			PageToken: badUUID})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.Nil(t, listAlarms)
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

				raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
				testRes, err := raCli.TestRule(ctx, &api.TestRuleRequest{
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

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		testRes, err := raCli.TestRule(ctx, &api.TestRuleRequest{Point: point,
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

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		testRes, err := raCli.TestRule(ctx, &api.TestRuleRequest{
			Point: point, Rule: rule})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"data point and rule attribute mismatch")
	})
}

func TestTestAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Test valid and invalid alarms", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-alarm", uuid.NewString())
		rule.Name = "test rule"

		dev := random.Device("api-alarm", uuid.NewString())
		dev.Status = common.Status_ACTIVE

		tests := []struct {
			inpPoint *common.DataPoint
			inpRule  *common.Rule
			inpDev   *common.Device
			inpTempl string
			res      string
			err      string
		}{
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				nil, nil, `test`, "test", ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				rule, nil, `point value is an integer: {{.pointVal}}, rule ` +
					`name is: {{.rule.Name}}`, "point value is an integer: " +
					"40, rule name is: test rule", ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_Fl64Val{
				Fl64Val: 37.7}}, nil, dev, `point value is a float: ` +
				`{{.pointVal}}, device status is: {{.device.Status}}`,
				"point value is a float: 37.7, device status is: ACTIVE", ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_StrVal{
				StrVal: "batt"}}, nil, nil, `point value is a string: ` +
				`{{.pointVal}}`, "point value is a string: batt", ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_BoolVal{
				BoolVal: true}}, nil, nil, `point value is a bool: ` +
				`{{.pointVal}}`, "point value is a bool: true", ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_BytesVal{
				BytesVal: []byte{0x00, 0x01}}}, nil, nil, `point value is a ` +
				`byte slice: {{.pointVal}}`, "point value is a byte slice: " +
				"[0 1]", ""},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				nil, nil, `{{if`, "", "unclosed action"},
			{&common.DataPoint{ValOneof: &common.DataPoint_IntVal{IntVal: 40}},
				nil, nil, `{{template "aaa"}}`, "", "no such template"},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Can evaluate %+v", lTest), func(t *testing.T) {
				t.Parallel()

				lTest.inpPoint.UniqId = "api-alarm-" + random.String(16)
				lTest.inpPoint.Attr = "api-alarm" + random.String(10)

				if lTest.inpRule == nil {
					lTest.inpRule = random.Rule("api-alarm", uuid.NewString())
				}

				if lTest.inpDev == nil {
					lTest.inpDev = random.Device("api-alarm", uuid.NewString())
				}

				alarm := random.Alarm("api-alarm", uuid.NewString(),
					uuid.NewString())
				alarm.SubjectTemplate = lTest.inpTempl
				alarm.BodyTemplate = lTest.inpTempl

				ctx, cancel := context.WithTimeout(context.Background(),
					testTimeout)
				defer cancel()

				raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
				testRes, err := raCli.TestAlarm(ctx, &api.TestAlarmRequest{
					Point: lTest.inpPoint, Rule: lTest.inpRule,
					Device: lTest.inpDev, Alarm: alarm})
				t.Logf("testRes, err: %+v, %v", testRes, err)
				if lTest.err == "" {
					require.Equal(t, lTest.res+" - "+lTest.res, testRes.Result)
					require.NoError(t, err)
				} else {
					require.Nil(t, testRes)
					require.Equal(t, codes.InvalidArgument, status.Code(err))
					require.Contains(t, err.Error(), lTest.err)
				}
			})
		}
	})

	t.Run("Test alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "api-alarm-" + random.String(16),
			Attr:     "api-alarm" + random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123}}
		rule := random.Rule("api-alarm", uuid.NewString())
		dev := random.Device("api-alarm", uuid.NewString())
		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		testRes, err := raCli.TestAlarm(ctx, &api.TestAlarmRequest{
			Point: point, Rule: rule, Device: dev, Alarm: alarm})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Test alarm with invalid body template", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{UniqId: "api-alarm-" + random.String(16),
			Attr:     "api-alarm" + random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123}}
		rule := random.Rule("api-alarm", uuid.NewString())
		dev := random.Device("api-alarm", uuid.NewString())
		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		alarm.BodyTemplate = `{{if`

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		testRes, err := raCli.TestAlarm(ctx, &api.TestAlarmRequest{
			Point: point, Rule: rule, Device: dev, Alarm: alarm})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"template: alarm:1: unclosed action")
	})
}
