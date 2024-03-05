//go:build !unit

package event

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
)

const testTimeout = 6 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid events", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		err = globalEvDAO.Create(ctx, random.Event("dao-event", createOrg.GetId()))
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Create invalid event", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		badUniqID := random.Event("dao-event", createOrg.GetId())
		badUniqID.UniqId = "dao-event-" + random.String(40)

		badTraceID := random.Event("dao-event", createOrg.GetId())
		badTraceID.TraceId = random.String(10)

		tests := []struct {
			inpEvent *api.Event
			err      error
		}{
			{badUniqID, dao.ErrInvalidFormat},
			{badTraceID, dao.ErrInvalidFormat},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Cannot create %+v", test), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(context.Background(),
					testTimeout)
				defer cancel()

				err := globalEvDAO.Create(ctx, test.inpEvent)
				t.Logf("err: %#v", err)
				require.ErrorIs(t, err, test.err)
			})
		}
	})

	t.Run("Dedupe event", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		event := random.Event("dao-event", createOrg.GetId())

		err = globalEvDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		err = globalEvDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.Equal(t, dao.ErrAlreadyExists, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("List events by UniqID, dev ID, and rule ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-event",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		events := []*api.Event{}

		for range 5 {
			event := random.Event("dao-event", createOrg.GetId())
			event.UniqId = createDev.GetUniqId()
			events = append(events, event)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			err := globalEvDAO.Create(ctx, event)
			t.Logf("err: %v", err)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		// Flip events to descending timestamp order.
		sort.Slice(events, func(i, j int) bool {
			return events[i].GetCreatedAt().AsTime().After(
				events[j].GetCreatedAt().AsTime())
		})

		ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		// Verify results by UniqID.
		listEventsUniqID, err := globalEvDAO.List(ctx, createOrg.GetId(),
			createDev.GetUniqId(), "", "", events[0].GetCreatedAt().AsTime(),
			events[len(events)-1].GetCreatedAt().AsTime().
				Add(-time.Millisecond))
		t.Logf("listEventsUniqID, err: %+v, %v", listEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listEventsUniqID, len(events))

		for i, event := range events {
			require.EqualExportedValues(t, event, listEventsUniqID[i])
		}

		// Verify results by dev ID without oldest event.
		listEventsDevID, err := globalEvDAO.List(ctx, createOrg.GetId(), "",
			createDev.GetId(), "", events[0].GetCreatedAt().AsTime(),
			events[len(events)-1].GetCreatedAt().AsTime())
		t.Logf("listEventsDevID, err: %+v, %v", listEventsDevID, err)
		require.NoError(t, err)
		require.Len(t, listEventsDevID, len(events)-1)

		for i, event := range events[:len(events)-1] {
			require.EqualExportedValues(t, event, listEventsDevID[i])
		}

		// Verify results by UniqID and rule ID.
		listEventsUniqID, err = globalEvDAO.List(ctx, createOrg.GetId(),
			createDev.GetUniqId(), "", events[len(events)-1].GetRuleId(),
			events[0].GetCreatedAt().AsTime(),
			events[len(events)-1].GetCreatedAt().AsTime().
				Add(-time.Millisecond))
		t.Logf("listEventsUniqID, err: %+v, %v", listEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listEventsUniqID, 1)
		require.EqualExportedValues(t, events[len(events)-1],
			listEventsUniqID[0])
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		event := random.Event("dao-event", createOrg.GetId())

		err = globalEvDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		listEvents, err := globalEvDAO.List(ctx, uuid.NewString(),
			event.GetUniqId(), "", "", event.GetCreatedAt().AsTime(),
			event.GetCreatedAt().AsTime().Add(-time.Millisecond))
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.NoError(t, err)
		require.Empty(t, listEvents)
	})

	t.Run("List events by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listEvents, err := globalEvDAO.List(ctx, random.String(10),
			uuid.NewString(), "", "", time.Now(), time.Now())
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.Nil(t, listEvents)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestLatest(t *testing.T) {
	t.Parallel()

	t.Run("Latest events", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		events := []*api.Event{}

		for range 5 {
			event := random.Event("dao-event", createOrg.GetId())
			events = append(events, event)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			err := globalEvDAO.Create(ctx, event)
			t.Logf("err: %v", err)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		// Flip events to descending timestamp order.
		sort.Slice(events, func(i, j int) bool {
			return events[i].GetCreatedAt().AsTime().After(
				events[j].GetCreatedAt().AsTime())
		})

		ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		// Verify results.
		latEvents, err := globalEvDAO.Latest(ctx, createOrg.GetId(), "")
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.NoError(t, err)
		require.Len(t, latEvents, len(events))

		for i, event := range events {
			require.EqualExportedValues(t, event, latEvents[i])
		}

		// Verify results by rule ID.
		latEventsRuleID, err := globalEvDAO.Latest(ctx, createOrg.GetId(),
			events[len(events)-1].GetRuleId())
		t.Logf("latEventsRuleID, err: %+v, %v", latEventsRuleID, err)
		require.NoError(t, err)
		require.Len(t, latEventsRuleID, 1)
		require.EqualExportedValues(t, events[len(events)-1],
			latEventsRuleID[0])
	})

	t.Run("Latest events are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		event := random.Event("dao-event", createOrg.GetId())

		err = globalEvDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		latEvents, err := globalEvDAO.Latest(ctx, uuid.NewString(), "")
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.NoError(t, err)
		require.Empty(t, latEvents)
	})

	t.Run("Latest events by invalid rule ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		latEvents, err := globalEvDAO.Latest(ctx, uuid.NewString(),
			random.String(10))
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.Nil(t, latEvents)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
