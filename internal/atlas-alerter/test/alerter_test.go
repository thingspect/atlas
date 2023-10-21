//go:build !unit

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
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/atlas-alerter/alerter"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 8 * time.Second

func TestAlertMessages(t *testing.T) {
	t.Parallel()

	for _, alarmType := range []api.AlarmType{
		api.AlarmType_APP,
		api.AlarmType_SMS,
		api.AlarmType_EMAIL,
	} {
		lAlarmType := alarmType

		t.Run(fmt.Sprintf("Can alert %v", lAlarmType), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			createOrg, err := globalOrgDAO.Create(ctx, random.Org("ale"))
			t.Logf("createOrg, err: %+v, %v", createOrg, err)
			require.NoError(t, err)

			rule := random.Rule("ale", createOrg.GetId())
			rule.Status = api.Status_ACTIVE
			createRule, err := globalRuleDAO.Create(ctx, rule)
			t.Logf("createRule, err: %+v, %v", createRule, err)
			require.NoError(t, err)

			user := random.User("dao-user", createOrg.GetId())
			user.Status = api.Status_ACTIVE
			createUser, err := globalUserDAO.Create(ctx, user)
			t.Logf("createUser, err: %+v, %v", createUser, err)
			require.NoError(t, err)

			alarm := random.Alarm("ale", createOrg.GetId(), createRule.GetId())
			alarm.Status = api.Status_ACTIVE
			alarm.Type = lAlarmType
			alarm.UserTags = createUser.GetTags()
			createAlarm, err := globalAlarmDAO.Create(ctx, alarm)
			t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
			require.NoError(t, err)

			dev := random.Device("ale", createOrg.GetId())
			dev.Tags = []string{createRule.GetDeviceTag()}

			eOut := &message.EventerOut{
				Point:  &common.DataPoint{TraceId: uuid.NewString()},
				Device: dev, Rule: createRule,
			}
			bEOut, err := proto.Marshal(eOut)
			require.NoError(t, err)
			t.Logf("bEOut: %s", bEOut)

			require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic,
				bEOut))
			time.Sleep(2 * time.Second)

			ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
			defer cancel()

			listAlerts, err := globalAleDAO.List(ctx, createOrg.GetId(), dev.GetUniqId(),
				"", createAlarm.GetId(), createUser.GetId(), time.Now(),
				time.Now().Add(-4*time.Second))
			t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
			require.NoError(t, err)
			require.Len(t, listAlerts, 1)

			alert := &api.Alert{
				OrgId:   createOrg.GetId(),
				UniqId:  dev.GetUniqId(),
				AlarmId: createAlarm.GetId(),
				UserId:  createUser.GetId(),
				Status:  api.AlertStatus_SENT,
				TraceId: eOut.GetPoint().GetTraceId(),
			}

			// Normalize timestamp.
			require.WithinDuration(t, time.Now(),
				listAlerts[0].GetCreatedAt().AsTime(), testTimeout)
			alert.CreatedAt = listAlerts[0].GetCreatedAt()

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(alert, listAlerts[0]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", alert, listAlerts[0])
			}
		})
	}
}

func TestAlertMessagesRepeat(t *testing.T) {
	t.Parallel()

	for _, alarmType := range []api.AlarmType{
		api.AlarmType_APP,
		api.AlarmType_SMS,
		api.AlarmType_EMAIL,
	} {
		lAlarmType := alarmType

		t.Run(fmt.Sprintf("Can repeat %v", lAlarmType), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			createOrg, err := globalOrgDAO.Create(ctx, random.Org("ale"))
			t.Logf("createOrg, err: %+v, %v", createOrg, err)
			require.NoError(t, err)

			rule := random.Rule("ale", createOrg.GetId())
			rule.Status = api.Status_ACTIVE
			createRule, err := globalRuleDAO.Create(ctx, rule)
			t.Logf("createRule, err: %+v, %v", createRule, err)
			require.NoError(t, err)

			user := random.User("dao-user", createOrg.GetId())
			user.Status = api.Status_ACTIVE
			createUser, err := globalUserDAO.Create(ctx, user)
			t.Logf("createUser, err: %+v, %v", createUser, err)
			require.NoError(t, err)

			alarm := random.Alarm("ale", createOrg.GetId(), createRule.GetId())
			alarm.Status = api.Status_ACTIVE
			alarm.Type = lAlarmType
			alarm.UserTags = createUser.GetTags()
			createAlarm, err := globalAlarmDAO.Create(ctx, alarm)
			t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
			require.NoError(t, err)

			dev := random.Device("ale", createOrg.GetId())
			dev.Tags = []string{rule.GetDeviceTag()}

			eOut := &message.EventerOut{
				Point:  &common.DataPoint{TraceId: uuid.NewString()},
				Device: dev, Rule: createRule,
			}
			bEOut, err := proto.Marshal(eOut)
			require.NoError(t, err)
			t.Logf("bEOut: %s", bEOut)

			// Publish twice. Don't stagger to validate cache locking.
			require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic,
				bEOut))
			require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic,
				bEOut))
			time.Sleep(2 * time.Second)

			ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
			defer cancel()

			listAlerts, err := globalAleDAO.List(ctx, createOrg.GetId(), dev.GetUniqId(),
				"", createAlarm.GetId(), createUser.GetId(), time.Now(),
				time.Now().Add(-4*time.Second))
			t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
			require.NoError(t, err)
			require.Len(t, listAlerts, 1)

			alert := &api.Alert{
				OrgId:   createOrg.GetId(),
				UniqId:  dev.GetUniqId(),
				AlarmId: createAlarm.GetId(),
				UserId:  createUser.GetId(),
				Status:  api.AlertStatus_SENT,
				TraceId: eOut.GetPoint().GetTraceId(),
			}

			// Normalize timestamp.
			require.WithinDuration(t, time.Now(),
				listAlerts[0].GetCreatedAt().AsTime(), testTimeout)
			alert.CreatedAt = listAlerts[0].GetCreatedAt()

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(alert, listAlerts[0]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", alert, listAlerts[0])
			}
		})
	}
}

