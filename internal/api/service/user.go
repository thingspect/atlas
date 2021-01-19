package service

//go:generate mockgen -source user.go -destination mock_userer_test.go -package service

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mennanov/fmutils"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Userer defines the methods provided by a user.DAO.
type Userer interface {
	Create(ctx context.Context, user *api.User) (*api.User, error)
	Read(ctx context.Context, userID, orgID string) (*api.User, error)
	ReadByEmail(ctx context.Context, email, orgName string) (*api.User, []byte,
		error)
	Update(ctx context.Context, user *api.User) (*api.User, error)
	UpdatePassword(ctx context.Context, userID, orgID string,
		passHash []byte) error
	Delete(ctx context.Context, userID, orgID string) error
	List(ctx context.Context, orgID string, lboundTS time.Time, prevID string,
		limit int32) ([]*api.User, int32, error)
}

// User service contains functions to query and modify users.
type User struct {
	api.UnimplementedUserServiceServer

	userDAO Userer
}

// NewUser instantiates and returns a new User service.
func NewUser(userDAO Userer) *User {
	return &User{
		userDAO: userDAO,
	}
}

// CreateUser creates a user.
func (u *User) CreateUser(ctx context.Context,
	req *api.CreateUserRequest) (*api.User, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	req.User.OrgId = sess.OrgID

	user, err := u.userDAO.Create(ctx, req.User)
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"201")); err != nil {
		logger.Errorf("CreateUser grpc.SetHeader: %v", err)
	}
	return user, nil
}

// GetUser retrieves a user by ID.
func (u *User) GetUser(ctx context.Context,
	req *api.GetUserRequest) (*api.User, error) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	user, err := u.userDAO.Read(ctx, req.Id, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return user, nil
}

// UpdateUser updates a user.
func (u *User) UpdateUser(ctx context.Context,
	req *api.UpdateUserRequest) (*api.User, error) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.User == nil {
		return nil, status.Error(codes.InvalidArgument, req.Validate().Error())
	}
	req.User.OrgId = sess.OrgID

	// Perform partial update if directed.
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		// Normalize and validate field mask.
		req.UpdateMask.Normalize()
		if !req.UpdateMask.IsValid(req.User) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		user, err := u.userDAO.Read(ctx, req.User.Id, sess.OrgID)
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.User, req.UpdateMask.Paths)
		proto.Merge(user, req.User)
		req.User = user
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := u.userDAO.Update(ctx, req.User)
	if err != nil {
		return nil, errToStatus(err)
	}

	return user, nil
}

// UpdateUserPassword updates a user's password by ID.
func (u *User) UpdateUserPassword(ctx context.Context,
	req *api.UpdateUserPasswordRequest) (*empty.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if err := crypto.CheckPass(req.Password); err != nil {
		return nil, errToStatus(err)
	}

	hash, err := crypto.HashPass(req.Password)
	if err != nil {
		logger.Errorf("UpdateUserPassword crypto.HashPass: %v", err)
		return nil, errToStatus(crypto.ErrWeakPass)
	}

	if err := u.userDAO.UpdatePassword(ctx, req.Id, sess.OrgID,
		hash); err != nil {
		return nil, errToStatus(err)
	}

	return &empty.Empty{}, nil
}

// DeleteUser deletes a user by ID.
func (u *User) DeleteUser(ctx context.Context,
	req *api.DeleteUserRequest) (*empty.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if err := u.userDAO.Delete(ctx, req.Id, sess.OrgID); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"204")); err != nil {
		logger.Errorf("DeleteUser grpc.SetHeader: %v", err)
	}
	return &empty.Empty{}, nil
}

// ListUsers retrieves all users.
func (u *User) ListUsers(ctx context.Context,
	req *api.ListUsersRequest) (*api.ListUsersResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}

	lboundTS, prevID, err := session.ParsePageToken(req.PageToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	users, count, err := u.userDAO.List(ctx, sess.OrgID, lboundTS, prevID,
		req.PageSize+1)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListUsersResponse{
		Users:         users,
		PrevPageToken: req.PageToken,
		TotalSize:     count,
	}

	// Populate next page token.
	if len(users) == int(req.PageSize+1) {
		resp.Users = users[:len(users)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			users[len(users)-2].CreatedAt.AsTime(),
			users[len(users)-2].Id); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger.Errorf("ListUsers session.GeneratePageToken user, err: "+
				"%+v, %v", users[len(users)-2], err)
		}
	}

	return resp, nil
}
