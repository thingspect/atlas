// +build !unit

package event

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
	"google.golang.org/protobuf/types/known/timestamppb"
)

const testTimeout = 8 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid events", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			lTest := i

			t.Run(fmt.Sprintf("Can create %v", lTest), func(t *testing.T) {
				t.Parallel()

				event := &api.Event{OrgId: createOrg.Id,
					RuleId: uuid.NewString(), UniqId: "dao-event-" +
						random.String(16), CreatedAt: timestamppb.Now(),
					TraceId: uuid.NewString()}

				ctx, cancel := context.WithTimeout(context.Background(),
					testTimeout)
				defer cancel()

				err := globalEventDAO.Create(ctx, event)
				t.Logf("err: %v", err)
				require.NoError(t, err)
			})
		}
	})

	t.Run("Create invalid event", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		tests := []struct {
			inpEvent *api.Event
			err      error
		}{
			{&api.Event{OrgId: createOrg.Id, RuleId: uuid.NewString(),
				UniqId:    "dao-event-" + random.String(40),
				CreatedAt: timestamppb.Now(), TraceId: uuid.NewString()},
				dao.ErrInvalidFormat},
			{&api.Event{OrgId: createOrg.Id, RuleId: uuid.NewString(),
				UniqId:    "dao-event-" + random.String(16),
				CreatedAt: timestamppb.Now(), TraceId: random.String(10)},
				dao.ErrInvalidFormat},
		}

		for _, test := range tests {
			lTest := test

			t.Run(fmt.Sprintf("Cannot create %+v", lTest), func(t *testing.T) {
				t.Parallel()

				ctx, cancel := context.WithTimeout(context.Background(),
					testTimeout)
				defer cancel()

				err := globalEventDAO.Create(ctx, lTest.inpEvent)
				t.Logf("err: %#v", err)
				require.ErrorIs(t, err, lTest.err)
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

		event := &api.Event{OrgId: createOrg.Id, RuleId: uuid.NewString(),
			UniqId:    "dao-event-" + random.String(16),
			CreatedAt: timestamppb.Now(), TraceId: uuid.NewString()}

		err = globalEventDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		err = globalEventDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.Equal(t, dao.ErrAlreadyExists, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("List events by UniqID, dev ID, and attr", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-event",
			createOrg.Id))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		events := []*api.Event{}

		for i := 0; i < 5; i++ {
			event := &api.Event{OrgId: createOrg.Id, RuleId: uuid.NewString(),
				UniqId: createDev.UniqId, CreatedAt: timestamppb.New(
					time.Now().UTC().Truncate(time.Millisecond)),
				TraceId: uuid.NewString()}
			events = append(events, event)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			err := globalEventDAO.Create(ctx, event)
			t.Logf("err: %v", err)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		sort.Slice(events, func(i, j int) bool {
			return events[i].CreatedAt.AsTime().After(
				events[j].CreatedAt.AsTime())
		})

		ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		// Verify results by UniqID.
		listEventsUniqID, err := globalEventDAO.List(ctx, createOrg.Id,
			createDev.UniqId, "", "", events[0].CreatedAt.AsTime(),
			events[len(events)-1].CreatedAt.AsTime().Add(-time.Millisecond))
		t.Logf("listEventsUniqID, err: %+v, %v", listEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listEventsUniqID, len(events))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, event := range events {
			if !proto.Equal(event, listEventsUniqID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", event,
					listEventsUniqID[i])
			}
		}

		// Verify results by dev ID without oldest event.
		listEventsDevID, err := globalEventDAO.List(ctx, createOrg.Id, "",
			createDev.Id, "", events[0].CreatedAt.AsTime(),
			events[len(events)-1].CreatedAt.AsTime())
		t.Logf("listEventsDevID, err: %+v, %v", listEventsDevID, err)
		require.NoError(t, err)
		require.Len(t, listEventsDevID, len(events)-1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, event := range events[:len(events)-1] {
			if !proto.Equal(event, listEventsDevID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", event,
					listEventsDevID[i])
			}
		}

		// Verify results by UniqID and rule ID.
		listEventsUniqID, err = globalEventDAO.List(ctx, createOrg.Id,
			createDev.UniqId, "", events[len(events)-1].RuleId,
			events[0].CreatedAt.AsTime(),
			events[len(events)-1].CreatedAt.AsTime().Add(-time.Millisecond))
		t.Logf("listEventsUniqID, err: %+v, %v", listEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listEventsUniqID, 1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(events[len(events)-1], listEventsUniqID[0]) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", events[len(events)-1],
				listEventsUniqID[0])
		}
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		event := &api.Event{OrgId: createOrg.Id, RuleId: uuid.NewString(),
			UniqId: "dao-event-" + random.String(16),
			CreatedAt: timestamppb.New(time.Now().UTC().Truncate(
				time.Millisecond)), TraceId: uuid.NewString()}

		err = globalEventDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		listEvents, err := globalEventDAO.List(ctx, uuid.NewString(),
			event.UniqId, "", "", event.CreatedAt.AsTime(),
			event.CreatedAt.AsTime().Add(-time.Millisecond))
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.NoError(t, err)
		require.Len(t, listEvents, 0)
	})

	t.Run("List events by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listEvents, err := globalEventDAO.List(ctx, random.String(10),
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

		for i := 0; i < 5; i++ {
			event := &api.Event{OrgId: createOrg.Id, RuleId: uuid.NewString(),
				UniqId: "dao-event-" + random.String(16),
				CreatedAt: timestamppb.New(time.Now().UTC().Truncate(
					time.Millisecond)), TraceId: uuid.NewString()}
			events = append(events, event)

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			err := globalEventDAO.Create(ctx, event)
			t.Logf("err: %v", err)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		sort.Slice(events, func(i, j int) bool {
			return events[i].CreatedAt.AsTime().After(
				events[j].CreatedAt.AsTime())
		})

		ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		// Verify results.
		latEventsUniqID, err := globalEventDAO.Latest(ctx, createOrg.Id, "")
		t.Logf("latEventsUniqID, err: %+v, %v", latEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, latEventsUniqID, len(events))

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		for i, event := range events {
			if !proto.Equal(event, latEventsUniqID[i]) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", event,
					latEventsUniqID[i])
			}
		}

		// Verify results by rule ID.
		latEventsUniqID, err = globalEventDAO.Latest(ctx, createOrg.Id,
			events[len(events)-1].RuleId)
		t.Logf("latEventsUniqID, err: %+v, %v", latEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, latEventsUniqID, 1)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(events[len(events)-1], latEventsUniqID[0]) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", events[len(events)-1],
				latEventsUniqID[0])
		}
	})

	t.Run("Latest events are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		event := &api.Event{OrgId: createOrg.Id, RuleId: uuid.NewString(),
			UniqId: "dao-event-" + random.String(16),
			CreatedAt: timestamppb.New(time.Now().UTC().Truncate(
				time.Millisecond)), TraceId: uuid.NewString()}

		err = globalEventDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		latEvents, err := globalEventDAO.Latest(ctx, uuid.NewString(), "")
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.NoError(t, err)
		require.Len(t, latEvents, 0)
	})

	t.Run("Latest events by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		latEvents, err := globalEventDAO.Latest(ctx, random.String(10), "")
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.Nil(t, latEvents)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
