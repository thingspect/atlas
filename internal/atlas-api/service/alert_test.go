//go:build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListAlerts(t *testing.T) {
	t.Parallel()

	t.Run("List alerts by valid UniqID with ts", func(t *testing.T) {
		t.Parallel()

		alert := random.Alert("dao-alert", uuid.NewString())
		retAlert, _ := proto.Clone(alert).(*api.Alert)
		end := time.Now().UTC()
		start := time.Now().UTC().Add(-15 * time.Minute)

		alerter := NewMockAlerter(gomock.NewController(t))
		alerter.EXPECT().List(gomock.Any(), alert.GetOrgId(), alert.GetUniqId(), "", "",
			"", end, start).Return([]*api.Alert{retAlert}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: alert.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		aleSvc := NewAlert(alerter)
		listAlerts, err := aleSvc.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_UniqId{UniqId: alert.GetUniqId()},
			EndTime: timestamppb.New(end), StartTime: timestamppb.New(start),
		})
		t.Logf("alert, listAlerts, err: %+v, %+v, %v", alert, listAlerts, err)
		require.NoError(t, err)
		require.EqualExportedValues(t,
			&api.ListAlertsResponse{Alerts: []*api.Alert{alert}}, listAlerts)
	})

	t.Run("List alerts by valid dev ID with alarm ID", func(t *testing.T) {
		t.Parallel()

		alert := random.Alert("dao-alert", uuid.NewString())
		retAlert, _ := proto.Clone(alert).(*api.Alert)
		devID := uuid.NewString()
		alarmID := uuid.NewString()

		alerter := NewMockAlerter(gomock.NewController(t))
		alerter.EXPECT().List(gomock.Any(), alert.GetOrgId(), "", devID,
			alarmID, "", matcher.NewRecentMatcher(2*time.Second),
			matcher.NewRecentMatcher(24*time.Hour+2*time.Second)).
			Return([]*api.Alert{retAlert}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: alert.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		aleSvc := NewAlert(alerter)
		listAlerts, err := aleSvc.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_DeviceId{DeviceId: devID},
			AlarmId: alarmID,
		})
		t.Logf("alert, listAlerts, err: %+v, %+v, %v", alert, listAlerts, err)
		require.NoError(t, err)
		require.EqualExportedValues(t,
			&api.ListAlertsResponse{Alerts: []*api.Alert{alert}}, listAlerts)
	})

	t.Run("List alerts by user ID", func(t *testing.T) {
		t.Parallel()

		alert := random.Alert("dao-alert", uuid.NewString())
		retAlert, _ := proto.Clone(alert).(*api.Alert)
		userID := uuid.NewString()

		alerter := NewMockAlerter(gomock.NewController(t))
		alerter.EXPECT().List(gomock.Any(), alert.GetOrgId(), "", "",
			"", userID, matcher.NewRecentMatcher(2*time.Second),
			matcher.NewRecentMatcher(24*time.Hour+2*time.Second)).
			Return([]*api.Alert{retAlert}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: alert.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		aleSvc := NewAlert(alerter)
		listAlerts, err := aleSvc.ListAlerts(ctx,
			&api.ListAlertsRequest{UserId: userID})
		t.Logf("alert, listAlerts, err: %+v, %+v, %v", alert, listAlerts, err)
		require.NoError(t, err)
		require.EqualExportedValues(t,
			&api.ListAlertsResponse{Alerts: []*api.Alert{alert}}, listAlerts)
	})

	t.Run("List alerts with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		aleSvc := NewAlert(nil)
		listAlerts, err := aleSvc.ListAlerts(ctx, &api.ListAlertsRequest{})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List alerts with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		aleSvc := NewAlert(nil)
		listAlerts, err := aleSvc.ListAlerts(ctx, &api.ListAlertsRequest{})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List alerts by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		aleSvc := NewAlert(nil)
		listAlerts, err := aleSvc.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_UniqId{UniqId: "api-alert-" +
				random.String(16)}, EndTime: timestamppb.Now(),
			StartTime: timestamppb.New(time.Now().Add(-91 * 24 * time.Hour)),
		})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"maximum time range exceeded"), err)
	})

	t.Run("List alerts by invalid org ID", func(t *testing.T) {
		t.Parallel()

		alerter := NewMockAlerter(gomock.NewController(t))
		alerter.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		aleSvc := NewAlert(alerter)
		listAlerts, err := aleSvc.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_UniqId{UniqId: "api-alert-" +
				random.String(16)},
		})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"dao: invalid format"), err)
	})
}
