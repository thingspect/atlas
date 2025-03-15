//go:build !unit

package test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("Create valid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)
		require.NotEqual(t, user.GetId(), createUser.GetId())
		require.WithinDuration(t, time.Now(), createUser.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createUser.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerGRPCConn)
		createUser, err := userCli.CreateUser(ctx, &api.CreateUserRequest{
			User: random.User("api-user", uuid.NewString()),
		})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Create sysadmin user as non-sysadmin", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_SYS_ADMIN

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, role modification not allowed")
	})

	t.Run("Create invalid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Email = "api-user-" + random.String(80)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateUserRequest.User: embedded message failed "+
			"validation | caused by: invalid User.Email: value must be a "+
			"valid email address | caused by: mail: missing '@' or angle-addr")
	})

	t.Run("Create invalid user with invalid phone", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER
		user.Phone = random.String(10)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid E.164 phone number")
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	user := random.User("api-user", uuid.NewString())
	user.Role = api.Role_BUILDER

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	userCli := api.NewUserServiceClient(globalAdminGRPCConn)
	createUser, err := userCli.CreateUser(ctx,
		&api.CreateUserRequest{User: user})
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	t.Run("Get user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		getUser, err := userCli.GetUser(ctx,
			&api.GetUserRequest{Id: createUser.GetId()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, createUser, getUser)
	})

	t.Run("Get user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerGRPCConn)
		getUser, err := userCli.GetUser(ctx,
			&api.GetUserRequest{Id: createUser.GetId()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.Nil(t, getUser)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Get user with insufficient key role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerKeyGRPCConn)
		getUser, err := userCli.GetUser(ctx,
			&api.GetUserRequest{Id: createUser.GetId()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.Nil(t, getUser)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Get user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		getUser, err := userCli.GetUser(ctx,
			&api.GetUserRequest{Id: uuid.NewString()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.Nil(t, getUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Gets are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewUserServiceClient(secondaryAdminGRPCConn)
		getUser, err := secCli.GetUser(ctx,
			&api.GetUserRequest{Id: createUser.GetId()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.Nil(t, getUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("Update user by valid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Name = "api-user-" + random.String(10)
		createUser.Email = "api-user-" + random.Email()
		createUser.Role = api.Role_ADMIN
		createUser.Status = api.Status_DISABLED
		createUser.Tags = random.Tags("api-user", 2)

		updateUser, err := userCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: createUser})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.NoError(t, err)
		require.Equal(t, createUser.GetName(), updateUser.GetName())
		require.Equal(t, createUser.GetEmail(), updateUser.GetEmail())
		require.Equal(t, createUser.GetRole(), updateUser.GetRole())
		require.Equal(t, createUser.GetStatus(), updateUser.GetStatus())
		require.Equal(t, createUser.GetTags(), updateUser.GetTags())
		require.True(t, updateUser.GetUpdatedAt().AsTime().After(
			updateUser.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createUser.GetCreatedAt().AsTime(),
			updateUser.GetUpdatedAt().AsTime(), 2*time.Second)

		getUser, err := userCli.GetUser(ctx,
			&api.GetUserRequest{Id: createUser.GetId()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateUser, getUser)
	})

	t.Run("Partial update user by valid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminKeyGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		part := &api.User{
			Id: createUser.GetId(), Name: "api-user-" + random.String(10),
			Email: "api-user-" + random.Email(), Role: api.Role_ADMIN,
			Status: api.Status_DISABLED, Tags: random.Tags("api-user", 2),
		}

		updateUser, err := userCli.UpdateUser(ctx, &api.UpdateUserRequest{
			User: part, UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{
				"name", "email", "role", "status", "tags",
			}},
		})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.NoError(t, err)
		require.Equal(t, part.GetName(), updateUser.GetName())
		require.Equal(t, part.GetEmail(), updateUser.GetEmail())
		require.Equal(t, part.GetRole(), updateUser.GetRole())
		require.Equal(t, part.GetStatus(), updateUser.GetStatus())
		require.Equal(t, part.GetTags(), updateUser.GetTags())
		require.True(t, updateUser.GetUpdatedAt().AsTime().After(
			updateUser.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createUser.GetCreatedAt().AsTime(),
			updateUser.GetUpdatedAt().AsTime(), 2*time.Second)

		getUser, err := userCli.GetUser(ctx,
			&api.GetUserRequest{Id: createUser.GetId()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateUser, getUser)
	})

	t.Run("Update nil user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		updateUser, err := userCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: nil})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateUserRequest.User: value is required")
	})

	t.Run("Update user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerGRPCConn)
		updateUser, err := userCli.UpdateUser(ctx, &api.UpdateUserRequest{
			User: random.User("api-user", uuid.NewString()),
		})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Update user with insufficient key role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerKeyGRPCConn)
		updateUser, err := userCli.UpdateUser(ctx, &api.UpdateUserRequest{
			User: random.User("api-user", uuid.NewString()),
		})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Update user role to sysadmin as non-sysadmin", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Role = api.Role_SYS_ADMIN

		updateUser, err := userCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: createUser})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, role modification not allowed")
	})

	t.Run("Update user with invalid phone", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Phone = random.String(10)

		updateUser, err := userCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: createUser})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid E.164 phone number")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_ADMIN

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		updateUser, err := userCli.UpdateUser(ctx, &api.UpdateUserRequest{
			User: user, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"},
			},
		})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update user by unknown user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_ADMIN

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		updateUser, err := userCli.UpdateUser(ctx, &api.UpdateUserRequest{
			User: user, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"email"},
			},
		})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update user by unknown user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		updateUser, err := userCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: user})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.OrgId = uuid.NewString()
		createUser.Email = "api-user-" + random.Email()

		secCli := api.NewUserServiceClient(secondaryAdminGRPCConn)
		updateUser, err := secCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: createUser})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update user validation failure", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = "api-user-" + random.String(10)

		updateUser, err := userCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: createUser})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateUserRequest.User: embedded message failed "+
			"validation | caused by: invalid User.Email: value must be a "+
			"valid email address | caused by: mail: missing '@' or angle-addr")
	})

	t.Run("Update user by invalid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.Email = fmt.Sprintf("%s@%s.com", random.String(64),
			random.String(15))

		updateUser, err := userCli.UpdateUser(ctx,
			&api.UpdateUserRequest{User: createUser})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid format: value too long")
	})
}

