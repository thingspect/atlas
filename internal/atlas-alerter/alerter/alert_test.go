//go:build !integration

package alerter

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/pkg/notify"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/message"
	"github.com/thingspect/proto/go/api"
	"github.com/thingspect/proto/go/common"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/proto"
)

const errTestProc consterr.Error = "alerter: test processor error"

func TestAlertMessages(t *testing.T) {
	t.Parallel()

	org := random.Org("ale")

	appAlarm := random.Alarm("ale", org.GetId(), uuid.NewString())
	appAlarm.Status = api.Status_ACTIVE
	appAlarm.Type = api.AlarmType_APP

	smsAlarm := random.Alarm("ale", org.GetId(), uuid.NewString())
	smsAlarm.Status = api.Status_ACTIVE
	smsAlarm.Type = api.AlarmType_SMS

	emailAlarm := random.Alarm("ale", org.GetId(), uuid.NewString())
	emailAlarm.Status = api.Status_ACTIVE
	emailAlarm.Type = api.AlarmType_EMAIL

	disAlarm := random.Alarm("ale", org.GetId(), uuid.NewString())
	disAlarm.Status = api.Status_DISABLED

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
				Device: random.Device("ale", org.GetId()),
				Rule:   random.Rule("ale", org.GetId()),
			},
			[]*api.Alarm{appAlarm},
			[]*api.User{random.User("ale", org.GetId())},
			1, 1, 0, 0, 1, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.GetId()),
				Rule:   random.Rule("ale", org.GetId()),
			},
			[]*api.Alarm{smsAlarm},
			[]*api.User{random.User("ale", org.GetId())},
			1, 0, 1, 0, 1, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.GetId()),
				Rule:   random.Rule("ale", org.GetId()),
			}, []*api.Alarm{emailAlarm}, []*api.User{
				random.User("ale", org.GetId()),
			}, 1, 0, 0, 1, 1, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.GetId()),
				Rule:   random.Rule("ale", org.GetId()),
			},
			[]*api.Alarm{disAlarm},
			[]*api.User{random.User("ale", org.GetId())},
			0, 0, 0, 0, 0, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.GetId()),
				Rule:   random.Rule("ale", org.GetId()),
			}, []*api.Alarm{appAlarm}, []*api.User{}, 1, 0, 0, 0, 0, false,
		},
		{
			&message.EventerOut{
				Point:  &common.DataPoint{},
				Device: random.Device("ale", org.GetId()),
				Rule:   random.Rule("ale", org.GetId()),
			},
			[]*api.Alarm{appAlarm},
			[]*api.User{random.User("ale", org.GetId())},
			1, 0, 0, 0, 0, true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can alert %+v", test), func(t *testing.T) {
			t.Parallel()

			eOutQueue := queue.NewFake()
			eOutSub, err := eOutQueue.Subscribe("")
			require.NoError(t, err)

			var wg sync.WaitGroup
			wg.Add(1)

			ctrl := gomock.NewController(t)
			orger := NewMockorger(ctrl)
			orger.EXPECT().Read(gomock.Any(), test.inpEOut.GetDevice().
				GetOrgId()).Return(org, nil).Times(1)

			alarmer := NewMockalarmer(ctrl)
			alarmer.EXPECT().List(gomock.Any(), test.inpEOut.GetDevice().
				GetOrgId(), time.Time{}, "", int32(0), test.inpEOut.GetRule().
				GetId()).Return(test.inpAlarms, int32(0), nil).Times(1)

			userer := NewMockuserer(ctrl)
			userer.EXPECT().ListByTags(gomock.Any(), test.inpEOut.GetDevice().
				GetOrgId(), test.inpAlarms[0].GetUserTags()).Return(test.
				inpUsers, nil).Times(test.inpUserTimes)

			var alert *api.Alert
			if test.inpAlertTimes > 0 {
				alert = &api.Alert{
					OrgId:   test.inpEOut.GetDevice().GetOrgId(),
					UniqId:  test.inpEOut.GetDevice().GetUniqId(),
					AlarmId: test.inpAlarms[0].GetId(),
					UserId:  test.inpUsers[0].GetId(),
					Status:  api.AlertStatus_SENT,
					TraceId: test.inpEOut.GetPoint().GetTraceId(),
				}
			}

			notifier := notify.NewMockNotifier(ctrl)
			notifier.EXPECT().App(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(nil).Times(test.inpAppTimes)
			notifier.EXPECT().SMS(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(nil).Times(test.inpSMSTimes)
			notifier.EXPECT().Email(gomock.Any(), org.GetDisplayName(),
				org.GetEmail(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).Times(test.inpEmailTimes)

			alerter := NewMockalerter(ctrl)
			alerter.EXPECT().Create(gomock.Any(),
				matcher.NewProtoMatcher(alert)).
				DoAndReturn(func(_ interface{}, _ interface{}) error {
					defer wg.Done()

					return nil
				}).Times(test.inpAlertTimes)

			cache := cache.NewMemory()

			if test.inpSeedCache {
				ctx, cancel := context.WithTimeout(context.Background(),
					5*time.Second)
				defer cancel()

				key := repeatKey(test.inpEOut.GetDevice().GetOrgId(),
					test.inpEOut.GetDevice().GetId(),
					test.inpAlarms[0].GetId(), test.inpUsers[0].GetId())
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

			bEOut, err := proto.Marshal(test.inpEOut)
			require.NoError(t, err)
			t.Logf("bEOut: %s", bEOut)

			require.NoError(t, eOutQueue.Publish("", bEOut))
			if test.inpAlertTimes > 0 {
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
	appAlarm.Status = api.Status_ACTIVE
	appAlarm.Type = api.AlarmType_APP

	badSubj := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	badSubj.Status = api.Status_ACTIVE
	badSubj.SubjectTemplate = `{{if`

	badBody := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	badBody.Status = api.Status_ACTIVE
	badBody.BodyTemplate = `{{if`

	unspecType := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	unspecType.Status = api.Status_ACTIVE
	unspecType.Type = api.AlarmType_ALARM_TYPE_UNSPECIFIED

	unknownType := random.Alarm("ale", uuid.NewString(), uuid.NewString())
	unknownType.Status = api.Status_ACTIVE
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
				Point: &common.DataPoint{}, Device: &api.Device{},
			}, nil, nil, nil, nil, nil, nil, true, nil, nil, nil, 0, 0, 0, 0, 0,
			0,
		},
		// Orger error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, nil, errTestProc, nil, nil, nil, nil, true, nil, nil, nil, 1, 0,
			0, 0, 0, 0,
		},
		// Alarmer error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, nil, errTestProc, nil, nil, true, nil, nil, nil, 1, 1,
			0, 0, 0, 0,
		},
		// Userer error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil,
			[]*api.Alarm{appAlarm},
			nil, nil, errTestProc, true,
			nil, nil, nil, 1, 1, 1, 0, 0, 0,
		},
		// Bad alarm subject.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, []*api.Alarm{badSubj}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 0, 0, 0,
		},
		// Bad alarm body.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, []*api.Alarm{badBody}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 0, 0, 0,
		},
		// Cacher error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, []*api.Alarm{appAlarm}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, false, errTestProc, nil, nil, 1, 1, 1, 1, 0, 0,
		},
		// Notifier error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, []*api.Alarm{appAlarm}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, errTestProc, nil, 1, 1, 1, 1, 1, 1,
		},
		// Unspecified alarm type.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, []*api.Alarm{unspecType}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 1, 0, 1,
		},
		// Unknown alarm type.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, []*api.Alarm{unknownType}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, nil, 1, 1, 1, 1, 0, 1,
		},
		// Alerter error.
		{
			&message.EventerOut{
				Point: &common.DataPoint{}, Device: &api.Device{},
				Rule: &api.Rule{},
			}, org, nil, []*api.Alarm{appAlarm}, nil, []*api.User{
				random.User("ale", uuid.NewString()),
			}, nil, true, nil, nil, errTestProc, 1, 1, 1, 1, 1, 1,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can alert %+v", test), func(t *testing.T) {
			t.Parallel()

			eOutQueue := queue.NewFake()
			eOutSub, err := eOutQueue.Subscribe("")
			require.NoError(t, err)

			var wg sync.WaitGroup
			wg.Add(1)

			ctrl := gomock.NewController(t)
			orger := NewMockorger(ctrl)
			orger.EXPECT().Read(gomock.Any(), gomock.Any()).Return(test.inpOrg,
				test.inpOrgErr).Times(test.inpOrgTimes)

			alarmer := NewMockalarmer(ctrl)
			alarmer.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any(), gomock.Any()).
				Return(test.inpAlarms, int32(0), test.inpAlarmsErr).
				Times(test.inpAlarmTimes)

			userer := NewMockuserer(ctrl)
			userer.EXPECT().ListByTags(gomock.Any(), gomock.Any(),
				gomock.Any()).Return(test.inpUsers, test.inpUsersErr).
				Times(test.inpUserTimes)

			cacher := cache.NewMockCacher(ctrl)
			cacher.EXPECT().SetIfNotExistTTL(gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any()).Return(test.inpCache,
				test.inpCacheErr).Times(test.inpCacheTimes)

			notifier := notify.NewMockNotifier(ctrl)
			notifier.EXPECT().App(gomock.Any(), gomock.Any(), gomock.Any(),
				gomock.Any()).Return(test.inpAppErr).Times(test.inpAppTimes)

			alerter := NewMockalerter(ctrl)
			alerter.EXPECT().Create(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ interface{}, _ interface{}) error {
					defer wg.Done()

					return test.inpAlertErr
				}).Times(test.inpAlertTimes)

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
			if test.inpEOut != nil {
				bEOut, err = proto.Marshal(test.inpEOut)
				require.NoError(t, err)
				t.Logf("bEOut: %s", bEOut)
			}

			require.NoError(t, eOutQueue.Publish("", bEOut))
			if test.inpAlertTimes > 0 {
				wg.Wait()
			} else {
				// If the success mode isn't supported by WaitGroup operation,
				// give it time to traverse the code.
				time.Sleep(100 * time.Millisecond)
			}
		})
	}
}
