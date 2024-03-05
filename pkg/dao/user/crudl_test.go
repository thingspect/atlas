//go:build !unit

package user

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
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

		user := random.User("dao-user", createOrg.GetId())
		createUser, _ := proto.Clone(user).(*api.User)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, createUser)
		t.Logf("user, createUser, err: %+v, %+v, %v", user, createUser, err)
		require.NoError(t, err)
		require.NotEqual(t, user.GetId(), createUser.GetId())
		require.WithinDuration(t, time.Now(), createUser.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createUser.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("dao-user", createOrg.GetId())
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
		createOrg.GetId()))
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	t.Run("Read user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, createUser.GetId(),
			createUser.GetOrgId())
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

		readUser, err := globalUserDAO.Read(ctx, createUser.GetId(),
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
			createUser.GetOrgId())
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
		createOrg.GetId()))
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	err = globalUserDAO.UpdatePassword(ctx, createUser.GetId(), createOrg.GetId(),
		hash)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	t.Run("Read user by valid email", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readUser, readHash, err := globalUserDAO.ReadByEmail(ctx,
			createUser.GetEmail(), createOrg.GetName())
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.NoError(t, err)
		require.Equal(t, hash, readHash)

		// Normalize timestamp.
		require.True(t, readUser.GetUpdatedAt().AsTime().After(
			createUser.GetCreatedAt().AsTime()))
		require.WithinDuration(t, readUser.GetUpdatedAt().AsTime(),
			createUser.GetUpdatedAt().AsTime(), 2*time.Second)
		createUser.UpdatedAt = readUser.GetUpdatedAt()
		require.EqualExportedValues(t, createUser, readUser)
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
			createUser.GetEmail(), random.String(10))
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
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Name = "dao-user-" + random.String(10)
		createUser.Email = "dao-user-" + random.Email()
		createUser.Phone = "+15125550101"
		createUser.Role = api.Role_ADMIN
		createUser.Status = api.Status_DISABLED
		createUser.Tags = nil
		createUser.AppKey = random.String(30)
		updateUser, _ := proto.Clone(createUser).(*api.User)

		updateUser, err = globalUserDAO.Update(ctx, updateUser)
		t.Logf("createUser, updateUser, err: %+v, %+v, %v", createUser,
			updateUser, err)
		require.NoError(t, err)
		require.Equal(t, createUser.GetName(), updateUser.GetName())
		require.Equal(t, createUser.GetEmail(), updateUser.GetEmail())
		require.Equal(t, createUser.GetPhone(), updateUser.GetPhone())
		require.Equal(t, createUser.GetRole(), updateUser.GetRole())
		require.Equal(t, createUser.GetStatus(), updateUser.GetStatus())
		require.Equal(t, createUser.GetTags(), updateUser.GetTags())
		require.Equal(t, createUser.GetAppKey(), updateUser.GetAppKey())
		require.True(t, updateUser.GetUpdatedAt().AsTime().After(
			updateUser.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createUser.GetCreatedAt().AsTime(),
			updateUser.GetUpdatedAt().AsTime(), 2*time.Second)

		readUser, err := globalUserDAO.Read(ctx, createUser.GetId(),
			createUser.GetOrgId())
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.NoError(t, err)
		require.Equal(t, updateUser, readUser)
	})

	t.Run("Update unknown user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateUser, err := globalUserDAO.Update(ctx, random.User("dao-user",
			createOrg.GetId()))
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.GetId()))
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
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "dao-user-" + random.String(80)
		createUser.Status = api.Status_DISABLED
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
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.UpdatePassword(ctx, createUser.GetId(), createOrg.GetId(),
			hash)
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Update user password by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalUserDAO.UpdatePassword(ctx, uuid.NewString(),
			createOrg.GetId(), hash)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Password updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.UpdatePassword(ctx, createUser.GetId(),
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
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.Delete(ctx, createUser.GetId(), createOrg.GetId())
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read user by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readUser, err := globalUserDAO.Read(ctx, createUser.GetId(),
				createOrg.GetId())
			t.Logf("readUser, err: %+v, %v", readUser, err)
			require.Nil(t, readUser)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalUserDAO.Delete(ctx, uuid.NewString(), createOrg.GetId())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.Delete(ctx, createUser.GetId(), uuid.NewString())
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
	userNames := []string{}
	userRoles := []api.Role{}
	userTags := [][]string{}
	userAppKeys := []string{}
	userTSes := []time.Time{}
	for range 3 {
		createUser, err := globalUserDAO.Create(ctx, random.User("dao-user",
			createOrg.GetId()))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		userIDs = append(userIDs, createUser.GetId())
		userNames = append(userNames, createUser.GetName())
		userRoles = append(userRoles, createUser.GetRole())
		userTags = append(userTags, createUser.GetTags())
		userAppKeys = append(userAppKeys, createUser.GetAppKey())
		userTSes = append(userTSes, createUser.GetCreatedAt().AsTime())
	}

	t.Run("List users by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.GetId(),
			time.Time{}, "", 0, "")
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, user := range listUsers {
			if user.GetId() == userIDs[len(userIDs)-1] &&
				user.GetName() == userNames[len(userNames)-1] &&
				user.GetRole() == userRoles[len(userRoles)-1] &&
				reflect.DeepEqual(user.GetTags(), userTags[len(userTags)-1]) &&
				user.GetAppKey() == userAppKeys[len(userAppKeys)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.GetId(),
			userTSes[0], userIDs[0], 5, "")
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, user := range listUsers {
			if user.GetId() == userIDs[len(userIDs)-1] &&
				user.GetName() == userNames[len(userNames)-1] &&
				user.GetRole() == userRoles[len(userRoles)-1] &&
				reflect.DeepEqual(user.GetTags(), userTags[len(userTags)-1]) &&
				user.GetAppKey() == userAppKeys[len(userAppKeys)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.GetId(),
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

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.GetId(),
			time.Time{}, "", 5, userTags[2][0])
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, userIDs[len(userIDs)-1], listUsers[0].GetId())
		require.Equal(t, userNames[len(userNames)-1], listUsers[0].GetName())
		require.Equal(t, userRoles[len(userRoles)-1], listUsers[0].GetRole())
		require.Equal(t, userTags[len(userTags)-1], listUsers[0].GetTags())
		require.Equal(t, userAppKeys[len(userAppKeys)-1], listUsers[0].GetAppKey())
	})

	t.Run("List users with tag filter and pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.GetId(),
			userTSes[0], userIDs[0], 5, userTags[2][0])
		t.Logf("listUsers, listCount, err: %+v, %v, %v", listUsers, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, int32(1), listCount)

		require.Equal(t, userIDs[len(userIDs)-1], listUsers[0].GetId())
		require.Equal(t, userNames[len(userNames)-1], listUsers[0].GetName())
		require.Equal(t, userRoles[len(userRoles)-1], listUsers[0].GetRole())
		require.Equal(t, userTags[len(userTags)-1], listUsers[0].GetTags())
		require.Equal(t, userAppKeys[len(userAppKeys)-1], listUsers[0].GetAppKey())
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
		require.Empty(t, listUsers)
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
	userNames := []string{}
	userTags := [][]string{}
	for range 3 {
		user := random.User("dao-user", createOrg.GetId())
		user.Status = api.Status_ACTIVE
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		userIDs = append(userIDs, createUser.GetId())
		userNames = append(userNames, createUser.GetName())
		userTags = append(userTags, createUser.GetTags())
	}

	t.Run("List users by valid org ID and unique tag", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listUsers, err := globalUserDAO.ListByTags(ctx, createOrg.GetId(),
			userTags[len(userTags)-1])
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, listUsers[0].GetId(), userIDs[len(userIDs)-1])
		require.Equal(t, listUsers[0].GetName(), userNames[len(userNames)-1])
		require.Equal(t, listUsers[0].GetTags(), userTags[len(userTags)-1])
	})

	t.Run("List users by valid org ID and api tag", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		user := random.User("dao-user", createOrg.GetId())
		user.Status = api.Status_ACTIVE
		user.Tags = []string{userTags[0][0]}

		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		listUsers, err := globalUserDAO.ListByTags(ctx, createOrg.GetId(),
			createUser.GetTags())
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 2)
		for _, user := range listUsers {
			require.Contains(t, user.GetTags(), createUser.GetTags()[0])
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
		require.Empty(t, listUsers)
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
