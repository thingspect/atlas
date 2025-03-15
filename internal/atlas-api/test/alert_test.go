//go:build !unit

package test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListAlerts(t *testing.T) {
	t.Parallel()

	t.Run("List alerts by UniqID, dev, alarm, and user ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-alert", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		alerts := []*api.Alert{}

		for range 5 {
			alert := random.Alert("dao-alert", globalAdminOrgID)
			alert.UniqId = createDev.GetUniqId()
			alerts = append(alerts, alert)

			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()

			err := globalAleDAO.Create(ctx, alert)
			t.Logf("err: %v", err)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		// Flip alerts to descending timestamp order.
		sort.Slice(alerts, func(i, j int) bool {
			return alerts[i].GetCreatedAt().AsTime().After(
				alerts[j].GetCreatedAt().AsTime())
		})

		ctx, cancel = context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		// Verify results by UniqID.
		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlertsUniqID, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_UniqId{UniqId: createDev.GetUniqId()},
			EndTime: alerts[0].GetCreatedAt(), StartTime: timestamppb.New(
				alerts[len(alerts)-1].GetCreatedAt().AsTime().Add(
					-time.Millisecond)),
		})
		t.Logf("listAlertsUniqID, err: %+v, %v", listAlertsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsUniqID.GetAlerts(), len(alerts))
		require.EqualExportedValues(t, &api.ListAlertsResponse{Alerts: alerts},
			listAlertsUniqID)

		// Verify results by dev ID without oldest alert.
		listAlertsDevID, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof:   &api.ListAlertsRequest_DeviceId{DeviceId: createDev.GetId()},
			StartTime: alerts[len(alerts)-1].GetCreatedAt(),
		})
		t.Logf("listAlertsDevID, err: %+v, %v", listAlertsDevID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsDevID.GetAlerts(), len(alerts)-1)
		require.EqualExportedValues(t, &api.ListAlertsResponse{
			Alerts: alerts[:len(alerts)-1],
		}, listAlertsDevID)

		// Verify results by alarm ID and user ID.
		listAlertsUniqID, err = aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			AlarmId: alerts[len(alerts)-1].GetAlarmId(),
			UserId:  alerts[len(alerts)-1].GetUserId(),
			EndTime: alerts[0].GetCreatedAt(),
			StartTime: timestamppb.New(alerts[len(alerts)-1].GetCreatedAt().
				AsTime().Add(-time.Millisecond)),
		})
		t.Logf("listAlertsUniqID, err: %+v, %v", listAlertsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsUniqID.GetAlerts(), 1)
		require.EqualExportedValues(t, &api.ListAlertsResponse{
			Alerts: []*api.Alert{alerts[len(alerts)-1]},
		}, listAlertsUniqID)
	})

	t.Run("List alerts are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-alert"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		alert := random.Alert("dao-alert", createOrg.GetId())

		err = globalAleDAO.Create(ctx, alert)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlerts, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_UniqId{UniqId: alert.GetUniqId()},
		})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.NoError(t, err)
		require.Empty(t, listAlerts.GetAlerts())
	})

	t.Run("List alerts by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlerts, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_DeviceId{
				DeviceId: uuid.NewString(),
			}, EndTime: timestamppb.Now(),
			StartTime: timestamppb.New(time.Now().Add(-91 * 24 * time.Hour)),
		})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"maximum time range exceeded")
	})

	t.Run("List alerts by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlerts, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_DeviceId{
				DeviceId: random.String(10),
			},
		})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid ListAlertsRequest.DeviceId: value must be a valid UUID | "+
			"caused by: invalid uuid format")
	})
}
