// +build !unit

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("Create valid user", func(t *testing.T) {
		t.Parallel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)
		require.NotNil(t, createUser.User)
		require.Equal(t, globalAuthOrgID, createUser.User.OrgId)
		require.Equal(t, user.Email, createUser.User.Email)
		require.Equal(t, user.Status, createUser.User.Status)
		require.WithinDuration(t, time.Now(),
			createUser.User.CreatedAt.AsTime(), 2*time.Second)
		require.WithinDuration(t, time.Now(),
			createUser.User.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Create invalid user", func(t *testing.T) {
		t.Parallel()

		user := &api.User{Email: "api-user-" + random.String(10),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateUserRequest.User: embedded message failed "+
			"validation | caused by: invalid User.Email: value must be a "+
			"valid email address | caused by: mail: missing '@' or angle-addr")
	})
}

func TestReadUser(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := &api.User{Email: "api-user-" + random.Email(),
		Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}

	userCli := api.NewUserServiceClient(globalAuthGRPCConn)
	createUser, err := userCli.Create(ctx, &api.CreateUserRequest{User: user})
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	t.Run("Read user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		readUser, err := userCli.Read(ctx, &api.ReadUserRequest{
			Id: createUser.User.Id})
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ReadUserResponse{User: createUser.User},
			readUser) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ReadUserResponse{User: createUser.User}, readUser)
		}
	})

	t.Run("Read user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		readUser, err := userCli.Read(ctx, &api.ReadUserRequest{
			Id: uuid.New().String()})
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.Nil(t, readUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Reads are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		secCli := api.NewUserServiceClient(secondaryAuthGRPCConn)
		readUser, err := secCli.Read(ctx, &api.ReadUserRequest{
			Id: createUser.User.Id})
		t.Logf("readUser, err: %+v, %v", readUser, err)
		require.Nil(t, readUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("Update user by valid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.User.Email = "api-user-" + random.Email()
		createUser.User.Status = []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]

		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: createUser.User})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.NoError(t, err)
		require.Equal(t, createUser.User.Email, updateUser.User.Email)
		require.Equal(t, createUser.User.Status, updateUser.User.Status)
		require.Equal(t, createUser.User.CreatedAt.AsTime(),
			updateUser.User.CreatedAt.AsTime())
		require.True(t, updateUser.User.UpdatedAt.AsTime().After(
			updateUser.User.CreatedAt.AsTime()))
		require.WithinDuration(t, createUser.User.CreatedAt.AsTime(),
			updateUser.User.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Partial update user by valid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		part := &api.User{Id: createUser.User.Id, Email: "api-user-" +
			random.Email(), Status: []common.Status{common.Status_ACTIVE,
			common.Status_DISABLED}[random.Intn(2)]}

		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"email", "status"}}})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.NoError(t, err)
		require.Equal(t, part.Email, updateUser.User.Email)
		require.Equal(t, part.Status, updateUser.User.Status)
		require.Equal(t, createUser.User.CreatedAt.AsTime(),
			updateUser.User.CreatedAt.AsTime())
		require.True(t, updateUser.User.UpdatedAt.AsTime().After(
			updateUser.User.CreatedAt.AsTime()))
		require.WithinDuration(t, createUser.User.CreatedAt.AsTime(),
			updateUser.User.UpdatedAt.AsTime(), 2*time.Second)
	})

	t.Run("Update nil user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: nil})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateUserRequest.User: value is required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownUser := &api.User{Id: uuid.New().String(), Email: "api-user-" +
			random.Email()}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: unknownUser, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"}}})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update user by unknown user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownUser := &api.User{Id: uuid.New().String(), Email: "api-user-" +
			random.Email()}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: unknownUser, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"email"}}})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update user by unknown user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownUser := &api.User{Id: uuid.New().String(),
			Email: "api-user-" + random.Email(), Status: []common.Status{
				common.Status_ACTIVE, common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: unknownUser})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.User.OrgId = uuid.New().String()
		createUser.User.Email = "api-user-" + random.Email()

		secCli := api.NewUserServiceClient(secondaryAuthGRPCConn)
		updateUser, err := secCli.Update(ctx, &api.UpdateUserRequest{
			User: createUser.User})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update user validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.User.Email = "api-user-" + random.String(10)

		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: createUser.User})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateUserRequest.User: embedded message failed "+
			"validation | caused by: invalid User.Email: value must be a "+
			"valid email address | caused by: mail: missing '@' or angle-addr")
	})

	t.Run("Update user by invalid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Update user fields.
		createUser.User.Email = fmt.Sprintf("%s@%s.com", random.String(64),
			random.String(15))

		updateUser, err := userCli.Update(ctx, &api.UpdateUserRequest{
			User: createUser.User})
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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		_, err = userCli.UpdatePassword(ctx, &api.UpdateUserPasswordRequest{
			Id: createUser.User.Id, Password: random.String(20)})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Update user password with weak password", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		_, err = userCli.UpdatePassword(ctx, &api.UpdateUserPasswordRequest{
			Id: createUser.User.Id, Password: "1234567890"})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"weak password, see NIST password guidelines")
	})

	t.Run("Update user password by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		_, err := userCli.UpdatePassword(ctx, &api.UpdateUserPasswordRequest{
			Id: uuid.New().String(), Password: random.String(20)})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Password updates are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		secCli := api.NewUserServiceClient(secondaryAuthGRPCConn)
		_, err = secCli.UpdatePassword(ctx, &api.UpdateUserPasswordRequest{
			Id: createUser.User.Id, Password: random.String(20)})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("Delete user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		_, err = userCli.Delete(ctx, &api.DeleteUserRequest{
			Id: createUser.User.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read user by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			userCli := api.NewUserServiceClient(globalAuthGRPCConn)
			readUser, err := userCli.Read(ctx, &api.ReadUserRequest{
				Id: createUser.User.Id})
			t.Logf("readUser, err: %+v, %v", readUser, err)
			require.Nil(t, readUser)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		_, err := userCli.Delete(ctx, &api.DeleteUserRequest{
			Id: uuid.New().String()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		secCli := api.NewUserServiceClient(secondaryAuthGRPCConn)
		_, err = secCli.Delete(ctx, &api.DeleteUserRequest{
			Id: createUser.User.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	userIDs := []string{}
	userStatuses := []common.Status{}
	for i := 0; i < 3; i++ {
		user := &api.User{Email: "api-user-" + random.Email(),
			Status: []common.Status{common.Status_ACTIVE,
				common.Status_DISABLED}[random.Intn(2)]}

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		createUser, err := userCli.Create(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)
		userIDs = append(userIDs, createUser.User.Id)
		userStatuses = append(userStatuses, createUser.User.Status)
	}

	t.Run("List users by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		listUsers, err := userCli.List(ctx, &api.ListUserRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listUsers.Users), 3)
		require.GreaterOrEqual(t, listUsers.TotalSize, int32(3))

		var found bool
		for _, user := range listUsers.Users {
			if user.Id == userIDs[len(userIDs)-1] &&
				user.Status == userStatuses[len(userIDs)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List users by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		listUsers, err := userCli.List(ctx, &api.ListUserRequest{PageSize: 2})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers.Users, 2)
		require.Empty(t, listUsers.PrevPageToken)
		require.NotEmpty(t, listUsers.NextPageToken)
		require.GreaterOrEqual(t, listUsers.TotalSize, int32(3))

		nextUsers, err := userCli.List(ctx, &api.ListUserRequest{PageSize: 2,
			PageToken: listUsers.NextPageToken})
		t.Logf("nextUsers, err: %+v, %v", nextUsers, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(nextUsers.Users), 1)
		require.NotEmpty(t, nextUsers.PrevPageToken)
		require.GreaterOrEqual(t, nextUsers.TotalSize, int32(3))
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		secCli := api.NewUserServiceClient(secondaryAuthGRPCConn)
		listUsers, err := secCli.List(ctx, &api.ListUserRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Len(t, listUsers.Users, 1)
		require.Equal(t, int32(1), listUsers.TotalSize)
	})

	t.Run("List users by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		userCli := api.NewUserServiceClient(globalAuthGRPCConn)
		listUsers, err := userCli.List(ctx, &api.ListUserRequest{
			PageToken: "..."})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}
