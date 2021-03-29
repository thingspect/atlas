// +build !unit

package test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListAlerts(t *testing.T) {
	t.Parallel()

	t.Run("List alerts by UniqID, dev, alarm, and user ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		daleCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := daleCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-alert", uuid.NewString())})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		alerts := []*api.Alert{}

		for i := 0; i < 5; i++ {
			alert := random.Alert("dao-alert", globalAdminOrgID)
			alert.UniqId = createDev.UniqId
			alerts = append(alerts, alert)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			err := globalAleDAO.Create(ctx, alert)
			t.Logf("err: %v", err)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		sort.Slice(alerts, func(i, j int) bool {
			return alerts[i].CreatedAt.AsTime().After(
				alerts[j].CreatedAt.AsTime())
		})

		ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		// Verify results by UniqID.
		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlertsUniqID, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_UniqId{UniqId: createDev.UniqId},
			EndTime: alerts[0].CreatedAt, StartTime: timestamppb.New(
				alerts[len(alerts)-1].CreatedAt.AsTime().Add(
					-time.Millisecond))})
		t.Logf("listAlertsUniqID, err: %+v, %v", listAlertsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsUniqID.Alerts, len(alerts))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListAlertsResponse{Alerts: alerts},
			listAlertsUniqID) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListAlertsResponse{Alerts: alerts}, listAlertsUniqID)
		}

		// Verify results by dev ID without oldest alert.
		listAlertsDevID, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof:   &api.ListAlertsRequest_DeviceId{DeviceId: createDev.Id},
			StartTime: alerts[len(alerts)-1].CreatedAt})
		t.Logf("listAlertsDevID, err: %+v, %v", listAlertsDevID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsDevID.Alerts, len(alerts)-1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListAlertsResponse{Alerts: alerts[:len(alerts)-1]},
			listAlertsDevID) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListAlertsResponse{
				Alerts: alerts[:len(alerts)-1]}, listAlertsDevID)
		}

		// Verify results by alarm ID and user ID.
		listAlertsUniqID, err = aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			AlarmId: alerts[len(alerts)-1].AlarmId,
			UserId:  alerts[len(alerts)-1].UserId, EndTime: alerts[0].CreatedAt,
			StartTime: timestamppb.New(alerts[len(alerts)-1].CreatedAt.AsTime().
				Add(-time.Millisecond))})
		t.Logf("listAlertsUniqID, err: %+v, %v", listAlertsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsUniqID.Alerts, 1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListAlertsResponse{Alerts: []*api.Alert{
			alerts[len(alerts)-1]}}, listAlertsUniqID) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListAlertsResponse{
				Alerts: []*api.Alert{alerts[len(alerts)-1]}}, listAlertsUniqID)
		}
	})

	t.Run("List alerts are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-alert"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		alert := random.Alert("dao-alert", createOrg.Id)

		err = globalAleDAO.Create(ctx, alert)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlerts, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_UniqId{UniqId: alert.UniqId}})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.NoError(t, err)
		require.Len(t, listAlerts.Alerts, 0)
	})

	t.Run("List alerts by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlerts, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_DeviceId{
				DeviceId: random.String(10)}, EndTime: timestamppb.Now(),
			StartTime: timestamppb.New(time.Now().Add(-91 * 24 * time.Hour))})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"maximum time range exceeded")
	})

	t.Run("List alerts by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		aleCli := api.NewAlertServiceClient(globalAdminGRPCConn)
		listAlerts, err := aleCli.ListAlerts(ctx, &api.ListAlertsRequest{
			IdOneof: &api.ListAlertsRequest_DeviceId{
				DeviceId: random.String(10)}})
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: UUID")
	})
}
