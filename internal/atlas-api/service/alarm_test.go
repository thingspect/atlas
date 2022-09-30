//go:build !integration

package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Create valid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		retAlarm, _ := proto.Clone(alarm).(*api.Alarm)

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Create(gomock.Any(), alarm).Return(retAlarm, nil).
			Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: alarm.OrgId, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		createAlarm, err := raSvc.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: alarm,
		})
		t.Logf("alarm, createAlarm, err: %+v, %+v, %v", alarm, createAlarm, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(alarm, createAlarm) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", alarm, createAlarm)
		}
	})

	t.Run("Create alarm with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		createAlarm, err := raSvc.CreateAlarm(ctx, &api.CreateAlarmRequest{})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		createAlarm, err := raSvc.CreateAlarm(ctx, &api.CreateAlarmRequest{})
		t.Logf("createAlarm, err: %+v, %v", createAlarm, err)
		require.Nil(t, createAlarm)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Create invalid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		alarm.Name = random.String(81)

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Create(gomock.Any(), alarm).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: alarm.OrgId, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		createAlarm, err := raSvc.CreateAlarm(ctx, &api.CreateAlarmRequest{
			Alarm: alarm,
		})
		t.Logf("alarm, createAlarm, err: %+v, %+v, %v", alarm, createAlarm, err)
		require.Nil(t, createAlarm)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestGetAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Get alarm by valid ID", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		retAlarm, _ := proto.Clone(alarm).(*api.Alarm)

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Read(gomock.Any(), alarm.Id, alarm.OrgId,
			alarm.RuleId).Return(retAlarm, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: alarm.OrgId, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		getAlarm, err := raSvc.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: alarm.Id, RuleId: alarm.RuleId,
		})
		t.Logf("rule, getAlarm, err: %+v, %+v, %v", alarm, getAlarm, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(alarm, getAlarm) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", alarm, getAlarm)
		}
	})

	t.Run("Get rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		getAlarm, err := raSvc.GetAlarm(ctx, &api.GetAlarmRequest{})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Get rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		getAlarm, err := raSvc.GetAlarm(ctx, &api.GetAlarmRequest{})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Get rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		getAlarm, err := raSvc.GetAlarm(ctx, &api.GetAlarmRequest{
			Id: uuid.NewString(), RuleId: uuid.NewString(),
		})
		t.Logf("getAlarm, err: %+v, %v", getAlarm, err)
		require.Nil(t, getAlarm)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestUpdateAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		retAlarm, _ := proto.Clone(alarm).(*api.Alarm)

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Update(gomock.Any(), alarm).Return(retAlarm, nil).
			Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: alarm.OrgId, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: alarm,
		})
		t.Logf("alarm, updateAlarm, err: %+v, %+v, %v", alarm, updateAlarm, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(alarm, updateAlarm) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", alarm, updateAlarm)
		}
	})

	t.Run("Partial update alarm by valid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		retAlarm, _ := proto.Clone(alarm).(*api.Alarm)
		part := &api.Alarm{
			Id: alarm.Id, RuleId: alarm.RuleId, Status: api.Status_ACTIVE,
			SubjectTemplate: `test`, UserTags: random.Tags("api-alarm", 2),
		}
		merged := &api.Alarm{
			Id: alarm.Id, OrgId: alarm.OrgId, RuleId: alarm.RuleId,
			Name: alarm.Name, Status: part.Status, Type: alarm.Type,
			UserTags: part.UserTags, SubjectTemplate: part.SubjectTemplate,
			BodyTemplate:   alarm.BodyTemplate,
			RepeatInterval: alarm.RepeatInterval,
		}
		retMerged, _ := proto.Clone(merged).(*api.Alarm)

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Read(gomock.Any(), alarm.Id, alarm.OrgId,
			alarm.RuleId).Return(retAlarm, nil).Times(1)
		alarmer.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(merged)).
			Return(retMerged, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: alarm.OrgId, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status", "subject_template", "user_tags"},
			},
		})
		t.Logf("merged, updateAlarm, err: %+v, %+v, %v", merged, updateAlarm,
			err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(merged, updateAlarm) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", merged, updateAlarm)
		}
	})

	t.Run("Update alarm with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Update alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Update nil alarm", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: nil,
		})
		t.Logf("updateAlarm, err: %+v, %v", updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateAlarmRequest.Alarm: value is required"), err)
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: alarm, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"},
			},
		})
		t.Logf("alarm, updateAlarm, err: %+v, %+v, %v", alarm, updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid field mask"), err)
	})

	t.Run("Partial update alarm by unknown alarm", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		part := &api.Alarm{
			Id: uuid.NewString(), RuleId: uuid.NewString(),
			Status: api.Status_ACTIVE,
		}

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Read(gomock.Any(), part.Id, orgID, part.RuleId).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"},
			},
		})
		t.Logf("part, updateAlarm, err: %+v, %+v, %v", part, updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Update alarm validation failure", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())
		alarm.Name = random.String(81)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: alarm.OrgId, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: alarm,
		})
		t.Logf("alarm, updateAlarm, err: %+v, %+v, %v", alarm, updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid "+
			"UpdateAlarmRequest.Alarm: embedded message failed validation | "+
			"caused by: invalid Alarm.Name: value length must be between 5 "+
			"and 80 runes, inclusive"), err)
	})

	t.Run("Update alarm by invalid alarm", func(t *testing.T) {
		t.Parallel()

		alarm := random.Alarm("api-alarm", uuid.NewString(), uuid.NewString())

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Update(gomock.Any(), alarm).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: alarm.OrgId, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		updateAlarm, err := raSvc.UpdateAlarm(ctx, &api.UpdateAlarmRequest{
			Alarm: alarm,
		})
		t.Logf("alarm, updateAlarm, err: %+v, %+v, %v", alarm, updateAlarm, err)
		require.Nil(t, updateAlarm)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestDeleteAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Delete rule by valid ID", func(t *testing.T) {
		t.Parallel()

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		_, err := raSvc.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete rule with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		_, err := raSvc.DeleteAlarm(ctx, &api.DeleteAlarmRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete rule with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		_, err := raSvc.DeleteAlarm(ctx, &api.DeleteAlarmRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Delete rule by unknown ID", func(t *testing.T) {
		t.Parallel()

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any()).Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		_, err := raSvc.DeleteAlarm(ctx, &api.DeleteAlarmRequest{
			Id: uuid.NewString(),
		})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestListAlarms(t *testing.T) {
	t.Parallel()

	t.Run("List alarms by valid org ID", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		alarms := []*api.Alarm{
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
		}

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(51),
			"").Return(alarms, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		listAlarms, err := raSvc.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listAlarms.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListAlarmsResponse{Alarms: alarms, TotalSize: 3},
			listAlarms) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListAlarmsResponse{Alarms: alarms, TotalSize: 3},
				listAlarms)
		}
	})

	t.Run("List alarms by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		alarms := []*api.Alarm{
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
		}

		next, err := session.GeneratePageToken(alarms[1].CreatedAt.AsTime(),
			alarms[1].Id)
		require.NoError(t, err)

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3),
			"").Return(alarms, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		listAlarms, err := raSvc.ListAlarms(ctx, &api.ListAlarmsRequest{
			PageSize: 2,
		})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listAlarms.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListAlarmsResponse{
			Alarms: alarms[:2], NextPageToken: next, TotalSize: 3,
		}, listAlarms) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListAlarmsResponse{
				Alarms: alarms[:2], NextPageToken: next, TotalSize: 3,
			}, listAlarms)
		}
	})

	t.Run("List alarms with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		listAlarms, err := raSvc.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.Nil(t, listAlarms)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List alarms with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		listAlarms, err := raSvc.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.Nil(t, listAlarms)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List alarms by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		listAlarms, err := raSvc.ListAlarms(ctx, &api.ListAlarmsRequest{
			PageToken: badUUID,
		})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.Nil(t, listAlarms)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid page token"), err)
	})

	t.Run("List alarms by invalid org ID", func(t *testing.T) {
		t.Parallel()

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any()).Return(nil, int32(0),
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		listAlarms, err := raSvc.ListAlarms(ctx, &api.ListAlarmsRequest{})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.Nil(t, listAlarms)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("List alarms with generation failure", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		alarms := []*api.Alarm{
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
			random.Alarm("api-alarm", uuid.NewString(), uuid.NewString()),
		}
		alarms[1].Id = badUUID

		alarmer := NewMockAlarmer(gomock.NewController(t))
		alarmer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3),
			"").Return(alarms, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, alarmer)
		listAlarms, err := raSvc.ListAlarms(ctx, &api.ListAlarmsRequest{
			PageSize: 2,
		})
		t.Logf("listAlarms, err: %+v, %v", listAlarms, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listAlarms.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListAlarmsResponse{
			Alarms: alarms[:2], TotalSize: 3,
		}, listAlarms) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListAlarmsResponse{
				Alarms: alarms[:2], TotalSize: 3,
			}, listAlarms)
		}
	})
}

