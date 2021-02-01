// +build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/crypto"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("Create valid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Create(gomock.Any(), user).Return(user, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: user.OrgId,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		createUser, err := userSvc.CreateUser(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(user, createUser) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", user, createUser)
		}
	})

	t.Run("Create user with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		createUser, err := userSvc.CreateUser(ctx, &api.CreateUserRequest{
			User: nil})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Create user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_BUILDER}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		createUser, err := userSvc.CreateUser(ctx, &api.CreateUserRequest{
			User: nil})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Create invalid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Create(gomock.Any(), user).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: user.OrgId,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		createUser, err := userSvc.CreateUser(ctx, &api.CreateUserRequest{
			User: user})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.Nil(t, createUser)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	t.Run("Get user by valid ID", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Read(gomock.Any(), user.Id, user.OrgId).Return(user,
			nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: user.OrgId,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		getUser, err := userSvc.GetUser(ctx, &api.GetUserRequest{Id: user.Id})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(user, getUser) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", user, getUser)
		}
	})

	t.Run("Get user with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		getUser, err := userSvc.GetUser(ctx, &api.GetUserRequest{
			Id: uuid.NewString()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.Nil(t, getUser)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Get user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{UserID: uuid.NewString(),
				OrgID: uuid.NewString(), Role: common.Role_VIEWER}),
			2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		getUser, err := userSvc.GetUser(ctx, &api.GetUserRequest{
			Id: uuid.NewString()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.Nil(t, getUser)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Get user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		getUser, err := userSvc.GetUser(ctx, &api.GetUserRequest{
			Id: uuid.NewString()})
		t.Logf("getUser, err: %+v, %v", getUser, err)
		require.Nil(t, getUser)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("Update user by valid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Update(gomock.Any(), user).Return(user, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: user.OrgId,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: user})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(user, updateUser) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", user, updateUser)
		}
	})

	t.Run("Partial update user by valid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		part := &api.User{Id: user.Id, Status: api.Status_ACTIVE}
		merged := &api.User{Id: user.Id, OrgId: user.OrgId, Email: user.Email,
			Role: user.Role, Status: api.Status_ACTIVE}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Read(gomock.Any(), user.Id, user.OrgId).Return(user,
			nil).Times(1)
		userer.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(merged)).
			Return(merged, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: user.OrgId,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"}}})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(merged, updateUser) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", merged, updateUser)
		}
	})

	t.Run("Update user with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: nil})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Update nil user", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: nil})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateUserRequest.User: value is required"), err)
	})

	t.Run("Update user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{UserID: uuid.NewString(),
				OrgID: uuid.NewString(), Role: common.Role_VIEWER}),
			2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: random.User("api-user", uuid.NewString())})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: user, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"}}})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid field mask"), err)
	})

	t.Run("Partial update user by unknown user", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		part := &api.User{Id: uuid.NewString(), Status: api.Status_ACTIVE}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Read(gomock.Any(), part.Id, orgID).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: part, UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"status"}}})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Update user validation failure", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Email = random.String(10)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: user.OrgId,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: user})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid "+
			"UpdateUserRequest.User: embedded message failed validation | "+
			"caused by: invalid User.Email: value must be a valid email "+
			"address | caused by: mail: missing '@' or angle-addr"), err)
	})

	t.Run("Update user by invalid user", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())
		user.Email = random.String(54) + random.Email()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Update(gomock.Any(), user).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: user.OrgId,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		updateUser, err := userSvc.UpdateUser(ctx, &api.UpdateUserRequest{
			User: user})
		t.Logf("updateUser, err: %+v, %v", updateUser, err)
		require.Nil(t, updateUser)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestUpdateUserPassword(t *testing.T) {
	t.Parallel()

	t.Run("Update user password by valid ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{Id: uuid.NewString(),
				Password: random.String(20)})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Update user password with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{Id: uuid.NewString(),
				Password: random.String(20)})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Update user password with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{UserID: uuid.NewString(),
				OrgID: uuid.NewString(), Role: common.Role_VIEWER}),
			2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{Id: uuid.NewString(),
				Password: random.String(20)})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Update user password with weak password", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{Id: uuid.NewString(),
				Password: "1234567890"})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.InvalidArgument,
			crypto.ErrWeakPass.Error()), err)
	})

	t.Run("Update user password by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any()).Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 6*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.UpdateUserPassword(ctx,
			&api.UpdateUserPasswordRequest{Id: uuid.NewString(),
				Password: random.String(20)})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("Delete user by valid ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.DeleteUser(ctx, &api.DeleteUserRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete user with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.DeleteUser(ctx, &api.DeleteUserRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Delete user with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_BUILDER}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.DeleteUser(ctx, &api.DeleteUserRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("Delete user by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		_, err := userSvc.DeleteUser(ctx, &api.DeleteUserRequest{
			Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestListUsers(t *testing.T) {
	t.Parallel()

	t.Run("List users by valid org ID", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		users := []*api.User{
			random.User("api-user", uuid.NewString()),
			random.User("api-user", uuid.NewString()),
			random.User("api-user", uuid.NewString()),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(51)).
			Return(users, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listUsers.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListUsersResponse{Users: users, TotalSize: 3},
			listUsers) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListUsersResponse{Users: users, TotalSize: 3}, listUsers)
		}
	})

	t.Run("List users by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		users := []*api.User{
			random.User("api-user", uuid.NewString()),
			random.User("api-user", uuid.NewString()),
			random.User("api-user", uuid.NewString()),
		}

		next, err := session.GeneratePageToken(users[1].CreatedAt.AsTime(),
			users[1].Id)
		require.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(users, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{
			PageSize: 2})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listUsers.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListUsersResponse{Users: users[:2],
			NextPageToken: next, TotalSize: 3}, listUsers) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListUsersResponse{Users: users[:2], NextPageToken: next,
					TotalSize: 3}, listUsers)
		}
	})

	t.Run("List users with invalid session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.Equal(t, errPerm(common.Role_ADMIN), err)
	})

	t.Run("List users with insufficient role", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Read(gomock.Any(), user.Id, user.OrgId).Return(user,
			nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{UserID: user.Id,
				OrgID: user.OrgId, Role: common.Role_VIEWER}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListUsersResponse{Users: []*api.User{user},
			TotalSize: 1}, listUsers) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListUsersResponse{
				Users: []*api.User{user}, TotalSize: 1}, listUsers)
		}
	})

	t.Run("List users by unknown ID", func(t *testing.T) {
		t.Parallel()

		user := random.User("api-user", uuid.NewString())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{UserID: user.Id,
				OrgID: user.OrgId, Role: common.Role_VIEWER}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("List users by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userer := NewMockUserer(ctrl)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: uuid.NewString(),
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{
			PageToken: "..."})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid page token"), err)
	})

	t.Run("List users by invalid org ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil, int32(0), dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: "aaa",
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.Nil(t, listUsers)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("List users with generation failure", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		users := []*api.User{
			random.User("api-user", uuid.NewString()),
			random.User("api-user", uuid.NewString()),
			random.User("api-user", uuid.NewString()),
		}
		users[1].Id = "..."

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(users, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{OrgID: orgID,
				Role: common.Role_ADMIN}), 2*time.Second)
		defer cancel()

		userSvc := NewUser(userer)
		listUsers, err := userSvc.ListUsers(ctx, &api.ListUsersRequest{
			PageSize: 2})
		t.Logf("listUsers, err: %+v, %v", listUsers, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listUsers.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListUsersResponse{Users: users[:2],
			TotalSize: 3}, listUsers) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListUsersResponse{Users: users[:2], TotalSize: 3},
				listUsers)
		}
	})
}
