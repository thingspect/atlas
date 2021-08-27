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

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-tag"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	var tCount int
	var lastTag string
	for i := 0; i < 3; i++ {
		createDev, err := globalDevDAO.Create(ctx, random.Device("dao-tag",
			createOrg.Id))
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		tCount += len(createDev.Tags)
		lastTag = createDev.Tags[len(createDev.Tags)-1]
	}

	for i := 0; i < 3; i++ {
		createUser, err := globalUserDAO.Create(ctx, random.User("dao-tag",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		tCount += len(createUser.Tags)
		lastTag = createUser.Tags[len(createUser.Tags)-1]
	}

	t.Run("List tags by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listTags, err := globalTagDAO.List(ctx, createOrg.Id)
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dev := random.Device("dao-tag", createOrg.Id)
		dev.Tags = []string{lastTag}

		createDev, err := globalDevDAO.Create(ctx, dev)
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		user := random.User("dao-tag", createOrg.Id)
		user.Tags = []string{lastTag}

		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		listTags, err := globalTagDAO.List(ctx, createOrg.Id)
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

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listTags, err := globalTagDAO.List(ctx, uuid.NewString())
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.NoError(t, err)
		require.Len(t, listTags, 0)
	})

	t.Run("List tags by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listTags, err := globalTagDAO.List(ctx, random.String(10))
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.Nil(t, listTags)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
