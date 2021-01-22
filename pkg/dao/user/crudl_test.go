// +build !unit

package user

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Create valid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)
		require.Equal(t, user.OrgId, createUser.OrgId)
		require.Equal(t, user.Email, createUser.Email)
		require.Equal(t, user.Status, createUser.Status)
		require.WithinDuration(t, time.Now(), createUser.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createUser.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.String(80), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" + random.Email(),
		Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	t.Run("Read user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, createUser.Id,
			createUser.OrgId)
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.NoError(t, err)
		require.Equal(t, createUser, readUser)
	})

	t.Run("Read user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, uuid.NewString(),
			uuid.NewString())
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.Nil(t, readUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, err := globalUserDAO.Read(ctx, createUser.Id,
			uuid.NewString())
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.Nil(t, readUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read user by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" + random.Email(),
		Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	err = globalUserDAO.UpdatePassword(ctx, createUser.Id, createOrg.ID,
		globalHash)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	t.Run("Read user by valid email", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, readHash, err := globalUserDAO.ReadByEmail(ctx,
			createUser.Email, createOrg.Name)
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.NoError(t, err)
		require.Equal(t, globalHash, readHash)

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

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update user by valid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "dao-user-" + random.Email()
		createUser.Status = []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]

		updateUser, err := globalUserDAO.Update(ctx, createUser)
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.NoError(t, err)
		require.Equal(t, createUser.Email, updateUser.Email)
		require.Equal(t, createUser.Status, updateUser.Status)
		require.Equal(t, createUser.CreatedAt, updateUser.CreatedAt)
		require.True(t, updateUser.UpdatedAt.AsTime().After(
			updateUser.CreatedAt.AsTime()))
		require.WithinDuration(t, createUser.CreatedAt.AsTime(),
			updateUser.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Update unknown user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownUser := &api.User{Id: uuid.NewString(), OrgId: createOrg.ID,
			Email: "dao-user-" + random.Email(), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}
		updateUser, err := globalUserDAO.Update(ctx, unknownUser)
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.OrgId = uuid.NewString()
		createUser.Email = "dao-user-" + random.Email()

		updateUser, err := globalUserDAO.Update(ctx, createUser)
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update user by invalid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "dao-user-" + random.String(80)
		createUser.Status = []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]

		updateUser, err := globalUserDAO.Update(ctx, createUser)
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestUpdatePassword(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Update user password by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.UpdatePassword(ctx, createUser.Id, createOrg.ID,
			globalHash)
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Update user password by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalUserDAO.UpdatePassword(ctx, uuid.NewString(),
			createOrg.ID, globalHash)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Password updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.UpdatePassword(ctx, createUser.Id,
			uuid.NewString(), globalHash)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Delete user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.Delete(ctx, createUser.Id, createOrg.ID)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read user by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			readUser, err := globalUserDAO.Read(ctx, createUser.Id,
				createOrg.ID)
			t.Logf("readUser, err: %+v, %v", readUser, err)
			require.Nil(t, readUser)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalUserDAO.Delete(ctx, uuid.NewString(), createOrg.ID)
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.Delete(ctx, createUser.Id, uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	userIDs := []string{}
	userStatuses := []common.Status{}
	userTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)
		userIDs = append(userIDs, createUser.Id)
		userStatuses = append(userStatuses, createUser.Status)
		userTSes = append(userTSes, createUser.CreatedAt.AsTime())
	}

	t.Run("List users by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.ID,
			time.Time{}, "", 0)
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 3)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, user := range listUsers {
			if user.Id == userIDs[len(userIDs)-1] &&
				user.Status == userStatuses[len(userIDs)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.ID,
			userTSes[0], userIDs[0], 5)
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 2)
		require.Equal(t, int32(3), listCount)

		var found bool
		for _, user := range listUsers {
			if user.Id == userIDs[len(userIDs)-1] &&
				user.Status == userStatuses[len(userIDs)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, createOrg.ID,
			time.Time{}, "", 1)
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 1)
		require.Equal(t, int32(3), listCount)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx,
			uuid.NewString(), time.Time{}, "", 0)
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 0)
		require.Equal(t, int32(0), listCount)
	})

	t.Run("List users by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		listUsers, listCount, err := globalUserDAO.List(ctx, random.String(10),
			time.Time{}, "", 0)
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.Equal(t, int32(0), listCount)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}
