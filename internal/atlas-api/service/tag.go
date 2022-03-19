package service

//go:generate mockgen -source tag.go -destination mock_tagger_test.go -package service

import (
	"context"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/atlas-api/session"
)

// Tagger defines the methods provided by a tag.DAO.
type Tagger interface {
	List(ctx context.Context, orgID string) ([]string, error)
}

// Tag service contains functions to query tags.
type Tag struct {
	api.UnimplementedTagServiceServer

	tagDAO Tagger
}

// NewTag instantiates and returns a new Tag service.
func NewTag(tagDAO Tagger) *Tag {
	return &Tag{
		tagDAO: tagDAO,
	}
}

// ListTags retrieves all tags.
func (t *Tag) ListTags(ctx context.Context, req *api.ListTagsRequest) (
	*api.ListTagsResponse, error,
) {
	sess, ok := session.FromContext(ctx)
	if !ok || sess.Role < api.Role_VIEWER {
		return nil, errPerm(api.Role_VIEWER)
	}

	tags, err := t.tagDAO.List(ctx, sess.OrgID)
	if err != nil {
		return nil, errToStatus(err)
	}

	return &api.ListTagsResponse{Tags: tags}, nil
}
