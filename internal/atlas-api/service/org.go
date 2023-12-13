package service

//go:generate mockgen -source org.go -destination mock_orger_test.go -package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/mennanov/fmutils"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/proto/go/api"
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
	List(ctx context.Context, lBoundTS time.Time, prevID string,
		limit int32) ([]*api.Org, int32, error)
}

// Org service contains functions to query and modify organizations.
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
func (o *Org) CreateOrg(ctx context.Context, req *api.CreateOrgRequest) (
	*api.Org, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_SYS_ADMIN {
		return nil, errPerm(api.Role_SYS_ADMIN)
	}

	org, err := o.orgDAO.Create(ctx, req.GetOrg())
	if err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusCreated))); err != nil {
		logger := alog.FromContext(ctx)
		logger.Errorf("CreateOrg grpc.SetHeader: %v", err)
	}

	return org, nil
}

// GetOrg retrieves an organization by ID.
func (o *Org) GetOrg(ctx context.Context, req *api.GetOrgRequest) (
	*api.Org, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || (sess.Role < api.Role_SYS_ADMIN && req.GetId() != sess.OrgID) {
		return nil, errPerm(api.Role_SYS_ADMIN)
	}

	org, err := o.orgDAO.Read(ctx, req.GetId())
	if err != nil {
		return nil, errToStatus(err)
	}

	return org, nil
}

// UpdateOrg updates an organization. Update actions validate after merge to
// support partial updates.
func (o *Org) UpdateOrg(ctx context.Context, req *api.UpdateOrgRequest) (
	*api.Org, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, errPerm(api.Role_SYS_ADMIN)
	}

	if req.GetOrg() == nil {
		return nil, status.Error(codes.InvalidArgument,
			req.Validate().Error())
	}

	// Admins can only update their own org, system admins can update any org.
	if (sess.Role < api.Role_SYS_ADMIN && req.GetOrg().GetId() != sess.OrgID) ||
		(sess.Role < api.Role_ADMIN && req.GetOrg().GetId() == sess.OrgID) {
		return nil, errPerm(api.Role_SYS_ADMIN)
	}

	// Perform partial update if directed.
	if len(req.GetUpdateMask().GetPaths()) > 0 {
		// Normalize and validate field mask.
		req.GetUpdateMask().Normalize()
		if !req.GetUpdateMask().IsValid(req.GetOrg()) {
			return nil, status.Error(codes.InvalidArgument,
				"invalid field mask")
		}

		org, err := o.orgDAO.Read(ctx, req.GetOrg().GetId())
		if err != nil {
			return nil, errToStatus(err)
		}

		fmutils.Filter(req.GetOrg(), req.GetUpdateMask().GetPaths())
		proto.Merge(org, req.GetOrg())
		req.Org = org
	}

	// Validate after merge to support partial updates.
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	org, err := o.orgDAO.Update(ctx, req.GetOrg())
	if err != nil {
		return nil, errToStatus(err)
	}

	return org, nil
}

// DeleteOrg deletes an organization by ID.
func (o *Org) DeleteOrg(ctx context.Context, req *api.DeleteOrgRequest) (
	*emptypb.Empty, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_SYS_ADMIN {
		return nil, errPerm(api.Role_SYS_ADMIN)
	}

	if err := o.orgDAO.Delete(ctx, req.GetId()); err != nil {
		return nil, errToStatus(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs(StatusCodeKey,
		strconv.Itoa(http.StatusNoContent))); err != nil {
		logger := alog.FromContext(ctx)
		logger.Errorf("DeleteOrg grpc.SetHeader: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListOrgs retrieves all organizations.
func (o *Org) ListOrgs(ctx context.Context, req *api.ListOrgsRequest) (
	*api.ListOrgsResponse, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok {
		return nil, errPerm(api.Role_SYS_ADMIN)
	}

	// If the org does not have sufficient role, return only their org.
	if sess.Role < api.Role_SYS_ADMIN {
		org, err := o.orgDAO.Read(ctx, sess.OrgID)
		if err != nil {
			return nil, errToStatus(err)
		}

		return &api.ListOrgsResponse{
			Orgs:      []*api.Org{org},
			TotalSize: 1,
		}, nil
	}

	if req.GetPageSize() == 0 {
		req.PageSize = defaultPageSize
	}

	lBoundTS, prevID, err := session.ParsePageToken(req.GetPageToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid page token")
	}

	// Retrieve PageSize+1 entries to find last page.
	orgs, count, err := o.orgDAO.List(ctx, lBoundTS, prevID, req.GetPageSize()+1)
	if err != nil {
		return nil, errToStatus(err)
	}

	resp := &api.ListOrgsResponse{Orgs: orgs, TotalSize: count}

	// Populate next page token.
	if len(orgs) == int(req.GetPageSize()+1) {
		resp.Orgs = orgs[:len(orgs)-1]

		if resp.NextPageToken, err = session.GeneratePageToken(
			orgs[len(orgs)-2].GetCreatedAt().AsTime(),
			orgs[len(orgs)-2].GetId()); err != nil {
			// GeneratePageToken should not error based on a DB-derived UUID.
			// Log the error and include the usable empty token.
			logger := alog.FromContext(ctx)
			logger.Errorf("ListOrgs session.GeneratePageToken org, err: "+
				"%+v, %v", orgs[len(orgs)-2], err)
		}
	}

	return resp, nil
}