func TestAlertMessagesError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("ale"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	badSubjRule := random.Rule("ale", createOrg.GetId())
	badSubjRule.Status = api.Status_ACTIVE
	createBadSubjRule, err := globalRuleDAO.Create(ctx, badSubjRule)
	t.Logf("createBadSubjRule, err: %+v, %v", createBadSubjRule, err)
	require.NoError(t, err)

	unspecTypeRule := random.Rule("ale", createOrg.GetId())
	unspecTypeRule.Status = api.Status_ACTIVE
	createUnspecTypeRule, err := globalRuleDAO.Create(ctx, unspecTypeRule)
	t.Logf("createUnspecTypeRule, err: %+v, %v", createUnspecTypeRule, err)
	require.NoError(t, err)

	user := random.User("dao-user", createOrg.GetId())
	user.Status = api.Status_ACTIVE
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	badSubjAlarm := random.Alarm("ale", createOrg.GetId(), createBadSubjRule.GetId())
	badSubjAlarm.Status = api.Status_ACTIVE
	badSubjAlarm.Type = api.AlarmType_APP
	badSubjAlarm.UserTags = createUser.GetTags()
	badSubjAlarm.SubjectTemplate = `{{if`
	createBadSubjAlarm, err := globalAlarmDAO.Create(ctx, badSubjAlarm)
	t.Logf("createBadSubjAlarm, err: %+v, %v", createBadSubjAlarm, err)
	require.NoError(t, err)

	unspecTypeAlarm := random.Alarm("ale", createOrg.GetId(),
		createUnspecTypeRule.GetId())
	unspecTypeAlarm.Status = api.Status_ACTIVE
	unspecTypeAlarm.Type = api.AlarmType_ALARM_TYPE_UNSPECIFIED
	unspecTypeAlarm.UserTags = createUser.GetTags()
	createUnspecTypeAlarm, err := globalAlarmDAO.Create(ctx, unspecTypeAlarm)
	t.Logf("createUnspecTypeAlarm, err: %+v, %v", createUnspecTypeAlarm, err)
	require.NoError(t, err)

	dev := random.Device("ale", createOrg.GetId())
	dev.Tags = []string{
		createBadSubjRule.GetDeviceTag(), createUnspecTypeRule.GetDeviceTag(),
	}

	tests := []struct {
		inpEOut      *message.EventerOut
		inpAlarmID   string
		inpNotifyErr error
	}{
		// Bad payload.
		{nil, uuid.NewString(), nil},
		// Missing data point.
		{&message.EventerOut{Device: &api.Device{}}, uuid.NewString(), nil},
		// Missing device.
		{
			&message.EventerOut{Point: &common.DataPoint{}}, uuid.NewString(),
			nil,
		},
		// Missing rule.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
			}, uuid.NewString(), nil,
		},
		// Unknown org. If this fails due to msg.Requeue(), remove it.
		{
			&message.EventerOut{
				Point:  &common.DataPoint{TraceId: uuid.NewString()},
				Device: random.Device("ale", uuid.NewString()),
				Rule:   random.Rule("ale", createOrg.GetId()),
			}, uuid.NewString(), nil,
		},
		// Bad alarm subject.
		{
			&message.EventerOut{
				Point:  &common.DataPoint{TraceId: uuid.NewString()},
				Device: dev, Rule: createBadSubjRule,
			}, createBadSubjAlarm.GetId(), nil,
		},
		// Unspecified alarm type.
		{
			&message.EventerOut{
				Point:  &common.DataPoint{TraceId: uuid.NewString()},
				Device: dev, Rule: createUnspecTypeRule,
			}, createUnspecTypeAlarm.GetId(), alerter.ErrUnknownAlarm,
		},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot alert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bEOut := []byte("ale-aaa")
			if lTest.inpEOut != nil {
				var err error
				bEOut, err = proto.Marshal(lTest.inpEOut)
				require.NoError(t, err)
				t.Logf("bEOut: %s", bEOut)
			}

			require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic,
				bEOut))
			time.Sleep(2 * time.Second)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			listAlerts, err := globalAleDAO.List(ctx, createOrg.GetId(),
				dev.GetUniqId(), "", lTest.inpAlarmID, createUser.GetId(), time.Now(),
				time.Now().Add(-4*time.Second))
			t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
			require.NoError(t, err)

			if lTest.inpNotifyErr != nil {
				require.Len(t, listAlerts, 1)

				alert := &api.Alert{
					OrgId:   createOrg.GetId(),
					UniqId:  dev.GetUniqId(),
					AlarmId: lTest.inpAlarmID,
					UserId:  createUser.GetId(),
					Status:  api.AlertStatus_ERROR,
					Error:   lTest.inpNotifyErr.Error(),
					TraceId: lTest.inpEOut.GetPoint().GetTraceId(),
				}

				// Normalize timestamp.
				require.WithinDuration(t, time.Now(),
					listAlerts[0].GetCreatedAt().AsTime(), testTimeout)
				alert.CreatedAt = listAlerts[0].GetCreatedAt()

				// Testify does not currently support protobuf equality:
				// https://github.com/stretchr/testify/issues/758
				if !proto.Equal(alert, listAlerts[0]) {
					t.Fatalf("\nExpect: %+v\nActual: %+v", alert, listAlerts[0])
				}
			} else {
				require.Len(t, listAlerts, 0)
			}
		})
	}
}
