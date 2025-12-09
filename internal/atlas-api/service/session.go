package service

//go:generate mockgen -source session.go -destination mock_keyer_test.go -package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/thingspect/atlas/internal/atlas-api/auth"
	"github.com/thingspect/atlas/internal/atlas-api/key"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Keyer defines the methods provided by a key.DAO.
type Keyer interface {
	Create(ctx context.Context, key *api.Key) (*api.Key, error)
	Delete(ctx context.Context, keyID, orgID string) error
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32) ([]*api.Key, int32, error)
}

// Session service contains functions to create sessions and keys.
type Session struct {
	api.UnimplementedSessionServiceServer

	userDAO Userer
	keyDAO  Keyer
	cache   cache.Cacher[string]

	pwtKey []byte
}

// NewSession instantiates and returns a new Session service.
func NewSession(
	userDAO Userer, keyDAO Keyer, cache cache.Cacher[string], pwtKey []byte,
) *Session {
	return &Session{
		userDAO: userDAO,
		keyDAO:  keyDAO,
		cache:   cache,

		pwtKey: pwtKey,
	}
}

// Login logs in a user.
func (s *Session) Login(ctx context.Context, req *api.LoginRequest) (
	*api.LoginResponse, error,
) {
	logger := alog.FromContext(ctx)

	user, hash, err := s.userDAO.ReadByEmail(ctx, req.GetEmail(),
		req.GetOrgName())
	// Hash the provided password if an error is returned to prevent account
	// enumeration attacks.
	if err != nil {
		_, hashErr := auth.HashPass(req.GetPassword())
		logger.Debugf("Login s.userDAO.ReadByEmail Email, OrgName, err, "+
			"hashErr: %v, %v, %v, %v", req.GetEmail(), req.GetOrgName(), err,
			hashErr)

		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	logger.Logger = logger.WithField("userID", user.GetId()).WithField("orgID",
		user.GetOrgId())

	if err := auth.CompareHashPass(hash, req.GetPassword()); err != nil ||
		user.GetStatus() != api.Status_ACTIVE ||
		user.GetRole() < api.Role_VIEWER {
		logger.Debugf("Login crypto.CompareHashPass err, user.Status: %v, %s",
			err, user.GetStatus())

		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	token, exp, err := session.GenerateWebToken(s.pwtKey, user)
	if err != nil {
		logger.Errorf("Login session.GenerateWebToken: %v", err)

		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	return &api.LoginResponse{Token: token, ExpiresAt: exp}, nil
}

// CreateKey creates an API key.
func (s *Session) CreateKey(ctx context.Context, req *api.CreateKeyRequest) (
	*api.CreateKeyResponse, error,
) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_ADMIN {
		return nil, errPerm(api.Role_ADMIN)
	}

	// Only system admins can create keys with system admin role.
	if sess.Role < api.Role_SYS_ADMIN &&
		req.GetKey().GetRole() == api.Role_SYS_ADMIN {
		return nil, status.Error(codes.PermissionDenied,
			"permission denied, role modification not allowed")
	}

	req.Key.OrgId = sess.OrgID

	key, err := s.keyDAO.Create(ctx, req.GetKey())
	if err != nil {
		return nil, errToStatus(err)
	}

	token, err := session.GenerateKeyToken(s.pwtKey, key.GetId(),
		key.GetOrgId(), key.GetRole())
	if err != nil {
		logger.Errorf("CreateKey session.GenerateKeyToken: %v", err)

		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusCreated))); err != nil {
		logger.Errorf("CreateKey grpc.SetHeader: %v", err)
	}

	return &api.CreateKeyResponse{Key: key, Token: token}, nil
}

// DeleteKey deletes an API key by ID.
func (s *Session) DeleteKey(ctx context.Context, req *api.DeleteKeyRequest) (
	*emptypb.Empty, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_ADMIN {
		return nil, errPerm(api.Role_ADMIN)
	}

	// Disable API key before removing record. If an unclaimed key ID is given,
	// it will be confined to this org.
	if err := s.cache.Set(ctx, key.Disabled(sess.OrgID, req.GetId()),
		""); err != nil {
		return nil, errToStatus(err)
	}

	// Delete key record.
	if err := s.keyDAO.Delete(ctx, req.GetId(), sess.OrgID); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusNoContent))); err != nil {
		logger := alog.FromContext(ctx)
		logger.Errorf("DeleteKey grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListKeys retrieves all API keys.
func (s *Session) ListKeys(ctx context.Context, req *api.ListKeysRequest) (
	*api.ListKeysResponse, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_ADMIN {
		return nil, errPerm(api.Role_ADMIN)
	}

	if req.GetPageSize() == 0 {
		req.PageSize = defaultPageSize
	}

	lBoundTS, prevID, err := session.ParsePageToken(req.GetPageToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	keys, count, err := s.keyDAO.List(ctx, sess.OrgID, lBoundTS, prevID,
		req.GetPageSize()+1)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListKeysResponse{Keys: keys, TotalSize: count}

	// Populate next page token.
	if len(keys) == int(req.GetPageSize()+1) {
		resp.Keys = keys[:len(keys)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			keys[len(keys)-2].GetCreatedAt().AsTime(),
			keys[len(keys)-2].GetId()); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger := alog.FromContext(ctx)
			logger.Errorf("ListKeys session.GeneratePageToken key, err: "+
				"%+v, %v", keys[len(keys)-2], err)
		}
	}

	return resp, nil
}
