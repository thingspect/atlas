// +build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListEvents(t *testing.T) {
	t.Parallel()

	t.Run("List events by valid UniqID with ts", func(t *testing.T) {
		t.Parallel()

		event := random.Event("dao-event", uuid.NewString())
		retEvent, _ := proto.Clone(event).(*api.Event)
		end := time.Now().UTC()
		start := time.Now().UTC().Add(-15 * time.Minute)

		eventer := NewMockEventer(gomock.NewController(t))
		eventer.EXPECT().List(gomock.Any(), event.OrgId, event.UniqId, "", "",
			end, start).Return([]*api.Event{retEvent}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: event.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		evSvc := NewEvent(eventer)
		listEvents, err := evSvc.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_UniqId{UniqId: event.UniqId},
			EndTime: timestamppb.New(end), StartTime: timestamppb.New(start)})
		t.Logf("event, listEvents, err: %+v, %+v, %v", event, listEvents, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListEventsResponse{
			Events: []*api.Event{event}}, listEvents) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListEventsResponse{
					Events: []*api.Event{event}}, listEvents)
		}
	})

	t.Run("List events by valid dev ID with rule ID", func(t *testing.T) {
		t.Parallel()

		event := random.Event("dao-event", uuid.NewString())
		retEvent, _ := proto.Clone(event).(*api.Event)
		devID := uuid.NewString()

		eventer := NewMockEventer(gomock.NewController(t))
		eventer.EXPECT().List(gomock.Any(), event.OrgId, "", devID,
			event.RuleId, matcher.NewRecentMatcher(2*time.Second),
			matcher.NewRecentMatcher(24*time.Hour+2*time.Second)).
			Return([]*api.Event{retEvent}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: event.OrgId,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		evSvc := NewEvent(eventer)
		listEvents, err := evSvc.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_DeviceId{DeviceId: devID},
			RuleId:  event.RuleId})
		t.Logf("event, listEvents, err: %+v, %+v, %v", event, listEvents, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListEventsResponse{
			Events: []*api.Event{event}}, listEvents) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListEventsResponse{
					Events: []*api.Event{event}}, listEvents)
		}
	})

	t.Run("List events with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		evSvc := NewEvent(nil)
		listEvents, err := evSvc.ListEvents(ctx, &api.ListEventsRequest{})
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.Nil(t, listEvents)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("List events with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_CONTACT}), testTimeout)
		defer cancel()

		evSvc := NewEvent(nil)
		listEvents, err := evSvc.ListEvents(ctx, &api.ListEventsRequest{})
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.Nil(t, listEvents)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("List events by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		evSvc := NewEvent(nil)
		listEvents, err := evSvc.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_UniqId{UniqId: "api-event-" +
				random.String(16)}, EndTime: timestamppb.Now(),
			StartTime: timestamppb.New(time.Now().Add(-91 * 24 * time.Hour))})
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.Nil(t, listEvents)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"maximum time range exceeded"), err)
	})

	t.Run("List events by invalid org ID", func(t *testing.T) {
		t.Parallel()

		eventer := NewMockEventer(gomock.NewController(t))
		eventer.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: "aaa",
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		evSvc := NewEvent(eventer)
		listEvents, err := evSvc.ListEvents(ctx, &api.ListEventsRequest{
			IdOneof: &api.ListEventsRequest_UniqId{UniqId: "api-event-" +
				random.String(16)}})
		t.Logf("listEvents, err: %+v, %v", listEvents, err)
		require.Nil(t, listEvents)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestLatestEvents(t *testing.T) {
	t.Parallel()

	t.Run("Latest events by valid rule ID", func(t *testing.T) {
		t.Parallel()

		event := random.Event("dao-event", uuid.NewString())
		retEvent, _ := proto.Clone(event).(*api.Event)
		orgID := uuid.NewString()

		eventer := NewMockEventer(gomock.NewController(t))
		eventer.EXPECT().Latest(gomock.Any(), orgID, event.RuleId).
			Return([]*api.Event{retEvent}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		evSvc := NewEvent(eventer)
		latEvents, err := evSvc.LatestEvents(ctx, &api.LatestEventsRequest{
			RuleId: event.RuleId})
		t.Logf("event, latEvents, err: %+v, %+v, %v", event, latEvents, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.LatestEventsResponse{Events: []*api.Event{event}},
			latEvents) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.LatestEventsResponse{
				Events: []*api.Event{event}}, latEvents)
		}
	})

	t.Run("Latest events with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		evSvc := NewEvent(nil)
		latEvents, err := evSvc.LatestEvents(ctx, &api.LatestEventsRequest{})
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.Nil(t, latEvents)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("Latest events with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_CONTACT}), testTimeout)
		defer cancel()

		evSvc := NewEvent(nil)
		latEvents, err := evSvc.LatestEvents(ctx, &api.LatestEventsRequest{})
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.Nil(t, latEvents)
		require.Equal(t, errPerm(common.Role_VIEWER), err)
	})

	t.Run("Latest events by invalid org ID", func(t *testing.T) {
		t.Parallel()

		eventer := NewMockEventer(gomock.NewController(t))
		eventer.EXPECT().Latest(gomock.Any(), "aaa", gomock.Any()).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: "aaa",
				Role: common.Role_ADMIN}), testTimeout)
		defer cancel()

		evSvc := NewEvent(eventer)
		latEvents, err := evSvc.LatestEvents(ctx, &api.LatestEventsRequest{})
		t.Logf("latEvents, err: %+v, %v", latEvents, err)
		require.Nil(t, latEvents)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}
