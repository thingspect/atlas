//go:build !integration

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
	"github.com/thingspect/atlas/internal/atlas-alerter/notify"
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

	org := random.Org("ale")

	appAlarm := random.Alarm("ale", org.Id, uuid.NewString())
	appAlarm.Status = common.Status_ACTIVE
	appAlarm.Type = api.AlarmType_APP

	smsAlarm := random.Alarm("ale", org.Id, uuid.NewString())
	smsAlarm.Status = common.Status_ACTIVE
	smsAlarm.Type = api.AlarmType_SMS

	emailAlarm := random.Alarm("ale", org.Id, uuid.NewString())
	emailAlarm.Status = common.Status_ACTIVE
	emailAlarm.Type = api.AlarmType_EMAIL

	disAlarm := random.Alarm("ale", org.Id, uuid.NewString())
	disAlarm.Status = common.Status_DISABLED

	tests := []struct {
		inpEOut       *message.EventerOut
		inpAlarms     []*api.Alarm
		inpUsers      []*api.User
		inpUserTimes  int
		inpAppTimes   int
		inpSMSTimes   int
		inpEmailTimes int
		inpAlertTimes int
		inpSeedCache  bool
	}{
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.Id), Rule: random.Rule("ale",
					org.Id),
			},
			[]*api.Alarm{appAlarm},
			[]*api.User{random.User("ale", org.Id)},
			1, 1, 0, 0, 1, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.Id), Rule: random.Rule("ale",
					org.Id),
			},
			[]*api.Alarm{smsAlarm},
			[]*api.User{random.User("ale", org.Id)},
			1, 0, 1, 0, 1, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.Id), Rule: random.Rule("ale",
					org.Id),
			}, []*api.Alarm{emailAlarm}, []*api.User{
				random.User("ale", org.Id),
			}, 1, 0, 0, 1, 1, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.Id), Rule: random.Rule("ale",
					org.Id),
			},
			[]*api.Alarm{disAlarm},
			[]*api.User{random.User("ale", org.Id)},
			0, 0, 0, 0, 0, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.Id), Rule: random.Rule("ale",
					org.Id),
			}, []*api.Alarm{appAlarm}, []*api.User{}, 1, 0, 0, 0, 0, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.Id), Rule: random.Rule("ale",
					org.Id),
			},
			[]*api.Alarm{appAlarm},
			[]*api.User{random.User("ale", org.Id)},
			1, 0, 0, 0, 0, true,
		},
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

			ctrl := gomock.NewController(t)
			orger := NewMockorger(ctrl)
			orger.EXPECT().Read(gomock.Any(), lTest.inpEOut.Device.OrgId).
				Return(org, nil).Times(1)

			alarmer := NewMockalarmer(ctrl)
			alarmer.EXPECT().List(gomock.Any(), lTest.inpEOut.Device.OrgId,
				time.Time{}, "", int32(0), lTest.inpEOut.Rule.Id).Return(
				lTest.inpAlarms, int32(0), nil).Times(1)

			userer := NewMockuserer(ctrl)
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

			notifier := notify.NewMockNotifier(ctrl)
			notifier.EXPECT().App(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(nil).Times(lTest.inpAppTimes)
			notifier.EXPECT().SMS(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(nil).Times(lTest.inpSMSTimes)
			notifier.EXPECT().Email(gomock.Any(), org.DisplayName, org.Email,
				gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).
				Times(lTest.inpEmailTimes)

			alerter := NewMockalerter(ctrl)
			alerter.EXPECT().Create(gomock.Any(),
				matcher.NewProtoMatcher(alert)).
				DoAndReturn(func(ctx interface{}, alert interface{}) error {
					defer wg.Done()

					return nil
				}).Times(lTest.inpAlertTimes)

			cache := cache.NewMemory()

			if lTest.inpSeedCache {
				ctx, cancel := context.WithTimeout(context.Background(),
					5*time.Second)
				defer cancel()

				key := repeatKey(lTest.inpEOut.Device.OrgId,
					lTest.inpEOut.Device.Id, lTest.inpAlarms[0].Id,
					lTest.inpUsers[0].Id)
				ok, err := cache.SetIfNotExistTTL(ctx, key, 1, time.Minute)
				require.True(t, ok)
				require.NoError(t, err)
			}

			ale := Alerter{
				orgDAO:   orger,
				alarmDAO: alarmer,
				userDAO:  userer,
				aleDAO:   alerter,
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

	org := random.Org("ale")

	appAlarm := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	appAlarm.Status = common.Status_ACTIVE
	appAlarm.Type = api.AlarmType_APP

	badSubj := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	badSubj.Status = common.Status_ACTIVE
	badSubj.SubjectTemplate = `{{if`

	badBody := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	badBody.Status = common.Status_ACTIVE
	badBody.BodyTemplate = `{{if`

	unspecType := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	unspecType.Status = common.Status_ACTIVE
	unspecType.Type = api.AlarmType_ALARM_TYPE_UNSPECIFIED

	unknownType := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	unknownType.Status = common.Status_ACTIVE
	unknownType.Type = 999

	tests := []struct {
		inpEOut       *message.EventerOut
		inpOrg        *api.Org
		inpOrgErr     error
		inpAlarms     []*api.Alarm
		inpAlarmsErr  error
		inpUsers      []*api.User
		inpUsersErr   error
		inpCache      bool
		inpCacheErr   error
		inpAppErr     error
		inpAlertErr   error
		inpOrgTimes   int
		inpAlarmTimes int
		inpUserTimes  int
		inpCacheTimes int
		inpAppTimes   int
		inpAlertTimes int
	}{
		// Bad payload.
		{
			nil, nil, nil, nil, nil, nil, nil, true, nil, nil, nil, 0, 0, 0, 0,
			0, 0,
		},
		// Missing data point.
		{
			&message.EventerOut{}, nil, nil, nil, nil, nil, nil, true, nil, nil,
			nil, 0, 0, 0, 0, 0, 0,
		},
		// Missing device.
		{
			&message.EventerOut{Point: &common.DataPoint{}}, nil, nil, nil, nil,
			nil, nil, true, nil, nil, nil, 0, 0, 0, 0, 0, 0,
		},
		// Missing rule.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
			}, nil, nil, nil, nil, nil, nil, true, nil, nil, nil, 0, 0, 0, 0, 0,
			0,
		},
		// Orger error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, nil, errTestProc, nil, nil, nil, nil, true, nil, nil, nil, 1, 0,
			0, 0, 0, 0,
		},
		// Alarmer error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, nil, errTestProc, nil, nil, true, nil, nil, nil, 1, 1,
			0, 0, 0, 0,
		},
		// Userer error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil,
			[]*api.Alarm{appAlarm},
			nil, nil, errTestProc, true,
			nil, nil, nil, 1, 1, 1, 0, 0, 0,
		},
		// Bad alarm subject.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, []*api.Alarm{badSubj}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 0, 0, 0,
		},
		// Bad alarm body.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, []*api.Alarm{badBody}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 0, 0, 0,
		},
		// Cacher error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, []*api.Alarm{appAlarm}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, false, errTestProc, nil, nil, 1, 1, 1, 1, 0, 0,
		},
		// Notifier error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, []*api.Alarm{appAlarm}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, errTestProc, nil, 1, 1, 1, 1, 1, 1,
		},
		// Unspecified alarm type.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, []*api.Alarm{unspecType}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 1, 0, 1,
		},
		// Unknown alarm type.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, []*api.Alarm{unknownType}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 1, 0, 1,
		},
		// Alerter error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &common.Device{},
				Rule: &common.Rule{},
			}, org, nil, []*api.Alarm{appAlarm}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, errTestProc, 1, 1, 1, 1, 1, 1,
		},
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

			ctrl := gomock.NewController(t)
			orger := NewMockorger(ctrl)
			orger.EXPECT().Read(gomock.Any(), gomock.Any()).Return(lTest.inpOrg,
				lTest.inpOrgErr).Times(lTest.inpOrgTimes)

			alarmer := NewMockalarmer(ctrl)
			alarmer.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any(), gomock.Any()).
				Return(lTest.inpAlarms, int32(0), lTest.inpAlarmsErr).
				Times(lTest.inpAlarmTimes)

			userer := NewMockuserer(ctrl)
			userer.EXPECT().ListByTags(gomock.Any(), gomock.Any(),
				gomock.Any()).Return(lTest.inpUsers, lTest.inpUsersErr).
				Times(lTest.inpUserTimes)

			cacher := cache.NewMockCacher(ctrl)
			cacher.EXPECT().SetIfNotExistTTL(gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any()).Return(lTest.inpCache,
				lTest.inpCacheErr).Times(lTest.inpCacheTimes)

			notifier := notify.NewMockNotifier(ctrl)
			notifier.EXPECT().App(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(lTest.inpAppErr).Times(lTest.inpAppTimes)

			alerter := NewMockalerter(ctrl)
			alerter.EXPECT().Create(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx interface{}, alert interface{}) error {
					defer wg.Done()

					return lTest.inpAlertErr
				}).Times(lTest.inpAlertTimes)

			ale := Alerter{
				orgDAO:   orger,
				alarmDAO: alarmer,
				userDAO:  userer,
				aleDAO:   alerter,
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
