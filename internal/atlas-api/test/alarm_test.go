//go:build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Create valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		alarm := random.Alarm("api-alarm", uuid.NewString(), createRule.GetId())

		createAlarm, err := raCli.CreateAlarm(ctx,
			&api.CreateAlarmRequest{Alarm: alarm})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)
		require.NotEqual(t, alarm.GetId(), createAlarm.GetId())
		require.WithinDuration(t, time.Now(), createAlarm.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createAlarm.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(),
				uuid.NewString()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Create invalid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		alarm.Name = "api-alarm-" + random.String(80)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createAlarm, err := raCli.CreateAlarm(ctx,
			&api.CreateAlarmRequest{Alarm: alarm})
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createAlarm, err := raCli.CreateAlarm(ctx,
			&api.CreateAlarmRequest{Alarm: alarm})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: alarms_rule_id_fkey")
	})
}

func TestGetAlarm(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
	createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
		Rule: random.Rule("api-alarm", uuid.NewString()),
	})
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
		Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
	})
	t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
	require.NoError(t, err)

	t.Run("Get alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.GetId(), RuleId: createRule.GetId(),
		})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, createAlarm, getAlarm)
	})

	t.Run("Get alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: uuid.NewString(), RuleId: createRule.GetId(),
		})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Get alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.GetId(), RuleId: uuid.NewString(),
		})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Gets are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		getAlarm, err := secCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.GetId(), RuleId: createRule.GetId(),
		})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestUpdateAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.Name = "api-alarm-" + random.String(10)
		createAlarm.Status = api.Status_DISABLED
		createAlarm.Type = api.AlarmType_SMS
		createAlarm.UserTags = random.Tags("api-alarm", 2)

		updateAlarm, err := raCli.UpdateAlarm(ctx,
			&api.UpdateAlarmRequest{Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.NoError(t, err)
		require.Equal(t, createAlarm.GetName(), updateAlarm.GetName())
		require.Equal(t, createAlarm.GetStatus(), updateAlarm.GetStatus())
		require.Equal(t, createAlarm.GetType(), updateAlarm.GetType())
		require.Equal(t, createAlarm.GetUserTags(), updateAlarm.GetUserTags())
		require.True(t, updateAlarm.GetUpdatedAt().AsTime().After(
			updateAlarm.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createAlarm.GetCreatedAt().AsTime(),
			updateAlarm.GetUpdatedAt().AsTime(), 2*time.Second)

		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.GetId(), RuleId: createRule.GetId(),
		})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateAlarm, getAlarm)
	})

	t.Run("Partial update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminKeyGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		part := &api.Alarm{
			Id: createAlarm.GetId(), RuleId: createRule.GetId(), Name: "api-alarm-" +
				random.String(10), Status: api.Status_DISABLED,
			Type: api.AlarmType_SMS, UserTags: random.Tags("api-alarm", 2),
		}

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status", "type", "user_tags"},
			},
		})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.NoError(t, err)
		require.Equal(t, part.GetName(), updateAlarm.GetName())
		require.Equal(t, part.GetStatus(), updateAlarm.GetStatus())
		require.Equal(t, part.GetType(), updateAlarm.GetType())
		require.Equal(t, part.GetUserTags(), updateAlarm.GetUserTags())
		require.True(t, updateAlarm.GetUpdatedAt().AsTime().After(
			updateAlarm.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createAlarm.GetCreatedAt().AsTime(),
			updateAlarm.GetUpdatedAt().AsTime(), 2*time.Second)

		getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: createAlarm.GetId(), RuleId: createRule.GetId(),
		})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateAlarm, getAlarm)
	})

	t.Run("Update alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
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

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateAlarm, err := raCli.UpdateAlarm(ctx,
			&api.UpdateAlarmRequest{Alarm: nil})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateAlarmRequest.Alarm: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(),
				uuid.NewString()), UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"},
			},
		})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update alarm by unknown alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(),
				createRule.GetId()), UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status"},
			},
		})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Partial update alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		part := &api.Alarm{
			Id: createAlarm.GetId(), RuleId: uuid.NewString(), Name: "api-alarm-" +
				random.String(10), Status: api.Status_DISABLED,
		}

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name", "status"},
			},
		})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update alarm by unknown alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		updateAlarm, err := raCli.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.RuleId = uuid.NewString()
		createAlarm.Name = "api-alarm-" + random.String(10)
		createAlarm.Status = api.Status_DISABLED

		updateAlarm, err := raCli.UpdateAlarm(ctx,
			&api.UpdateAlarmRequest{Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.OrgId = uuid.NewString()
		createAlarm.Name = "api-alarm-" + random.String(10)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		updateAlarm, err := secCli.UpdateAlarm(ctx,
			&api.UpdateAlarmRequest{Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update alarm validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.Name = "api-alarm-" + random.String(80)

		updateAlarm, err := raCli.UpdateAlarm(ctx,
			&api.UpdateAlarmRequest{Alarm: createAlarm})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateAlarmRequest.Alarm: embedded message failed "+
			"validation | caused by: invalid Alarm.Name: value length must be "+
			"between 5 and 80 runes, inclusive")
	})
}

func TestDeleteAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Delete alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		_, err = raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: createAlarm.GetId(), RuleId: createRule.GetId(),
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read alarm by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()

			raCli := api.NewRuleAlarmServiceClient(globalAdminKeyGRPCConn)
			getAlarm, err := raCli.GetAlarm(ctx, &api.GetAlarmRequest{
				Id: createAlarm.GetId(), RuleId: createRule.GetId(),
			})
			t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
			require.Nil(t, getAlarm)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		_, err := raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: uuid.NewString(), RuleId: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Delete alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		_, err := raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: uuid.NewString(), RuleId: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Delete alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		_, err = raCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: createAlarm.GetId(), RuleId: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
			Rule: random.Rule("api-alarm", uuid.NewString()),
		})
		t.Logf("createRule, err: %+v, %v", createRule, err)
		require.NoError(t, err)

		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: createAlarm.GetId(), RuleId: createRule.GetId(),
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListAlarms(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
	createRule, err := raCli.CreateRule(ctx, &api.CreateRuleRequest{
		Rule: random.Rule("api-alarm", uuid.NewString()),
	})
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	alarmIDs := []string{}
	alarmNames := []string{}
	alarmStatuses := []api.Status{}
	for range 3 {
		createAlarm, err := raCli.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: random.Alarm("api-alarm", uuid.NewString(), createRule.GetId()),
		})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		alarmIDs = append(alarmIDs, createAlarm.GetId())
		alarmNames = append(alarmNames, createAlarm.GetName())
		alarmStatuses = append(alarmStatuses, createAlarm.GetStatus())
	}

	t.Run("List alarms by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listAlarms.GetAlarms()), 3)
		require.GreaterOrEqual(t, listAlarms.GetTotalSize(), int32(3))

		var found bool
		for _, alarm := range listAlarms.GetAlarms() {
			if alarm.GetId() == alarmIDs[len(alarmIDs)-1] &&
				alarm.GetName() == alarmNames[len(alarmNames)-1] &&
				alarm.GetStatus() == alarmStatuses[len(alarmStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx,
			&api.ListAlarmsRequest{PageSize: 2})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.Len(t, listAlarms.GetAlarms(), 2)
		require.NotEmpty(t, listAlarms.GetNextPageToken())
		require.GreaterOrEqual(t, listAlarms.GetTotalSize(), int32(3))

		nextAlarms, err := raCli.ListAlarms(ctx, &api.ListAlarmsRequest{
			PageSize: 2, PageToken: listAlarms.GetNextPageToken(),
		})
		t.Logf("nextAlarms, err: %+v, %v", nextAlarms, err)
		require.NoError(t, err)
		require.NotEmpty(t, nextAlarms.GetAlarms())
		require.GreaterOrEqual(t, nextAlarms.GetTotalSize(), int32(3))
	})

	t.Run("List alarms with rule filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx,
			&api.ListAlarmsRequest{RuleId: createRule.GetId()})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listAlarms.GetAlarms()), 3)
		require.GreaterOrEqual(t, listAlarms.GetTotalSize(), int32(3))

		var found bool
		for _, alarm := range listAlarms.GetAlarms() {
			if alarm.GetId() == alarmIDs[len(alarmIDs)-1] &&
				alarm.GetName() == alarmNames[len(alarmNames)-1] &&
				alarm.GetStatus() == alarmStatuses[len(alarmStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewRuleAlarmServiceClient(secondaryAdminGRPCConn)
		listAlarms, err := secCli.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.Empty(t, listAlarms.GetAlarms())
		require.Equal(t, int32(0), listAlarms.GetTotalSize())
	})

	t.Run("List alarms by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		listAlarms, err := raCli.ListAlarms(ctx,
			&api.ListAlarmsRequest{PageToken: badUUID})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.Nil(t, listAlarms)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}

func TestTestAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Test valid and invalid alarms", func(t *testing.T) {
		t.Parallel()

		rule := random.Rule("api-alarm", uuid.NewString())
		rule.Name = "test rule"

		dev := random.Device("api-alarm", uuid.NewString())
		dev.Status = api.Status_ACTIVE

		tests := []struct {
			inpPoint *common.DataPoint
			inpRule  *api.Rule
			inpDev   *api.Device
			inpTempl string
			res      string
			err      string
		}{
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, nil, nil, `test`, "test", "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, rule, nil, `point value is an integer: {{.pointVal}}, ` +
					`rule name is: {{.rule.Name}}`, "point value is an " +
					"integer: 40, rule name is: test rule", "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_Fl64Val{
					Fl64Val: 37.7,
				}}, nil, dev, `point value is a float: {{.pointVal}}, device ` +
					`status is: {{.device.Status}}`, "point value is a " +
					"float: 37.7, device status is: ACTIVE", "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_StrVal{
					StrVal: "line",
				}}, nil, nil, `point value is a string: {{.pointVal}}`,
				"point value is a string: line", "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_BoolVal{
					BoolVal: true,
				}}, nil, nil, `point value is a bool: {{.pointVal}}`,
				"point value is a bool: true", "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_BytesVal{
					BytesVal: []byte{0x00, 0x01},
				}}, nil, nil, `point value is a byte slice: {{.pointVal}}`,
				"point value is a byte slice: [0 1]", "",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, nil, nil, `{{if`, "", "unclosed action",
			},
			{
				&common.DataPoint{ValOneof: &common.DataPoint_IntVal{
					IntVal: 40,
				}}, nil, nil, `{{template "aaa"}}`, "", "no such template",
			},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Can evaluate %+v", test), func(t *testing.T) {
				t.Parallel()

				test.inpPoint.UniqId = "api-alarm-" + random.String(16)
				test.inpPoint.Attr = "api-alarm" + random.String(10)

				if test.inpRule == nil {
					test.inpRule = random.Rule("api-alarm", uuid.NewString())
				}

				if test.inpDev == nil {
					test.inpDev = random.Device("api-alarm", uuid.NewString())
				}

				alarm := random.Alarm("api-alarm", uuid.NewString(),
					uuid.NewString())
				alarm.SubjectTemplate = test.inpTempl
				alarm.BodyTemplate = test.inpTempl

				ctx, cancel := context.WithTimeout(t.Context(),
					testTimeout)
				defer cancel()

				raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
				testRes, err := raCli.TestAlarm(ctx, &api.TestAlarmRequest{
					Point: test.inpPoint, Rule: test.inpRule,
					Device: test.inpDev, Alarm: alarm,
				})
				t.Logf("testRes, err: %+v, %v", testRes, err)
				if test.err == "" {
					require.Equal(t, test.res+" - "+test.res, testRes.GetResult())
					require.NoError(t, err)
				} else {
					require.Nil(t, testRes)
					require.Equal(t, codes.InvalidArgument, status.Code(err))
					require.Contains(t, err.Error(), test.err)
				}
			})
		}
	})

	t.Run("Test alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-alarm-" + random.String(16), Attr: "api-alarm" +
				random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
		}
		rule := random.Rule("api-alarm", uuid.NewString())
		dev := random.Device("api-alarm", uuid.NewString())
		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(secondaryViewerGRPCConn)
		testRes, err := raCli.TestAlarm(ctx, &api.TestAlarmRequest{
			Point: point, Rule: rule, Device: dev, Alarm: alarm,
		})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, BUILDER role required")
	})

	t.Run("Test alarm with invalid body template", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-alarm-" + random.String(16), Attr: "api-alarm" +
				random.String(10),
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
		}
		rule := random.Rule("api-alarm", uuid.NewString())
		dev := random.Device("api-alarm", uuid.NewString())
		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		alarm.BodyTemplate = `{{if`

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		raCli := api.NewRuleAlarmServiceClient(globalAdminGRPCConn)
		testRes, err := raCli.TestAlarm(ctx, &api.TestAlarmRequest{
			Point: point, Rule: rule, Device: dev, Alarm: alarm,
		})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"template: template:1: unclosed action")
	})
}
