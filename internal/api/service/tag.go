package service

//go:generate mockgen -source tag.go -destination mock_tager_test.go -package service

import (
	"context"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
)

// Tager defines the methods provided by a tag.DAO.
type Tager interface {
	List(ctx context.Context, orgID string) ([]string, error)
}

// Tag service contains functions to query tags.
type Tag struct {
	api.UnimplementedTagServiceServer

	tagDAO Tager
}

// NewTag instantiates and returns a new Tag service.
func NewTag(tagDAO Tager) *Tag {
	return &Tag{
		tagDAO: tagDAO,
	}
}

// ListTags retrieves all tags.
func (d *Tag) ListTags(ctx context.Context,
	req *api.ListTagsRequest) (*api.ListTagsResponse, error) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < common.Role_VIEWER {
		return nil, errPerm(common.Role_VIEWER)
	}

	tags, err := d.tagDAO.List(ctx, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.ListTagsResponse{Tags: tags}, nil
}