func TestUpdateUserPassword(t *testing.T) {
	t.Parallel()

	t.Run("Update user password by valid ID", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		_, err = userCli.UpdateUserPassword(ctx, &api.UpdateUserPasswordRequest{
			Id: createUser.GetId(), Password: random.String(20),
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Update user password with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerGRPCConn)
		_, err := userCli.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{
				Id: uuid.NewString(), Password: random.String(20),
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Update user password with insufficient key role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerKeyGRPCConn)
		_, err := userCli.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{
				Id: uuid.NewString(), Password: random.String(20),
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Update user password with weak password", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		_, err = userCli.UpdateUserPassword(ctx, &api.UpdateUserPasswordRequest{
			Id: createUser.GetId(), Password: "1234567890",
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"weak password, see NIST password guidelines")
	})

	t.Run("Update user password by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		_, err := userCli.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{
				Id: uuid.NewString(), Password: random.String(20),
			})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Password updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		secCli := api.NewUserServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.UpdateUserPassword(ctx, &api.UpdateUserPasswordRequest{
			Id: createUser.GetId(), Password: random.String(20),
		})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("Delete user by valid ID", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		_, err = userCli.DeleteUser(ctx,
			&api.DeleteUserRequest{Id: createUser.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read user by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(),
				testTimeout)
			defer cancel()

			userCli := api.NewUserServiceClient(globalAdminKeyGRPCConn)
			getUser, err := userCli.GetUser(ctx,
				&api.GetUserRequest{Id: createUser.GetId()})
			t.Logf("getUser, err: %+v, %v", getUser, err)
			require.Nil(t, getUser)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(secondaryViewerGRPCConn)
		_, err := userCli.DeleteUser(ctx,
			&api.DeleteUserRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied "+
			"desc = permission denied, ADMIN role required")
	})

	t.Run("Delete user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		_, err := userCli.DeleteUser(ctx,
			&api.DeleteUserRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		secCli := api.NewUserServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteUser(ctx,
			&api.DeleteUserRequest{Id: createUser.GetId()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	userIDs := []string{}
	userNames := []string{}
	userRoles := []api.Role{}
	userStatuses := []api.Status{}
	userTags := [][]string{}
	for range 3 {
		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx,
			&api.CreateUserRequest{User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		userIDs = append(userIDs, createUser.GetId())
		userNames = append(userNames, createUser.GetName())
		userRoles = append(userRoles, createUser.GetRole())
		userStatuses = append(userStatuses, createUser.GetStatus())
		userTags = append(userTags, createUser.GetTags())
	}

	t.Run("List users by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		listUsers, err := userCli.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listUsers.GetUsers()), 3)
		require.GreaterOrEqual(t, listUsers.GetTotalSize(), int32(3))

		var found bool
		for _, user := range listUsers.GetUsers() {
			if user.GetId() == userIDs[len(userIDs)-1] &&
				user.GetName() == userNames[len(userNames)-1] &&
				user.GetRole() == userRoles[len(userRoles)-1] &&
				user.GetStatus() == userStatuses[len(userStatuses)-1] &&
				reflect.DeepEqual(user.GetTags(), userTags[len(userTags)-1]) {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminKeyGRPCConn)
		listUsers, err := userCli.ListUsers(ctx,
			&api.ListUsersRequest{PageSize: 2})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers.GetUsers(), 2)
		require.NotEmpty(t, listUsers.GetNextPageToken())
		require.GreaterOrEqual(t, listUsers.GetTotalSize(), int32(3))

		nextUsers, err := userCli.ListUsers(ctx, &api.ListUsersRequest{
			PageSize: 2, PageToken: listUsers.GetNextPageToken(),
		})
		t.Logf("nextUsers, err: %+v, %v", nextUsers, err)
		require.NoError(t, err)
		require.NotEmpty(t, nextUsers.GetUsers())
		require.GreaterOrEqual(t, nextUsers.GetTotalSize(), int32(3))
	})

	t.Run("List users with tag filter", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		listUsers, err := userCli.ListUsers(ctx,
			&api.ListUsersRequest{Tag: userTags[2][0]})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers.GetUsers(), 1)
		require.Equal(t, int32(1), listUsers.GetTotalSize())

		require.Equal(t, userIDs[len(userIDs)-1], listUsers.GetUsers()[0].GetId())
		require.Equal(t, userNames[len(userNames)-1], listUsers.GetUsers()[0].GetName())
		require.Equal(t, userStatuses[len(userStatuses)-1],
			listUsers.GetUsers()[0].GetStatus())
		require.Equal(t, userRoles[len(userRoles)-1], listUsers.GetUsers()[0].GetRole())
		require.Equal(t, userTags[len(userTags)-1], listUsers.GetUsers()[0].GetTags())
	})

	t.Run("List no users by key role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewUserServiceClient(secondaryViewerKeyGRPCConn)
		listUsers, err := secCli.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Empty(t, listUsers.GetUsers())
		require.Equal(t, int32(0), listUsers.GetTotalSize())
	})

	t.Run("List own user with non-admin role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewUserServiceClient(secondaryViewerGRPCConn)
		listUsers, err := secCli.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers.GetUsers(), 1)
		require.Equal(t, int32(1), listUsers.GetTotalSize())
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		secCli := api.NewUserServiceClient(secondaryAdminGRPCConn)
		listUsers, err := secCli.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers.GetUsers(), 1)
		require.Equal(t, int32(1), listUsers.GetTotalSize())
	})

	t.Run("List users by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		listUsers, err := userCli.ListUsers(ctx,
			&api.ListUsersRequest{PageToken: badUUID})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}
