// +build !integration

package alerter

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/alerter/notify"
	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const errTestProc consterr.Error = "alerter: test processor error"

func TestAlertMessages(t *testing.T) {
	t.Parallel()

	orgID := uuid.NewString()

	appAlarm := random.Alarm("ale", orgID, uuid.NewString())
	appAlarm.Status = common.Status_ACTIVE
	appAlarm.Type = api.AlarmType_APP

	smsAlarm := random.Alarm("ale", orgID, uuid.NewString())
	smsAlarm.Status = common.Status_ACTIVE
	smsAlarm.Type = api.AlarmType_SMS

	disAlarm := random.Alarm("ale", orgID, uuid.NewString())
	disAlarm.Status = common.Status_DISABLED

	tests := []struct {
		inpEOut       *message.EventerOut
		inpAlarms     []*api.Alarm
		inpUsers      []*api.User
		inpUserTimes  int
		inpAppTimes   int
		inpSMSTimes   int
		inpAlertTimes int
		inpSeedCache  bool
	}{
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: random.Device("ale", orgID),
			Rule:   random.Rule("ale", orgID)}, []*api.Alarm{appAlarm},
			[]*api.User{random.User("ale", orgID)}, 1, 1, 0, 1, false},
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: random.Device("ale", orgID),
			Rule:   random.Rule("ale", orgID)}, []*api.Alarm{smsAlarm},
			[]*api.User{random.User("ale", orgID)}, 1, 0, 1, 1, false},
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: random.Device("ale", orgID),
			Rule:   random.Rule("ale", orgID)}, []*api.Alarm{disAlarm},
			[]*api.User{random.User("ale", orgID)}, 0, 0, 0, 0, false},
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: random.Device("ale", orgID),
			Rule:   random.Rule("ale", orgID)}, []*api.Alarm{appAlarm},
			[]*api.User{}, 1, 0, 0, 0, false},
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: random.Device("ale", orgID),
			Rule:   random.Rule("ale", orgID)}, []*api.Alarm{appAlarm},
			[]*api.User{random.User("ale", orgID)}, 1, 0, 0, 0, true},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can alert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			eOutQueue := queue.NewFake()
			eOutSub, err := eOutQueue.Subscribe("")
			require.NoError(t, err)

			var wg sync.WaitGroup
			wg.Add(1)

			alarmer := NewMockalarmer(gomock.NewController(t))
			alarmer.EXPECT().List(gomock.Any(), lTest.inpEOut.Device.OrgId,
				time.Time{}, "", int32(0), lTest.inpEOut.Rule.Id).Return(
				lTest.inpAlarms, int32(0), nil).Times(1)

			userer := NewMockuserer(gomock.NewController(t))
			userer.EXPECT().ListByTags(gomock.Any(), lTest.inpEOut.Device.OrgId,
				lTest.inpAlarms[0].UserTags).Return(lTest.inpUsers, nil).
				Times(lTest.inpUserTimes)

			var alert *api.Alert
			if lTest.inpAlertTimes > 0 {
				alert = &api.Alert{
					OrgId:   lTest.inpEOut.Device.OrgId,
					UniqId:  lTest.inpEOut.Device.UniqId,
					AlarmId: lTest.inpAlarms[0].Id,
					UserId:  lTest.inpUsers[0].Id,
					Status:  api.AlertStatus_SENT,
					TraceId: lTest.inpEOut.Point.TraceId,
				}
			}

			notifier := notify.NewMockNotifier(gomock.NewController(t))
			notifier.EXPECT().App(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(nil).Times(lTest.inpAppTimes)
			notifier.EXPECT().SMS(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(nil).Times(lTest.inpSMSTimes)

			alerter := NewMockalerter(gomock.NewController(t))
			alerter.EXPECT().Create(gomock.Any(), matcher.NewProtoMatcher(
				alert)).DoAndReturn(func(_ ...interface{}) error {
				defer wg.Done()

				return nil
			}).Times(lTest.inpAlertTimes)

			cache, err := cache.NewMemory()
			require.NoError(t, err)

			if lTest.inpSeedCache {
				ctx, cancel := context.WithTimeout(context.Background(),
					5*time.Second)
				defer cancel()

				key := Key(lTest.inpEOut.Device.OrgId, lTest.inpEOut.Device.Id,
					lTest.inpAlarms[0].Id, lTest.inpUsers[0].Id)
				ok, err := cache.SetIfNotExistTTL(ctx, key, 0, time.Minute)
				require.True(t, ok)
				require.NoError(t, err)
			}

			ale := Alerter{
				alarmDAO: alarmer,
				userDAO:  userer,
				alertDAO: alerter,
				cache:    cache,

				aleQueue: eOutQueue,
				eOutSub:  eOutSub,

				notify: notifier,
			}
			go func() {
				ale.alertMessages()
			}()

			bEOut, err := proto.Marshal(lTest.inpEOut)
			require.NoError(t, err)
			t.Logf("bEOut: %s", bEOut)

			require.NoError(t, eOutQueue.Publish("", bEOut))
			if lTest.inpAlertTimes > 0 {
				wg.Wait()
			} else {
				// If the success mode isn't supported by WaitGroup operation,
				// give it time to traverse the code.
				time.Sleep(100 * time.Millisecond)
			}
		})
	}
}

