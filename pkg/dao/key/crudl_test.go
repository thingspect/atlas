// +build !unit

package key

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 8 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-key"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid key", func(t *testing.T) {
		t.Parallel()

		key := random.Key("dao-key", createOrg.Id)
		createKey, _ := proto.Clone(key).(*api.Key)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createKey, err := globalKeyDAO.Create(ctx, createKey)
		t.Logf("key, createKey, err: %+v, %+v, %v", key, createKey, err)
		require.NoError(t, err)
		require.NotEqual(t, key.Id, createKey.Id)
		require.WithinDuration(t, time.Now(), createKey.CreatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid key", func(t *testing.T) {
		t.Parallel()

		key := random.Key("dao-key", createOrg.Id)
		key.Name = "dao-key-" + random.String(80)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createKey, err := globalKeyDAO.Create(ctx, key)
		t.Logf("key, createKey, err: %+v, %+v, %v", key, createKey, err)
		require.Nil(t, createKey)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-key"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createKey, err := globalKeyDAO.Create(ctx, random.Key("dao-key",
		createOrg.Id))
	t.Logf("createKey, err: %+v, %v", createKey, err)
	require.NoError(t, err)

	t.Run("Read key by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readKey, err := globalKeyDAO.read(ctx, createKey.Id, createKey.OrgId)
		t.Logf("readKey, err: %+v, %v", readKey, err)
		require.NoError(t, err)
		require.Equal(t, createKey, readKey)
	})

	t.Run("Read key by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readKey, err := globalKeyDAO.read(ctx, uuid.NewString(),
			uuid.NewString())
		t.Logf("readKey, err: %+v, %v", readKey, err)
		require.Nil(t, readKey)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readKey, err := globalKeyDAO.read(ctx, createKey.Id,
			uuid.NewString())
		t.Logf("readKey, err: %+v, %v", readKey, err)
		require.Nil(t, readKey)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read key by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readKey, err := globalKeyDAO.read(ctx, random.String(10),
			createKey.OrgId)
		t.Logf("readKey, err: %+v, %v", readKey, err)
		require.Nil(t, readKey)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-key"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete key by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createKey, err := globalKeyDAO.Create(ctx, random.Key("dao-key",
			createOrg.Id))
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)

		err = globalKeyDAO.Delete(ctx, createKey.Id, createOrg.Id)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read key by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readKey, err := globalKeyDAO.read(ctx, createKey.Id,
				createOrg.Id)
			t.Logf("readKey, err: %+v, %v", readKey, err)
			require.Nil(t, readKey)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete key by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalKeyDAO.Delete(ctx, uuid.NewString(), createOrg.Id)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createKey, err := globalKeyDAO.Create(ctx, random.Key("dao-key",
			createOrg.Id))
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)

		err = globalKeyDAO.Delete(ctx, createKey.Id, uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-key"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	keyIDs := []string{}
	keyNames := []string{}
	keyRoles := []common.Role{}
	keyTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createKey, err := globalKeyDAO.Create(ctx, random.Key("dao-key",
			createOrg.Id))
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)

		keyIDs = append(keyIDs, createKey.Id)
		keyNames = append(keyNames, createKey.Name)
		keyRoles = append(keyRoles, createKey.Role)
		keyTSes = append(keyTSes, createKey.CreatedAt.AsTime())
	}

	t.Run("List keys by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listKeys, listCount, err := globalKeyDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 0, "")
		t.Logf("listKeys, listCount, err: %+v, %v, %v", listKeys, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listKeys, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, key := range listKeys {
			if key.Id == keyIDs[len(keyIDs)-1] &&
				key.Name == keyNames[len(keyNames)-1] &&
				key.Role == keyRoles[len(keyRoles)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List keys by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listKeys, listCount, err := globalKeyDAO.List(ctx, createOrg.Id,
			keyTSes[0], keyIDs[0], 5, "")
		t.Logf("listKeys, listCount, err: %+v, %v, %v", listKeys, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listKeys, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, key := range listKeys {
			if key.Id == keyIDs[len(keyIDs)-1] &&
				key.Name == keyNames[len(keyNames)-1] &&
				key.Role == keyRoles[len(keyRoles)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List keys by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listKeys, listCount, err := globalKeyDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 1, "")
		t.Logf("listKeys, listCount, err: %+v, %v, %v", listKeys, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listKeys, 1)
		require.Equal(t, int32(3), listCount)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listKeys, listCount, err := globalKeyDAO.List(ctx, uuid.NewString(),
			time.Time{}, "", 0, "")
		t.Logf("listKeys, listCount, err: %+v, %v, %v", listKeys, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listKeys, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("List keys by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listKeys, listCount, err := globalKeyDAO.List(ctx, random.String(10),
			time.Time{}, "", 0, "")
		t.Logf("listKeys, listCount, err: %+v, %v, %v", listKeys, listCount,
			err)
		require.Nil(t, listKeys)
		require.Equal(t, int32(0), listCount)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
