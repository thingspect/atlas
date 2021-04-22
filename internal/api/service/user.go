package service

//go:generate mockgen -source user.go -destination mock_userer_test.go -package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/mennanov/fmutils"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/crypto"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// E.164 format: https://www.twilio.com/docs/glossary/what-e164
var rePhone = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

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
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32, tag string) ([]*api.User, int32, error)
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
	if !ok || sess.Role < common.Role_ADMIN {
		return nil, errPerm(common.Role_ADMIN)
	}

	// Only system admins can elevate to system admin.
	if sess.Role < common.Role_SYS_ADMIN &&
		req.User.Role == common.Role_SYS_ADMIN {
		return nil, status.Error(codes.PermissionDenied,
			"permission denied, role modification not allowed")
	}

	// Validate phone number.
	if req.User.Phone != "" && !rePhone.MatchString(req.User.Phone) {
		return nil, status.Error(codes.InvalidArgument,
			"invalid E.164 phone number")
	}

	req.User.OrgId = sess.OrgID
	req.User.Tags = append(req.User.Tags,
		strings.ToLower(req.User.Role.String()))

	user, err := u.userDAO.Create(ctx, req.User)
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		"201")); err != nil {
		logger.Errorf("CreateUser grpc.SetHeader: %v", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID.
func (u *User) GetUser(ctx context.Context,
	req *api.GetUserRequest) (*api.User, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || (sess.Role < common.Role_ADMIN && req.Id != sess.UserID) {
		return nil, errPerm(common.Role_ADMIN)
	}

	user, err := u.userDAO.Read(ctx, req.Id, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return user, nil
}

// UpdateUser updates a user. Update actions validate after merge to support
// partial updates.
func (u *User) UpdateUser(ctx context.Context,
	req *api.UpdateUserRequest) (*api.User, error) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, errPerm(common.Role_ADMIN)
	}

	if req.User == nil {
		return nil, status.Error(codes.InvalidArgument,
			req.Validate(false).Error())
	}
	req.User.OrgId = sess.OrgID

	// Non-admins can only update their own user.
	if sess.Role < common.Role_ADMIN && req.User.Id != sess.UserID {
		return nil, errPerm(common.Role_ADMIN)
	}

	// Only admins can update roles, and only system admins can elevate to
	// system admin.
	if (sess.Role < common.Role_ADMIN && req.User.Role != sess.Role) ||
		(sess.Role < common.Role_SYS_ADMIN &&
			req.User.Role == common.Role_SYS_ADMIN) {
		return nil, status.Error(codes.PermissionDenied,
			"permission denied, role modification not allowed")
	}

	// Validate phone number.
	if req.User.Phone != "" && !rePhone.MatchString(req.User.Phone) {
		return nil, status.Error(codes.InvalidArgument,
			"invalid E.164 phone number")
	}

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
		if req.User.Tags != nil {
			user.Tags = nil
		}
		proto.Merge(user, req.User)
		req.User = user
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(false); err != nil {
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
	req *api.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || (sess.Role < common.Role_ADMIN && req.Id != sess.UserID) {
		return nil, errPerm(common.Role_ADMIN)
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

	return &emptypb.Empty{}, nil
}

// DeleteUser deletes a user by ID.
func (u *User) DeleteUser(ctx context.Context,
	req *api.DeleteUserRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_ADMIN {
		return nil, errPerm(common.Role_ADMIN)
	}

	if err := u.userDAO.Delete(ctx, req.Id, sess.OrgID); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		"204")); err != nil {
		logger.Errorf("DeleteUser grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListUsers retrieves all users.
func (u *User) ListUsers(ctx context.Context,
	req *api.ListUsersRequest) (*api.ListUsersResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, errPerm(common.Role_ADMIN)
	}

	// If the user does not have sufficient role, return only their user. Will
	// not be found for API key tokens.
	if sess.Role < common.Role_ADMIN {
		user, err := u.userDAO.Read(ctx, sess.UserID, sess.OrgID)
		if err != nil {
			return nil, errToStatus(err)
		}

		return &api.ListUsersResponse{
			Users:     []*api.User{user},
			TotalSize: 1,
		}, nil
	}

	if req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}

	lBoundTS, prevID, err := session.ParsePageToken(req.PageToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	users, count, err := u.userDAO.List(ctx, sess.OrgID, lBoundTS, prevID,
		req.PageSize+1, req.Tag)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListUsersResponse{Users: users, TotalSize: count}

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
