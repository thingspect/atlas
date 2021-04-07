// +build !unit

package user

import (
	"context"
	"reflect"
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

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("dao-user", createOrg.Id)
		createUser, _ := proto.Clone(user).(*api.User)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, createUser)
		t.Logf("user, createUser, err: %+v, %+v, %v", user, createUser, err)
		require.NoError(t, err)
		require.NotEqual(t, user.Id, createUser.Id)
		require.WithinDuration(t, time.Now(), createUser.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createUser.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("dao-user", createOrg.Id)
		user.Email = "dao-user-" + random.String(80)
		createUser, _ := proto.Clone(user).(*api.User)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, createUser)
		t.Logf("user, createUser, err: %+v, %+v, %v", user, createUser, err)
		require.Nil(t, createUser)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
		createOrg.Id))
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	t.Run("Read user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, createUser.Id,
			createUser.OrgId)
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.NoError(t, err)
		require.Equal(t, createUser, readUser)
	})

	t.Run("Read user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, uuid.NewString(),
			uuid.NewString())
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.Nil(t, readUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, createUser.Id,
			uuid.NewString())
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.Nil(t, readUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read user by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, random.String(10),
			createUser.OrgId)
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.Nil(t, readUser)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestReadByEmail(t *testing.T) {
	t.Parallel()

	hash := random.Bytes(60)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
		createOrg.Id))
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	err = globalUserDAO.UpdatePassword(ctx, createUser.Id, createOrg.Id,
		hash)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	t.Run("Read user by valid email", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, readHash, err := globalUserDAO.ReadByEmail(ctx,
			createUser.Email, createOrg.Name)
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.NoError(t, err)
		require.Equal(t, hash, readHash)

		// Normalize timestamps.
		require.True(t, readUser.UpdatedAt.AsTime().After(
			createUser.CreatedAt.AsTime()))
		require.WithinDuration(t, readUser.UpdatedAt.AsTime(),
			createUser.UpdatedAt.AsTime(), 2*time.Second)
		createUser.UpdatedAt = readUser.UpdatedAt

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(createUser, readUser) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", createUser, readUser)
		}
	})

	t.Run("Read user by unknown email", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, readHash, err := globalUserDAO.ReadByEmail(ctx,
			random.Email(), random.String(10))
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.Nil(t, readUser)
		require.Nil(t, readHash)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org name", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, readHash, err := globalUserDAO.ReadByEmail(ctx,
			createUser.Email, random.String(10))
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.Nil(t, readUser)
		require.Nil(t, readHash)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update user by valid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "dao-user-" + random.Email()
		createUser.Phone = "+15125553434"
		createUser.Role = common.Role_ADMIN
		createUser.Status = common.Status_DISABLED
		createUser.Tags = nil
		createUser.AppKey = random.String(30)
		updateUser, _ := proto.Clone(createUser).(*api.User)

		updateUser, err = globalUserDAO.Update(ctx, updateUser)
		t.Logf("createUser, updateUser, err: %+v, %+v, %v", createUser,
			updateUser, err)
		require.NoError(t, err)
		require.Equal(t, createUser.Email, updateUser.Email)
		require.Equal(t, createUser.Phone, updateUser.Phone)
		require.Equal(t, createUser.Role, updateUser.Role)
		require.Equal(t, createUser.Status, updateUser.Status)
		require.Equal(t, createUser.Tags, updateUser.Tags)
		require.Equal(t, createUser.AppKey, updateUser.AppKey)
		require.True(t, updateUser.UpdatedAt.AsTime().After(
			updateUser.CreatedAt.AsTime()))
		require.WithinDuration(t, createUser.CreatedAt.AsTime(),
			updateUser.UpdatedAt.AsTime(), 2*time.Second)

		readUser, err := globalUserDAO.Read(ctx, createUser.Id,
			createUser.OrgId)
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.NoError(t, err)
		require.Equal(t, updateUser, readUser)
	})

	t.Run("Update unknown user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateUser, err := globalUserDAO.Update(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.OrgId = uuid.NewString()
		createUser.Email = "dao-user-" + random.Email()
		updateUser, _ := proto.Clone(createUser).(*api.User)

		updateUser, err = globalUserDAO.Update(ctx, updateUser)
		t.Logf("createUser, updateUser, err: %+v, %+v, %v", createUser,
			updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update user by invalid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "dao-user-" + random.String(80)
		createUser.Status = common.Status_DISABLED
		updateUser, _ := proto.Clone(createUser).(*api.User)

		updateUser, err = globalUserDAO.Update(ctx, updateUser)
		t.Logf("createUser, updateUser, err: %+v, %+v, %v", createUser,
			updateUser, err)
		require.Nil(t, updateUser)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Parallel()

	hash := random.Bytes(60)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update user password by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.UpdatePassword(ctx, createUser.Id, createOrg.Id,
			hash)
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Update user password by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalUserDAO.UpdatePassword(ctx, uuid.NewString(),
			createOrg.Id, hash)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Password updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.UpdatePassword(ctx, createUser.Id,
			uuid.NewString(), hash)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.Delete(ctx, createUser.Id, createOrg.Id)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read user by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readUser, err := globalUserDAO.Read(ctx, createUser.Id,
				createOrg.Id)
			t.Logf("readUser, err: %+v, %v", readUser, err)
			require.Nil(t, readUser)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalUserDAO.Delete(ctx, uuid.NewString(), createOrg.Id)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.Delete(ctx, createUser.Id, uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	userIDs := []string{}
	userRoles := []common.Role{}
	userTags := [][]string{}
	userAppKeys := []string{}
	userTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.Id))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		userIDs = append(userIDs, createUser.Id)
		userRoles = append(userRoles, createUser.Role)
		userTags = append(userTags, createUser.Tags)
		userAppKeys = append(userAppKeys, createUser.AppKey)
		userTSes = append(userTSes, createUser.CreatedAt.AsTime())
	}

	t.Run("List users by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 0, "")
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, user := range listUsers {
			if user.Id == userIDs[len(userIDs)-1] &&
				user.Role == userRoles[len(userRoles)-1] &&
				reflect.DeepEqual(user.Tags, userTags[len(userTags)-1]) &&
				user.AppKey == userAppKeys[len(userAppKeys)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.Id,
			userTSes[0], userIDs[0], 5, "")
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, user := range listUsers {
			if user.Id == userIDs[len(userIDs)-1] &&
				user.Role == userRoles[len(userRoles)-1] &&
				reflect.DeepEqual(user.Tags, userTags[len(userTags)-1]) &&
				user.AppKey == userAppKeys[len(userAppKeys)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 1, "")
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, int32(3), listCount)
	})

	t.Run("List users with tag filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.Id,
			time.Time{}, "", 5, userTags[2][0])
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, userIDs[len(userIDs)-1], listUsers[0].Id)
		require.Equal(t, userRoles[len(userRoles)-1], listUsers[0].Role)
		require.Equal(t, userTags[len(userTags)-1], listUsers[0].Tags)
		require.Equal(t, userAppKeys[len(userAppKeys)-1], listUsers[0].AppKey)
	})

	t.Run("List users with tag filter and pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.Id,
			userTSes[0], userIDs[0], 5, userTags[2][0])
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, userIDs[len(userIDs)-1], listUsers[0].Id)
		require.Equal(t, userRoles[len(userRoles)-1], listUsers[0].Role)
		require.Equal(t, userTags[len(userTags)-1], listUsers[0].Tags)
		require.Equal(t, userAppKeys[len(userAppKeys)-1], listUsers[0].AppKey)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx,
			uuid.NewString(), time.Time{}, "", 0, "")
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("List users by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, random.String(10),
			time.Time{}, "", 0, "")
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.Nil(t, listUsers)
		require.Equal(t, int32(0), listCount)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestListByTags(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-user"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	userIDs := []string{}
	userTags := [][]string{}
	for i := 0; i < 3; i++ {
		user := random.User("dao-user", createOrg.Id)
		user.Status = common.Status_ACTIVE
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		userIDs = append(userIDs, createUser.Id)
		userTags = append(userTags, createUser.Tags)
	}

	t.Run("List users by valid org ID and unique tag", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, err := globalUserDAO.ListByTags(ctx, createOrg.Id,
			userTags[len(userTags)-1])
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, listUsers[0].Id, userIDs[len(userIDs)-1])
		require.Equal(t, listUsers[0].Tags, userTags[len(userTags)-1])
	})

	t.Run("List users by valid org ID and common tag", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		user := random.User("dao-user", createOrg.Id)
		user.Status = common.Status_ACTIVE
		user.Tags = []string{userTags[0][0]}

		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		listUsers, err := globalUserDAO.ListByTags(ctx, createOrg.Id,
			createUser.Tags)
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 2)
		for _, user := range listUsers {
			require.Contains(t, user.Tags, createUser.Tags[0])
		}
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, err := globalUserDAO.ListByTags(ctx, uuid.NewString(),
			userTags[len(userTags)-1])
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 0)
	})

	t.Run("List users by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, err := globalUserDAO.ListByTags(ctx, random.String(10),
			userTags[len(userTags)-1])
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