func TestTestAlarm(t *testing.T) {
	t.Parallel()

	t.Run("Test valid and invalid alarms", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			inpPoint *common.DataPoint
			inpRule  *api.Rule
			inpDev   *api.Device
			inpTempl string
			res      string
			err      string
		}{
			{
				&common.DataPoint{}, nil, nil, `test`, "test", "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_IntVal{IntVal: 40},
				}, &api.Rule{Name: "test rule"}, nil, `point value is an ` +
					`integer: {{.pointVal}}, rule name is: {{.rule.Name}}`,
				"point value is an integer: 40, rule name is: test rule", "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_Fl64Val{Fl64Val: 37.7},
				}, nil, &api.Device{Status: api.Status_ACTIVE}, `point ` +
					`value is a float: {{.pointVal}}, device status is: ` +
					`{{.device.Status}}`,
				"point value is a float: 37.7, device status is: ACTIVE", "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_StrVal{StrVal: "line"},
				}, nil, nil, `point value is a string: {{.pointVal}}`,
				"point value is a string: line", "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_BoolVal{BoolVal: true},
				}, nil, nil, `point value is a bool: {{.pointVal}}`,
				"point value is a bool: true", "",
			},
			{
				&common.DataPoint{
					ValOneof: &common.DataPoint_BytesVal{
						BytesVal: []byte{0x00, 0x01},
					},
				}, nil, nil, `point value is a byte slice: {{.pointVal}}`,
				"point value is a byte slice: [0 1]", "",
			},
			{
				&common.DataPoint{}, nil, nil, `{{if`, "", "unclosed action",
			},
			{
				&common.DataPoint{}, nil, nil, `{{template "aaa"}}`, "",
				"no such template",
			},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Can evaluate %+v", lTest), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(session.NewContext(
					context.Background(), &session.Session{
						OrgID: uuid.NewString(), Role: api.Role_ADMIN,
					}),
					testTimeout)
				defer cancel()

				raSvc := NewRuleAlarm(nil, nil)
				testRes, err := raSvc.TestAlarm(ctx, &api.TestAlarmRequest{
					Point: lTest.inpPoint, Rule: lTest.inpRule,
					Device: lTest.inpDev, Alarm: &api.Alarm{
						SubjectTemplate: lTest.inpTempl,
						BodyTemplate:    lTest.inpTempl,
					},
				})
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

	t.Run("Test alarm with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		testRes, err := raSvc.TestAlarm(ctx, &api.TestAlarmRequest{})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Test alarm with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		testRes, err := raSvc.TestAlarm(ctx, &api.TestAlarmRequest{})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, errPerm(api.Role_BUILDER), err)
	})

	t.Run("Test alarm with invalid body template", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		raSvc := NewRuleAlarm(nil, nil)
		testRes, err := raSvc.TestAlarm(ctx, &api.TestAlarmRequest{
			Point: &common.DataPoint{}, Rule: nil, Device: nil,
			Alarm: &api.Alarm{SubjectTemplate: `test`, BodyTemplate: `{{if`},
		})
		t.Logf("testRes, err: %+v, %v", testRes, err)
		require.Nil(t, testRes)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"template: template:1: unclosed action"), err)
	})
}
