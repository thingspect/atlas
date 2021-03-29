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
		time.Now().Add(-3*time.Second))
	t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
	require.NoError(t, err)
	require.Len(t, listAlerts, 1)

	alert := &api.Alert{
		OrgId:   createOrg.Id,
		UniqId:  dev.UniqId,
		AlarmId: createAlarm.Id,
		UserId:  createUser.Id,
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
		time.Now().Add(-3*time.Second))
	t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
	require.NoError(t, err)
	require.Len(t, listAlerts, 1)

	alert := &api.Alert{
		OrgId:   createOrg.Id,
		UniqId:  dev.UniqId,
		AlarmId: createAlarm.Id,
		UserId:  createUser.Id,
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
	alarm.SubjectTemplate = `{{if`
	createAlarm, err := globalAlarmDAO.Create(ctx, alarm)
	t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
	require.NoError(t, err)

	dev := random.Device("ale", createOrg.Id)
	dev.Tags = []string{createRule.DeviceTag}

	tests := []struct {
		inp *message.EventerOut
	}{
		// Bad payload.
		{nil},
		// Missing data point.
		{&message.EventerOut{Device: &common.Device{}}},
		// Missing device.
		{&message.EventerOut{Point: &common.DataPoint{}}},
		// Missing rule.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}}},
		// Bad alarm subject.
		{&message.EventerOut{Point: &common.DataPoint{
			TraceId: uuid.NewString()}, Device: dev, Rule: createRule}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Cannot alert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			bEOut := []byte("ale-aaa")
			if lTest.inp != nil {
				// Make device unique per run.
				if lTest.inp.Device != nil {
					lTest.inp.Device.UniqId = "ale" + "-" + random.String(16)
				}

				var err error
				bEOut, err = proto.Marshal(lTest.inp)
				require.NoError(t, err)
				t.Logf("bEOut: %s", bEOut)
			}

			require.NoError(t, globalAleQueue.Publish(globalEOutSubTopic,
				bEOut))
			time.Sleep(2 * time.Second)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			listAlerts, err := globalAleDAO.List(ctx, createOrg.Id, dev.UniqId,
				"", createAlarm.Id, createUser.Id, time.Now(),
				time.Now().Add(-3*time.Second))
			t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
			require.NoError(t, err)
			require.Len(t, listAlerts, 0)
		})
	}
}
