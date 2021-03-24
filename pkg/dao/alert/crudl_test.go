// +build !unit

package alert

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 6 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid alerts", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alert"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		err = globalAlertDAO.Create(ctx, random.Alert("dao-alert", createOrg.Id))
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create invalid alert", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alert"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		badUniqID := random.Alert("dao-alert", createOrg.Id)
		badUniqID.UniqId = "dao-alert-" + random.String(40)

		badTraceID := random.Alert("dao-alert", createOrg.Id)
		badTraceID.TraceId = random.String(10)

		tests := []struct {
			inpAlert *api.Alert
			err      error
		}{
			{badUniqID, dao.ErrInvalidFormat},
			{badTraceID, dao.ErrInvalidFormat},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Cannot create %+v", lTest), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(context.Background(),
					testTimeout)
				defer cancel()

				err := globalAlertDAO.Create(ctx, lTest.inpAlert)
				t.Logf("err: %#v", err)
				require.ErrorIs(t, err, lTest.err)
			})
		}
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("List alerts by UniqID, dev, alarm, and user ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alert"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-alert",
			createOrg.Id))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		alerts := []*api.Alert{}

		for i := 0; i < 5; i++ {
			alert := random.Alert("dao-alert", createOrg.Id)
			alert.UniqId = createDev.UniqId
			alerts = append(alerts, alert)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			err := globalAlertDAO.Create(ctx, alert)
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
		listAlertsUniqID, err := globalAlertDAO.List(ctx, createOrg.Id,
			createDev.UniqId, "", "", "", alerts[0].CreatedAt.AsTime(),
			alerts[len(alerts)-1].CreatedAt.AsTime().Add(-time.Millisecond))
		t.Logf("listAlertsUniqID, err: %+v, %v", listAlertsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsUniqID, len(alerts))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, alert := range alerts {
			if !proto.Equal(alert, listAlertsUniqID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", alert,
					listAlertsUniqID[i])
			}
		}

		// Verify results by dev ID without oldest alert.
		listAlertsDevID, err := globalAlertDAO.List(ctx, createOrg.Id, "",
			createDev.Id, "", "", alerts[0].CreatedAt.AsTime(),
			alerts[len(alerts)-1].CreatedAt.AsTime())
		t.Logf("listAlertsDevID, err: %+v, %v", listAlertsDevID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsDevID, len(alerts)-1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, alert := range alerts[:len(alerts)-1] {
			if !proto.Equal(alert, listAlertsDevID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", alert,
					listAlertsDevID[i])
			}
		}

		// Verify results by alarm ID and user ID.
		listAlertsAlarmID, err := globalAlertDAO.List(ctx, createOrg.Id, "", "",
			alerts[len(alerts)-1].AlarmId, alerts[len(alerts)-1].UserId,
			alerts[0].CreatedAt.AsTime(),
			alerts[len(alerts)-1].CreatedAt.AsTime().Add(-time.Millisecond))
		t.Logf("listAlertsAlarmID, err: %+v, %v", listAlertsAlarmID, err)
		require.NoError(t, err)
		require.Len(t, listAlertsAlarmID, 1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(alerts[len(alerts)-1], listAlertsAlarmID[0]) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", alerts[len(alerts)-1],
				listAlertsAlarmID[0])
		}
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-alert"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		alert := random.Alert("dao-alert", createOrg.Id)

		err = globalAlertDAO.Create(ctx, alert)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		listAlerts, err := globalAlertDAO.List(ctx, uuid.NewString(),
			alert.UniqId, "", "", "", alert.CreatedAt.AsTime(),
			alert.CreatedAt.AsTime().Add(-time.Millisecond))
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.NoError(t, err)
		require.Len(t, listAlerts, 0)
	})

	t.Run("List alerts by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listAlerts, err := globalAlertDAO.List(ctx, random.String(10),
			uuid.NewString(), "", "", "", time.Now(), time.Now())
		t.Logf("listAlerts, err: %+v, %v", listAlerts, err)
		require.Nil(t, listAlerts)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
