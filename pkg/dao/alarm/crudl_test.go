// +build !unit

package alarm

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
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
		createOrg.Id))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Create valid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("dao-alarm", createOrg.Id, createRule.Id)
		createAlarm, _ := proto.Clone(alarm).(*api.Alarm)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, createAlarm)
		t.Logf("alarm, createAlarm, err: %+v, %+v, %v", alarm, createAlarm, err)
		require.NoError(t, err)
		require.NotEqual(t, alarm.Id, createAlarm.Id)
		require.WithinDuration(t, time.Now(), createAlarm.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createAlarm.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("dao-alarm", createOrg.Id, createRule.Id)
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
			createOrg.Id, uuid.NewString()))
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
		createOrg.Id))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
		createOrg.Id, createRule.Id))
	t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
	require.NoError(t, err)

	t.Run("Read alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.Id,
			createAlarm.OrgId, createRule.Id)
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.NoError(t, err)
		require.Equal(t, createAlarm, readAlarm)
	})

	t.Run("Read alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, uuid.NewString(),
			createAlarm.OrgId, createRule.Id)
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.Nil(t, readAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.Id,
			createAlarm.OrgId, uuid.NewString())
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.Nil(t, readAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.Id,
			uuid.NewString(), createRule.Id)
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.Nil(t, readAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read alarm by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readAlarm, err := globalAlarmDAO.Read(ctx, random.String(10),
			createAlarm.OrgId, createRule.Id)
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
		createOrg.Id))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.Id, createRule.Id))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		// Update alarm fields.
		createAlarm.Name = "dao-alarm-" + random.String(10)
		createAlarm.Status = common.Status_DISABLED
		createAlarm.UserTags = random.Tags("dao-alarm-", 2)
		updateAlarm, _ := proto.Clone(createAlarm).(*api.Alarm)

		updateAlarm, err = globalAlarmDAO.Update(ctx, updateAlarm)
		t.Logf("createAlarm, updateAlarm, err: %+v, %+v, %v", createAlarm,
			updateAlarm, err)
		require.NoError(t, err)
		require.Equal(t, createAlarm.Name, updateAlarm.Name)
		require.True(t, updateAlarm.UpdatedAt.AsTime().After(
			updateAlarm.CreatedAt.AsTime()))
		require.WithinDuration(t, createAlarm.CreatedAt.AsTime(),
			updateAlarm.UpdatedAt.AsTime(), 2*time.Second)

		readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.Id,
			createAlarm.OrgId, createRule.Id)
		t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
		require.NoError(t, err)
		require.Equal(t, updateAlarm, readAlarm)
	})

	t.Run("Update unknown alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateAlarm, err := globalAlarmDAO.Update(ctx, random.Alarm("dao-alarm",
			createOrg.Id, createRule.Id))
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.Id, createRule.Id))
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
			createOrg.Id, createRule.Id))
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
			createOrg.Id, createRule.Id))
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
		createOrg.Id))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	t.Run("Delete alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.Id, createRule.Id))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		err = globalAlarmDAO.Delete(ctx, createAlarm.Id, createOrg.Id,
			createRule.Id)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read alarm by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readAlarm, err := globalAlarmDAO.Read(ctx, createAlarm.Id,
				createOrg.Id, createRule.Id)
			t.Logf("readAlarm, err: %+v, %v", readAlarm, err)
			require.Nil(t, readAlarm)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete alarm by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalAlarmDAO.Delete(ctx, uuid.NewString(), createOrg.Id,
			createRule.Id)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Delete alarm by unknown rule", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.Id, createRule.Id))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		err = globalAlarmDAO.Delete(ctx, createAlarm.Id, createOrg.Id,
			uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.Id, createRule.Id))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		err = globalAlarmDAO.Delete(ctx, createAlarm.Id, uuid.NewString(),
			createRule.Id)
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
		createOrg.Id))
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	alarmIDs := []string{}
	alarmNames := []string{}
	alarmStatuses := []common.Status{}
	alarmTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createAlarm, err := globalAlarmDAO.Create(ctx, random.Alarm("dao-alarm",
			createOrg.Id, createRule.Id))
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.NoError(t, err)

		alarmIDs = append(alarmIDs, createAlarm.Id)
		alarmNames = append(alarmNames, createAlarm.Name)
		alarmStatuses = append(alarmStatuses, createAlarm.Status)
		alarmTSes = append(alarmTSes, createAlarm.CreatedAt.AsTime())
	}

	t.Run("List alarms by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 0, "")
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.Id == alarmIDs[len(alarmIDs)-1] &&
				alarm.Name == alarmNames[len(alarmNames)-1] &&
				alarm.Status == alarmStatuses[len(alarmStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.Id,
			alarmTSes[0], alarmIDs[0], 5, "")
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.Id == alarmIDs[len(alarmIDs)-1] &&
				alarm.Name == alarmNames[len(alarmNames)-1] &&
				alarm.Status == alarmStatuses[len(alarmStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.Id,
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

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 0, createRule.Id)
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.Id == alarmIDs[len(alarmIDs)-1] &&
				alarm.Name == alarmNames[len(alarmNames)-1] &&
				alarm.Status == alarmStatuses[len(alarmStatuses)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List alarms with rule filter and pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlarms, listCount, err := globalAlarmDAO.List(ctx, createOrg.Id,
			alarmTSes[0], alarmIDs[0], 5, createRule.Id)
		t.Logf("listAlarms, listCount, err: %+v, %v, %v", listAlarms, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listAlarms, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, alarm := range listAlarms {
			if alarm.Id == alarmIDs[len(alarmIDs)-1] &&
				alarm.Name == alarmNames[len(alarmNames)-1] &&
				alarm.Status == alarmStatuses[len(alarmStatuses)-1] {
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
