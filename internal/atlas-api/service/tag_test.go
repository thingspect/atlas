//go:build !integration

package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestListTags(t *testing.T) {
	t.Parallel()

	t.Run("List tags by valid org ID", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		tags := random.Tags("api-tag", 5)

		tagger := NewMockTagger(gomock.NewController(t))
		tagger.EXPECT().List(gomock.Any(), orgID).Return(tags, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		tagSvc := NewTag(tagger)
		listTags, err := tagSvc.ListTags(ctx, &api.ListTagsRequest{})
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListTagsResponse{Tags: tags}, listTags) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListTagsResponse{Tags: tags}, listTags)
		}
	})

	t.Run("List tags with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		tagSvc := NewTag(nil)
		listTags, err := tagSvc.ListTags(ctx, &api.ListTagsRequest{})
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.Nil(t, listTags)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List tags with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		tagSvc := NewTag(nil)
		listTags, err := tagSvc.ListTags(ctx, &api.ListTagsRequest{})
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.Nil(t, listTags)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List tags by invalid org ID", func(t *testing.T) {
		t.Parallel()

		tagger := NewMockTagger(gomock.NewController(t))
		tagger.EXPECT().List(gomock.Any(), "aaa").Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		tagSvc := NewTag(tagger)
		listTags, err := tagSvc.ListTags(ctx, &api.ListTagsRequest{})
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.Nil(t, listTags)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}
