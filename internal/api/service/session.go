package service

//go:generate mockgen -source session.go -destination mock_userer_test.go -package service

import (
	"context"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/crypto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Userer defines the methods provided by a user.DAO.
type Userer interface {
	ReadByEmail(ctx context.Context, email, orgName string) (*api.User, []byte,
		error)
}

// Session service contains functions to create sessions and keys.
type Session struct {
	api.UnimplementedSessionServiceServer

	userDAO Userer
	pwtKey  []byte
}

// NewSession instantiates and returns a new Session service.
func NewSession(userDAO Userer, pwtKey []byte) *Session {
	return &Session{
		userDAO: userDAO,
		pwtKey:  pwtKey,
	}
}

// Login logs in a user.
func (s *Session) Login(ctx context.Context,
	req *api.LoginRequest) (*api.LoginResponse, error) {
	user, hash, err := s.userDAO.ReadByEmail(ctx, req.Email, req.OrgName)
	// Hash the provided password if an error is returned to prevent account
	// enumeration attacks.
	if err != nil {
		_, hashErr := crypto.HashPass(req.Password)
		alog.Debugf("Login s.userDAO.ReadByEmail Email, OrgName, err, "+
			"hashErr: %v, %v, %v, %v", req.Email, req.OrgName, err, hashErr)
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if err := crypto.CompareHashPass(hash, req.Password); err != nil ||
		user.IsDisabled {
		alog.Debugf("Login crypto.CompareHashPass Email, OrgName, err, "+
			"user.IsDisabled: %v, %v, %v, %v", req.Email, req.OrgName, err,
			user.IsDisabled)
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	token, exp, err := session.GenerateToken(s.pwtKey, user.Id, user.OrgId)
	if err != nil {
		alog.Errorf("Login crypto.GenerateToken Email, OrgName, err: %v, %v, "+
			"%v", req.Email, req.OrgName, err)
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	return &api.LoginResponse{Token: token, ExpiresAt: exp}, nil
}
