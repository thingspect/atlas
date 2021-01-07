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
		createUser, err := globalUserDAO.Create(ctx, user, globalHash)
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
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user,
			[]byte(random.String(61)))
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.EqualError(t, err, "invalid format: users_password_hash_check")
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
	createUser, err := globalUserDAO.Create(ctx, user, globalHash)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	t.Run("Read user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, readHash, err := globalUserDAO.Read(ctx, createUser.Id,
			createUser.OrgId)
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.NoError(t, err)
		require.Equal(t, globalHash, readHash)
		require.Equal(t, createUser, readUser)
	})

	t.Run("Read user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, readHash, err := globalUserDAO.Read(ctx, uuid.New().String(),
			uuid.New().String())
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.Nil(t, readUser)
		require.Nil(t, readHash)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, readHash, err := globalUserDAO.Read(ctx, createUser.Id,
			uuid.New().String())
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.Nil(t, readUser)
		require.Nil(t, readHash)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read user by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readUser, readHash, err := globalUserDAO.Read(ctx, random.String(10),
			createUser.OrgId)
		t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash, err)
		require.Nil(t, readUser)
		require.Nil(t, readHash)
		require.EqualError(t, err, "invalid format: UUID")
	})
}

func TestReadByEmail(t *testing.T) {
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
	createUser, err := globalUserDAO.Create(ctx, user, globalHash)
	t.Logf("createUser, err: %+v, %v", createUser, err)
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
		require.Equal(t, createUser, readUser)
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

func TestUpdateUser(t *testing.T) {
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
		createUser, err := globalUserDAO.Create(ctx, user, globalHash)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "dao-user-" + random.Email()
		createUser.Status = []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]

		updateUser, err := globalUserDAO.Update(ctx, createUser, globalHash)
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

		unknownUser := &api.User{Id: uuid.New().String(), OrgId: createOrg.ID,
			Email: "dao-user-" + random.Email(), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}
		updateUser, err := globalUserDAO.Update(ctx, unknownUser, globalHash)
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
		createUser, err := globalUserDAO.Create(ctx, user, globalHash)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.OrgId = uuid.New().String()
		createUser.Email = "dao-user-" + random.Email()

		updateUser, err := globalUserDAO.Update(ctx, createUser, globalHash)
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
		createUser, err := globalUserDAO.Create(ctx, user, globalHash)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "dao-user-" + random.Email()
		createUser.Status = []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]

		updateUser, err := globalUserDAO.Update(ctx, createUser,
			[]byte(random.String(61)))
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "invalid format: users_password_hash_check")
	})
}

func TestDeleteUser(t *testing.T) {
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
		createUser, err := globalUserDAO.Create(ctx, user, globalHash)
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

			readUser, readHash, err := globalUserDAO.Read(ctx, createUser.Id,
				createOrg.ID)
			t.Logf("readUser, readHash, err: %+v, %s, %v", readUser, readHash,
				err)
			require.Nil(t, readUser)
			require.Nil(t, readHash)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalUserDAO.Delete(ctx, uuid.New().String(), createOrg.ID)
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
		createUser, err := globalUserDAO.Create(ctx, user, globalHash)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		err = globalUserDAO.Delete(ctx, createUser.Id, uuid.New().String())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	org := org.Org{Name: "dao-user-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	var lastUserID string
	var lastUserStatus common.Status
	for i := 0; i < 3; i++ {
		user := &api.User{OrgId: createOrg.ID, Email: "dao-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}
		createUser, err := globalUserDAO.Create(ctx, user, globalHash)
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)
		lastUserID = createUser.Id
		lastUserStatus = createUser.Status
	}

	t.Run("List users by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listUsers, err := globalUserDAO.List(ctx, createOrg.ID)
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 3)

		var found bool
		for _, user := range listUsers {
			if user.Id == lastUserID && user.Status == lastUserStatus {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listUsers, err := globalUserDAO.List(ctx, uuid.New().String())
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers, 0)
	})

	t.Run("List users by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listUsers, err := globalUserDAO.List(ctx, random.String(10))
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.EqualError(t, err, "invalid format: UUID")
	})
}
