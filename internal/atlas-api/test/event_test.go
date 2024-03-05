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

func TestListEvents(t *testing.T) {
	t.Parallel()

	t.Run("List events by UniqID, dev ID, and rule ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-event", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		events := []*api.Event{}

		for range 5 {
			event := random.Event("api-event", globalAdminOrgID)
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
		evCli := api.NewEventServiceClient(globalAdminGRPCConn)
		listEventsUniqID, err := evCli.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_UniqId{UniqId: createDev.GetUniqId()},
			EndTime: events[0].GetCreatedAt(), StartTime: timestamppb.New(
				events[len(events)-1].GetCreatedAt().AsTime().Add(
					-time.Millisecond)),
		})
		t.Logf("listEventsUniqID, err: %+v, %v", listEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listEventsUniqID.GetEvents(), len(events))
		require.EqualExportedValues(t, &api.ListEventsResponse{Events: events},
			listEventsUniqID)

		// Verify results by dev ID without oldest event.
		listEventsDevID, err := evCli.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof:   &api.ListEventsRequest_DeviceId{DeviceId: createDev.GetId()},
			StartTime: events[len(events)-1].GetCreatedAt(),
		})
		t.Logf("listEventsDevID, err: %+v, %v", listEventsDevID, err)
		require.NoError(t, err)
		require.Len(t, listEventsDevID.GetEvents(), len(events)-1)
		require.EqualExportedValues(t, &api.ListEventsResponse{
			Events: events[:len(events)-1],
		}, listEventsDevID)

		// Verify results by UniqID and rule ID.
		listEventsUniqID, err = evCli.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_UniqId{UniqId: createDev.GetUniqId()},
			RuleId:  events[len(events)-1].GetRuleId(),
			EndTime: events[0].GetCreatedAt(),
			StartTime: timestamppb.New(events[len(events)-1].GetCreatedAt().
				AsTime().Add(-time.Millisecond)),
		})
		t.Logf("listEventsUniqID, err: %+v, %v", listEventsUniqID, err)
		require.NoError(t, err)
		require.Len(t, listEventsUniqID.GetEvents(), 1)
		require.EqualExportedValues(t, &api.ListEventsResponse{
			Events: []*api.Event{events[len(events)-1]},
		}, listEventsUniqID)
	})

	t.Run("List events are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		event := random.Event("api-event", createOrg.GetId())

		err = globalEvDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		evCli := api.NewEventServiceClient(globalAdminGRPCConn)
		listEvents, err := evCli.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_UniqId{UniqId: event.GetUniqId()},
		})
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.NoError(t, err)
		require.Empty(t, listEvents.GetEvents())
	})

	t.Run("List events by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		evCli := api.NewEventServiceClient(globalAdminGRPCConn)
		listEvents, err := evCli.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_DeviceId{
				DeviceId: uuid.NewString(),
			}, EndTime: timestamppb.Now(), StartTime: timestamppb.New(
				time.Now().Add(-91 * 24 * time.Hour)),
		})
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.Nil(t, listEvents)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"maximum time range exceeded")
	})

	t.Run("List events by invalid dev ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		evCli := api.NewEventServiceClient(globalAdminGRPCConn)
		listEvents, err := evCli.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_DeviceId{
				DeviceId: random.String(10),
			},
		})
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.Nil(t, listEvents)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid ListEventsRequest.DeviceId: value must be a valid UUID | "+
			"caused by: invalid uuid format")
	})
}

func TestLatestEvents(t *testing.T) {
	t.Parallel()

	t.Run("Latest events", func(t *testing.T) {
		t.Parallel()

		events := []*api.Event{}

		for range 5 {
			event := random.Event("api-event", globalAdminOrgID)
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		// Verify results.
		evCli := api.NewEventServiceClient(globalAdminGRPCConn)
		latEvents, err := evCli.LatestEvents(ctx, &api.LatestEventsRequest{})
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(latEvents.GetEvents()), 5)

		var found bool
		for _, event := range latEvents.GetEvents() {
			if event.GetRuleId() == events[len(events)-1].GetRuleId() &&
				event.GetUniqId() == events[len(events)-1].GetUniqId() {
				found = true
			}
		}
		require.True(t, found)

		// Verify results by rule ID.
		latEventsRuleID, err := evCli.LatestEvents(ctx,
			&api.LatestEventsRequest{RuleId: events[len(events)-1].GetRuleId()})
		t.Logf("latEventsDevID, err: %+v, %v", latEventsRuleID, err)
		require.NoError(t, err)
		require.Len(t, latEventsRuleID.GetEvents(), 1)
		require.EqualExportedValues(t, &api.LatestEventsResponse{
			Events: []*api.Event{events[len(events)-1]},
		}, latEventsRuleID)
	})

	t.Run("Latest events are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-event"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		event := random.Event("api-event", createOrg.GetId())

		err = globalEvDAO.Create(ctx, event)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		evCli := api.NewEventServiceClient(secondaryAdminGRPCConn)
		latEvents, err := evCli.LatestEvents(ctx, &api.LatestEventsRequest{})
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.NoError(t, err)
		require.Empty(t, latEvents.GetEvents())
	})

	t.Run("Latest events by invalid rule ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		evCli := api.NewEventServiceClient(globalAdminGRPCConn)
		latEvents, err := evCli.LatestEvents(ctx, &api.LatestEventsRequest{
			RuleId: random.String(10),
		})
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.Nil(t, latEvents)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid LatestEventsRequest.RuleId: value must be a valid UUID | "+
			"caused by: invalid uuid format")
	})
}
