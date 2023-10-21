//go:build !unit

package alarm

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 8 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alarm"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-alarm",
		createOrg.GetId()))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Create valid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("dao-alarm", createOrg.GetId(), createRule.GetId())
		createAlarm, _ := proto.Clone(alarm).(*api.Alarm)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, createAlarm)
		t.Logf("alarm, createAlarm, err: %+v, %+v, %v", alarm, createAlarm, err)
		require.NoError(t, err)
		require.NotEqual(t, alarm.GetId(), createAlarm.GetId())
		require.WithinDuration(t, time.Now(), createAlarm.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createAlarm.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("dao-alarm", createOrg.GetId(), createRule.GetId())
		alarm.Name = "dao-alarm-" + random.String(80)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, alarm)
		t.Logf("alarm, createAlarm, err: %+v, %+v, %v", alarm, createAlarm, err)
		require.Nil(t, createAlarm)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})

	t.Run("Create valid alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), uuid.NewString()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alarm"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-alarm",
		createOrg.GetId()))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
		createOrg.GetId(), createRule.GetId()))
	t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
	require.NoError(t, err)

	t.Run("Read alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.GetId(),
			createAlarm.GetOrgId(), createRule.GetId())
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.NoError(t, err)
		require.Equal(t, createAlarm, readAlarm)
	})

	t.Run("Read alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, uuid.NewString(),
			createAlarm.GetOrgId(), createRule.GetId())
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.Nil(t, readAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.GetId(),
			createAlarm.GetOrgId(), uuid.NewString())
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.Nil(t, readAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.GetId(),
			uuid.NewString(), createRule.GetId())
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.Nil(t, readAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read alarm by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, random.String(10),
			createAlarm.GetOrgId(), createRule.GetId())
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.Nil(t, readAlarm)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alarm"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-alarm",
		createOrg.GetId()))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.Name = "dao-alarm-" + random.String(10)
		createAlarm.Status = api.Status_DISABLED
		createAlarm.Type = api.AlarmType_APP
		createAlarm.UserTags = random.Tags("dao-alarm", 2)
		updateAlarm, _ := proto.Clone(createAlarm).(*api.Alarm)

		updateAlarm, err = globalAlarmDAO.Update(ctx, updateAlarm)
		t.Logf("createAlarm, updateAlarm, err: %+v, %+v, %v", createAlarm,
			updateAlarm, err)
		require.NoError(t, err)
		require.Equal(t, createAlarm.GetName(), updateAlarm.GetName())
		require.Equal(t, createAlarm.GetStatus(), updateAlarm.GetStatus())
		require.Equal(t, createAlarm.GetType(), updateAlarm.GetType())
		require.Equal(t, createAlarm.GetUserTags(), updateAlarm.GetUserTags())
		require.True(t, updateAlarm.GetUpdatedAt().AsTime().After(
			updateAlarm.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createAlarm.GetCreatedAt().AsTime(),
			updateAlarm.GetUpdatedAt().AsTime(), 2*time.Second)

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.GetId(),
			createAlarm.GetOrgId(), createRule.GetId())
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.NoError(t, err)
		require.Equal(t, updateAlarm, readAlarm)
	})

	t.Run("Update unknown alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateAlarm, err := globalAlarmDAO.Update(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.RuleId = uuid.NewString()
		updateAlarm, _ := proto.Clone(createAlarm).(*api.Alarm)

		updateAlarm, err = globalAlarmDAO.Update(ctx, updateAlarm)
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.OrgId = uuid.NewString()
		createAlarm.Name = "dao-alarm-" + random.String(10)

		updateAlarm, err := globalAlarmDAO.Update(ctx, createAlarm)
		t.Logf("createAlarm, updateAlarm, err: %+v, %+v, %v", createAlarm,
			updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update alarm by invalid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.Name = "dao-alarm-" + random.String(80)

		updateAlarm, err := globalAlarmDAO.Update(ctx, createAlarm)
		t.Logf("createAlarm, updateAlarm, err: %+v, %+v, %v", createAlarm,
			updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alarm"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-alarm",
		createOrg.GetId()))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Delete alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		err = globalAlarmDAO.Delete(ctx, createAlarm.GetId(), createOrg.GetId(),
			createRule.GetId())
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read alarm by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.GetId(),
				createOrg.GetId(), createRule.GetId())
			t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
			require.Nil(t, readAlarm)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalAlarmDAO.Delete(ctx, uuid.NewString(), createOrg.GetId(),
			createRule.GetId())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Delete alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		err = globalAlarmDAO.Delete(ctx, createAlarm.GetId(), createOrg.GetId(),
			uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		err = globalAlarmDAO.Delete(ctx, createAlarm.GetId(), uuid.NewString(),
			createRule.GetId())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alarm"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createRule, err := globalRuleDAO.Create(ctx, random.Rule("dao-alarm",
		createOrg.GetId()))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	alarmIDs := []string{}
	alarmNames := []string{}
	alarmTypes := []api.AlarmType{}
	alarmTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.GetId(), createRule.GetId()))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		alarmIDs = append(alarmIDs, createAlarm.GetId())
		alarmNames = append(alarmNames, createAlarm.GetName())
		alarmTypes = append(alarmTypes, createAlarm.GetType())
		alarmTSes = append(alarmTSes, createAlarm.GetCreatedAt().AsTime())
	}

	t.Run("List alarms by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.GetId(),
			time.Time{}, "", 0, "")
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.GetId() == alarmIDs[len(alarmIDs)-1] &&
				alarm.GetName() == alarmNames[len(alarmNames)-1] &&
				alarm.GetType() == alarmTypes[len(alarmTypes)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.GetId(),
			alarmTSes[0], alarmIDs[0], 5, "")
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.GetId() == alarmIDs[len(alarmIDs)-1] &&
				alarm.GetName() == alarmNames[len(alarmNames)-1] &&
				alarm.GetType() == alarmTypes[len(alarmTypes)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.GetId(),
			time.Time{}, "", 1, "")
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 1)
		require.Equal(t, int32(3), listCount)
	})

	t.Run("List alarms with rule filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.GetId(),
			time.Time{}, "", 0, createRule.GetId())
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.GetId() == alarmIDs[len(alarmIDs)-1] &&
				alarm.GetName() == alarmNames[len(alarmNames)-1] &&
				alarm.GetType() == alarmTypes[len(alarmTypes)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms with rule filter and pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.GetId(),
			alarmTSes[0], alarmIDs[0], 5, createRule.GetId())
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.GetId() == alarmIDs[len(alarmIDs)-1] &&
				alarm.GetName() == alarmNames[len(alarmNames)-1] &&
				alarm.GetType() == alarmTypes[len(alarmTypes)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, uuid.NewString(),
			time.Time{}, "", 0, uuid.NewString())
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, uuid.NewString(),
			time.Time{}, "", 0, "")
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("List alarms by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, random.String(10),
			time.Time{}, "", 0, "")
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.Nil(t, listAlarms)
		require.Equal(t, int32(0), listCount)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
