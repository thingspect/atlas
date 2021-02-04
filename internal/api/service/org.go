package service

//go:generate mockgen -source org.go -destination mock_orger_test.go -package service

import (
	"context"
	"time"

	"github.com/mennanov/fmutils"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Orger defines the methods provided by an org.DAO.
type Orger interface {
	Create(ctx context.Context, org *api.Org) (*api.Org, error)
	Read(ctx context.Context, orgID string) (*api.Org, error)
	Update(ctx context.Context, org *api.Org) (*api.Org, error)
	Delete(ctx context.Context, orgID string) error
	List(ctx context.Context, lboundTS time.Time, prevID string,
		limit int32) ([]*api.Org, int32, error)
}

// Org service contains functions to query and modify orgs.
type Org struct {
	api.UnimplementedOrgServiceServer

	orgDAO Orger
}

// NewOrg instantiates and returns a new Org service.
func NewOrg(orgDAO Orger) *Org {
	return &Org{
		orgDAO: orgDAO,
	}
}

// CreateOrg creates an organization.
func (u *Org) CreateOrg(ctx context.Context,
	req *api.CreateOrgRequest) (*api.Org, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_SYS_ADMIN {
		return nil, errPerm(common.Role_SYS_ADMIN)
	}

	org, err := u.orgDAO.Create(ctx, req.Org)
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"201")); err != nil {
		logger.Errorf("CreateOrg grpc.SetHeader: %v", err)
	}
	return org, nil
}

// GetOrg retrieves an organization by ID.
func (u *Org) GetOrg(ctx context.Context, req *api.GetOrgRequest) (*api.Org,
	error) {
	sess, ok := session.FromContext(ctx)
	if !ok || (sess.Role < common.Role_SYS_ADMIN && req.Id != sess.OrgID) {
		return nil, errPerm(common.Role_SYS_ADMIN)
	}

	org, err := u.orgDAO.Read(ctx, req.Id)
	if err != nil {
		return nil, errToStatus(err)
	}

	return org, nil
}

// UpdateOrg updates an organization. Update actions validate after merge to
// support partial updates.
func (u *Org) UpdateOrg(ctx context.Context,
	req *api.UpdateOrgRequest) (*api.Org, error) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, errPerm(common.Role_SYS_ADMIN)
	}

	if req.Org == nil {
		return nil, status.Error(codes.InvalidArgument, req.Validate().Error())
	}

	// Admins can only update their own org, system admins can update any org.
	if (sess.Role < common.Role_SYS_ADMIN && req.Org.Id != sess.OrgID) ||
		(sess.Role < common.Role_ADMIN && req.Org.Id == sess.OrgID) {
		return nil, errPerm(common.Role_SYS_ADMIN)
	}

	// Perform partial update if directed.
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		// Normalize and validate field mask.
		req.UpdateMask.Normalize()
		if !req.UpdateMask.IsValid(req.Org) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		org, err := u.orgDAO.Read(ctx, req.Org.Id)
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.Org, req.UpdateMask.Paths)
		proto.Merge(org, req.Org)
		req.Org = org
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	org, err := u.orgDAO.Update(ctx, req.Org)
	if err != nil {
		return nil, errToStatus(err)
	}

	return org, nil
}

// DeleteOrg deletes an organization by ID.
func (u *Org) DeleteOrg(ctx context.Context,
	req *api.DeleteOrgRequest) (*emptypb.Empty, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_SYS_ADMIN {
		return nil, errPerm(common.Role_SYS_ADMIN)
	}

	if err := u.orgDAO.Delete(ctx, req.Id); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("atlas-status-code",
		"204")); err != nil {
		logger.Errorf("DeleteOrg grpc.SetHeader: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ListOrgs retrieves all orgs.
func (u *Org) ListOrgs(ctx context.Context,
	req *api.ListOrgsRequest) (*api.ListOrgsResponse, error) {
	logger := alog.FromContext(ctx)
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, errPerm(common.Role_SYS_ADMIN)
	}

	// If the org does not have sufficient role, return only their org.
	if sess.Role < common.Role_SYS_ADMIN {
		org, err := u.orgDAO.Read(ctx, sess.OrgID)
		if err != nil {
			return nil, errToStatus(err)
		}

		return &api.ListOrgsResponse{
			Orgs:      []*api.Org{org},
			TotalSize: 1,
		}, nil
	}

	if req.PageSize == 0 {
		req.PageSize = defaultPageSize
	}

	lboundTS, prevID, err := session.ParsePageToken(req.PageToken)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	orgs, count, err := u.orgDAO.List(ctx, lboundTS, prevID, req.PageSize+1)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListOrgsResponse{Orgs: orgs, TotalSize: count}

	// Populate next page token.
	if len(orgs) == int(req.PageSize+1) {
		resp.Orgs = orgs[:len(orgs)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			orgs[len(orgs)-2].CreatedAt.AsTime(),
			orgs[len(orgs)-2].Id); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger.Errorf("ListOrgs session.GeneratePageToken org, err: "+
				"%+v, %v", orgs[len(orgs)-2], err)
		}
	}

	return resp, nil
}
