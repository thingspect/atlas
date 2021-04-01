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
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/alerter/alerter"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 8 * time.Second

func TestAlertMessages(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("ale"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	rule := random.Rule("ale", createOrg.Id)
	rule.Status = common.Status_ACTIVE
	createRule, err := globalRuleDAO.Create(ctx, rule)
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	user := random.User("dao-user", createOrg.Id)
	user.Status = common.Status_ACTIVE
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	alarm := random.Alarm("ale", createOrg.Id, createRule.Id)
	alarm.Status = common.Status_ACTIVE
	alarm.UserTags = createUser.Tags
	createAlarm, err := globalAlarmDAO.Create(ctx, alarm)
	t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
	require.NoError(t, err)

	dev := random.Device("ale", createOrg.Id)
	dev.Tags = []string{createRule.DeviceTag}

	eOut := &message.EventerOut{Point: &common.DataPoint{
		TraceId: uuid.NewString()}, Device: dev, Rule: createRule}
	bEOut, err := proto.Marshal(eOut)
	require.NoError(t, err)
	t.Logf("bEOut: %s", bEOut)

	require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic, bEOut))
	time.Sleep(2 * time.Second)

	ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	listAlerts, err := globalAleDAO.List(ctx, createOrg.Id, dev.UniqId, "",
		createAlarm.Id, createUser.Id, time.Now(),
		time.Now().Add(-4*time.Second))
	t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
	require.NoError(t, err)
	require.Len(t, listAlerts, 1)

	alert := &api.Alert{
		OrgId:   createOrg.Id,
		UniqId:  dev.UniqId,
		AlarmId: createAlarm.Id,
		UserId:  createUser.Id,
		Status:  api.AlertStatus_SENT,
		TraceId: eOut.Point.TraceId,
	}

	// Normalize timestamp.
	require.WithinDuration(t, time.Now(), listAlerts[0].CreatedAt.AsTime(),
		testTimeout)
	alert.CreatedAt = listAlerts[0].CreatedAt

	// Testify does not currently support protobuf equality:
	// https://github.com/stretchr/testify/issues/758
	if !proto.Equal(alert, listAlerts[0]) {
		t.Fatalf("\nExpect: %+v\nActual: %+v", alert, listAlerts[0])
	}
}

func TestAlertMessagesRepeat(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("ale"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	rule := random.Rule("ale", createOrg.Id)
	rule.Status = common.Status_ACTIVE
	createRule, err := globalRuleDAO.Create(ctx, rule)
	t.Logf("createRule, err: %+v, %v", createRule, err)
	require.NoError(t, err)

	user := random.User("dao-user", createOrg.Id)
	user.Status = common.Status_ACTIVE
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	alarm := random.Alarm("ale", createOrg.Id, createRule.Id)
	alarm.Status = common.Status_ACTIVE
	alarm.UserTags = createUser.Tags
	createAlarm, err := globalAlarmDAO.Create(ctx, alarm)
	t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
	require.NoError(t, err)

	dev := random.Device("ale", createOrg.Id)
	dev.Tags = []string{rule.DeviceTag}

	eOut := &message.EventerOut{Point: &common.DataPoint{
		TraceId: uuid.NewString()}, Device: dev, Rule: createRule}
	bEOut, err := proto.Marshal(eOut)
	require.NoError(t, err)
	t.Logf("bEOut: %s", bEOut)

	// Publish twice. Don't stagger to validate cache locking.
	require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic, bEOut))
	require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic, bEOut))
	time.Sleep(2 * time.Second)

	ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	listAlerts, err := globalAleDAO.List(ctx, createOrg.Id, dev.UniqId, "",
		createAlarm.Id, createUser.Id, time.Now(),
		time.Now().Add(-4*time.Second))
	t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
	require.NoError(t, err)
	require.Len(t, listAlerts, 1)

	alert := &api.Alert{
		OrgId:   createOrg.Id,
		UniqId:  dev.UniqId,
		AlarmId: createAlarm.Id,
		UserId:  createUser.Id,
		Status:  api.AlertStatus_SENT,
		TraceId: eOut.Point.TraceId,
	}

	// Normalize timestamp.
	require.WithinDuration(t, time.Now(), listAlerts[0].CreatedAt.AsTime(),
		testTimeout)
	alert.CreatedAt = listAlerts[0].CreatedAt

	// Testify does not currently support protobuf equality:
	// https://github.com/stretchr/testify/issues/758
	if !proto.Equal(alert, listAlerts[0]) {
		t.Fatalf("\nExpect: %+v\nActual: %+v", alert, listAlerts[0])
	}
}

