//go:build !unit

package tag

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
)

const testTimeout = 14 * time.Second

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-tag"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	var tCount int
	var lastTag string
	for range 3 {
		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-tag",
			createOrg.GetId()))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		tCount += len(createDev.GetTags())
		lastTag = createDev.GetTags()[len(createDev.GetTags())-1]
	}

	for range 3 {
		createUser, err := globalUserDAO.Create(ctx, random.User("dao-tag",
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		tCount += len(createUser.GetTags())
		lastTag = createUser.GetTags()[len(createUser.GetTags())-1]
	}

	t.Run("List tags by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		listTags, err := globalTagDAO.List(ctx, createOrg.GetId())
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.NoError(t, err)
		require.Len(t, listTags, tCount)

		var found bool
		for _, tag := range listTags {
			if tag == lastTag {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List tags by valid org ID without duplicates", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		dev := random.Device("dao-tag", createOrg.GetId())
		dev.Tags = []string{lastTag}

		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		user := random.User("dao-tag", createOrg.GetId())
		user.Tags = []string{lastTag}

		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		listTags, err := globalTagDAO.List(ctx, createOrg.GetId())
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.NoError(t, err)
		require.Len(t, listTags, tCount)

		var found bool
		for _, tag := range listTags {
			if tag == lastTag {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		listTags, err := globalTagDAO.List(ctx, uuid.NewString())
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.NoError(t, err)
		require.Empty(t, listTags)
	})

	t.Run("List tags by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		listTags, err := globalTagDAO.List(ctx, random.String(10))
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.Nil(t, listTags)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