func TestAlertMessagesError(t *testing.T) {
	t.Parallel()

	appAlarm := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	appAlarm.Status = common.Status_ACTIVE
	appAlarm.Type = api.AlarmType_APP

	badSubj := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	badSubj.Status = common.Status_ACTIVE
	badSubj.SubjectTemplate = `{{if`

	badBody := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	badBody.Status = common.Status_ACTIVE
	badBody.BodyTemplate = `{{if`

	badType := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	badType.Status = common.Status_ACTIVE
	badType.Type = 999

	tests := []struct {
		inpEOut       *message.EventerOut
		inpAlarms     []*api.Alarm
		inpAlarmsErr  error
		inpUsers      []*api.User
		inpUsersErr   error
		inpCache      bool
		inpCacheErr   error
		inpAppErr     error
		inpAlertErr   error
		inpAlarmTimes int
		inpUserTimes  int
		inpCacheTimes int
		inpAppTimes   int
		inpAlertTimes int
	}{
		// Bad payload.
		{nil, nil, nil, nil, nil, true, nil, nil, nil, 0, 0, 0, 0, 0},
		// Missing data point.
		{&message.EventerOut{}, nil, nil, nil, nil, true, nil, nil, nil, 0, 0,
			0, 0, 0},
		// Missing device.
		{&message.EventerOut{Point: &common.DataPoint{}}, nil, nil, nil, nil,
			true, nil, nil, nil, 0, 0, 0, 0, 0},
		// Missing rule.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}}, nil, nil, nil, nil, true, nil, nil, nil,
			0, 0, 0, 0, 0},
		// Alarmer error.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}}, nil, errTestProc,
			nil, nil, true, nil, nil, nil, 1, 0, 0, 0, 0},
		// Userer error.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}},
			[]*api.Alarm{appAlarm}, nil, nil, errTestProc, true, nil, nil, nil, 1,
			1, 0, 0, 0},
		// Bad alarm subject.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}},
			[]*api.Alarm{badSubj}, nil, []*api.User{random.User("ale",
				uuid.NewString())}, nil, true, nil, nil, nil, 1, 1, 0, 0, 0},
		// Bad alarm body.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}},
			[]*api.Alarm{badBody}, nil, []*api.User{random.User("ale",
				uuid.NewString())}, nil, true, nil, nil, nil, 1, 1, 0, 0, 0},
		// Cacher error.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}},
			[]*api.Alarm{appAlarm}, nil, []*api.User{random.User("ale",
				uuid.NewString())}, nil, false, errTestProc, nil, nil, 1, 1, 1,
			0, 0},
		// Notifier error.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}},
			[]*api.Alarm{appAlarm}, nil, []*api.User{random.User("ale",
				uuid.NewString())}, nil, true, nil, errTestProc, nil, 1, 1, 1,
			1, 1},
		// Bad alarm type.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}},
			[]*api.Alarm{badType}, nil, []*api.User{random.User("ale",
				uuid.NewString())}, nil, true, nil, nil, nil, 1, 1, 1, 0, 1},
		// Alerter error.
		{&message.EventerOut{Point: &common.DataPoint{},
			Device: &common.Device{}, Rule: &common.Rule{}},
			[]*api.Alarm{appAlarm}, nil, []*api.User{random.User("ale",
				uuid.NewString())}, nil, true, nil, nil, errTestProc, 1, 1, 1,
			1, 1},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can alert %+v", lTest), func(t *testing.T) {
			t.Parallel()

			eOutQueue := queue.NewFake()
			eOutSub, err := eOutQueue.Subscribe("")
			require.NoError(t, err)

			var wg sync.WaitGroup
			wg.Add(1)

			alarmer := NewMockalarmer(gomock.NewController(t))
			alarmer.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any(), gomock.Any()).
				Return(lTest.inpAlarms, int32(0), lTest.inpAlarmsErr).
				Times(lTest.inpAlarmTimes)

			userer := NewMockuserer(gomock.NewController(t))
			userer.EXPECT().ListByTags(gomock.Any(), gomock.Any(),
				gomock.Any()).Return(lTest.inpUsers, lTest.inpUsersErr).
				Times(lTest.inpUserTimes)

			cacher := cache.NewMockCacher(gomock.NewController(t))
			cacher.EXPECT().SetIfNotExistTTL(gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any()).Return(lTest.inpCache,
				lTest.inpCacheErr).Times(lTest.inpCacheTimes)

			notifier := notify.NewMockNotifier(gomock.NewController(t))
			notifier.EXPECT().App(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(lTest.inpAppErr).Times(lTest.inpAppTimes)

			alerter := NewMockalerter(gomock.NewController(t))
			alerter.EXPECT().Create(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ ...interface{}) error {
					defer wg.Done()

					return lTest.inpAlertErr
				}).Times(lTest.inpAlertTimes)

			ale := Alerter{
				alarmDAO: alarmer,
				userDAO:  userer,
				alertDAO: alerter,
				cache:    cacher,

				aleQueue: eOutQueue,
				eOutSub:  eOutSub,

				notify: notifier,
			}
			go func() {
				ale.alertMessages()
			}()

			bEOut := []byte("ale-aaa")
			if lTest.inpEOut != nil {
				bEOut, err = proto.Marshal(lTest.inpEOut)
				require.NoError(t, err)
				t.Logf("bEOut: %s", bEOut)
			}

			require.NoError(t, eOutQueue.Publish("", bEOut))
			if lTest.inpAlertTimes > 0 {
				wg.Wait()
			} else {
				// If the success mode isn't supported by WaitGroup operation,
				// give it time to traverse the code.
				time.Sleep(100 * time.Millisecond)
			}
		})
	}
}