func TestAlertMessagesError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("ale"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	badSubjRule := random.Rule("ale", createOrg.Id)
	badSubjRule.Status = common.Status_ACTIVE
	createBadSubjRule, err := globalRuleDAO.Create(ctx, badSubjRule)
	t.Logf("createBadSubjRule, err: %+v, %v", createBadSubjRule, err)
	require.NoError(t, err)

	badTypeRule := random.Rule("ale", createOrg.Id)
	badTypeRule.Status = common.Status_ACTIVE
	createBadTypeRule, err := globalRuleDAO.Create(ctx, badTypeRule)
	t.Logf("createBadTypeRule, err: %+v, %v", createBadTypeRule, err)
	require.NoError(t, err)

	user := random.User("dao-user", createOrg.Id)
	user.Status = common.Status_ACTIVE
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	badSubjAlarm := random.Alarm("ale", createOrg.Id, createBadSubjRule.Id)
	badSubjAlarm.Status = common.Status_ACTIVE
	badSubjAlarm.UserTags = createUser.Tags
	badSubjAlarm.SubjectTemplate = `{{if`
	createBadSubjAlarm, err := globalAlarmDAO.Create(ctx, badSubjAlarm)
	t.Logf("createBadSubjAlarm, err: %+v, %v", createBadSubjAlarm, err)
	require.NoError(t, err)

	badTypeAlarm := random.Alarm("ale", createOrg.Id, createBadTypeRule.Id)
	badTypeAlarm.Status = common.Status_ACTIVE
	badTypeAlarm.Type = api.AlarmType_ALARM_TYPE_UNSPECIFIED
	badTypeAlarm.UserTags = createUser.Tags
	createBadTypeAlarm, err := globalAlarmDAO.Create(ctx, badTypeAlarm)
	t.Logf("createBadTypeAlarm, err: %+v, %v", createBadTypeAlarm, err)
	require.NoError(t, err)

	dev := random.Device("ale", createOrg.Id)
	dev.Tags = []string{createBadSubjRule.DeviceTag,
		createBadTypeRule.DeviceTag}

	tests := []struct {
		inpEOut      *message.EventerOut
		inpAlarmID   string
		inpNotifyErr error
	}{
		// Bad payload.
		{nil, uuid.NewString(), nil},
		// Missing data point.
		{&message.EventerOut{Device: &common.Device{}}, uuid.NewString(), nil},
		// Missing device.
		{&message.EventerOut{Point: &common.DataPoint{}}, uuid.NewString(),
			nil},
		// Missing rule.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}}, uuid.NewString(), nil},
		// Bad alarm subject.
		{&message.EventerOut{Point: &common.DataPoint{
			TraceId: uuid.NewString()}, Device: dev, Rule: createBadSubjRule},
			createBadSubjAlarm.Id, nil},
		// Bad alarm type.
		{&message.EventerOut{Point: &common.DataPoint{
			TraceId: uuid.NewString()}, Device: dev, Rule: createBadTypeRule},
			createBadTypeAlarm.Id, alerter.ErrUnknownAlarm},
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

			listAlerts, err := globalAleDAO.List(ctx, createOrg.Id,
				dev.UniqId, "", lTest.inpAlarmID, createUser.Id, time.Now(),
				time.Now().Add(-4*time.Second))
			t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
			require.NoError(t, err)

			if lTest.inpNotifyErr != nil {
				require.Len(t, listAlerts, 1)

				alert := &api.Alert{
					OrgId:   createOrg.Id,
					UniqId:  dev.UniqId,
					AlarmId: lTest.inpAlarmID,
					UserId:  createUser.Id,
					Status:  api.AlertStatus_ERROR,
					Error:   lTest.inpNotifyErr.Error(),
					TraceId: lTest.inpEOut.Point.TraceId,
				}

				// Normalize timestamp.
				require.WithinDuration(t, time.Now(),
					listAlerts[0].CreatedAt.AsTime(), testTimeout)
				alert.CreatedAt = listAlerts[0].CreatedAt

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
