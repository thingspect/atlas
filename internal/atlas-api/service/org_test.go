//go:build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateOrg(t *testing.T) {
	t.Parallel()

	t.Run("Create valid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")
		retOrg, _ := proto.Clone(org).(*api.Org)

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Create(gomock.Any(), org).Return(retOrg, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		createOrg, err := orgSvc.CreateOrg(ctx, &api.CreateOrgRequest{Org: org})
		t.Logf("org, createOrg, err: %+v, %+v, %v", org, createOrg, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(org, createOrg) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", org, createOrg)
		}
	})

	t.Run("Create org with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		createOrg, err := orgSvc.CreateOrg(ctx, &api.CreateOrgRequest{})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.Nil(t, createOrg)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Create org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		createOrg, err := orgSvc.CreateOrg(ctx, &api.CreateOrgRequest{})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.Nil(t, createOrg)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Create invalid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Create(gomock.Any(), org).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		createOrg, err := orgSvc.CreateOrg(ctx, &api.CreateOrgRequest{Org: org})
		t.Logf("org, createOrg, err: %+v, %+v, %v", org, createOrg, err)
		require.Nil(t, createOrg)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestGetOrg(t *testing.T) {
	t.Parallel()

	t.Run("Get org by valid ID", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")
		retOrg, _ := proto.Clone(org).(*api.Org)

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Read(gomock.Any(), org.Id).Return(retOrg, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		getOrg, err := orgSvc.GetOrg(ctx, &api.GetOrgRequest{Id: org.Id})
		t.Logf("org, getOrg, err: %+v, %+v, %v", org, getOrg, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(org, getOrg) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", org, getOrg)
		}
	})

	t.Run("Get org with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		getOrg, err := orgSvc.GetOrg(ctx, &api.GetOrgRequest{})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.Nil(t, getOrg)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Get org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		getOrg, err := orgSvc.GetOrg(ctx, &api.GetOrgRequest{})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.Nil(t, getOrg)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Get org by unknown ID", func(t *testing.T) {
		t.Parallel()

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Read(gomock.Any(), gomock.Any()).Return(nil,
			dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		getOrg, err := orgSvc.GetOrg(ctx,
			&api.GetOrgRequest{Id: uuid.NewString()})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.Nil(t, getOrg)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestUpdateOrg(t *testing.T) {
	t.Parallel()

	t.Run("Update org by valid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")
		retOrg, _ := proto.Clone(org).(*api.Org)

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Update(gomock.Any(), org).Return(retOrg, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{Org: org})
		t.Logf("org, updateOrg, err: %+v, %+v, %v", org, updateOrg, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(org, updateOrg) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", org, updateOrg)
		}
	})

	t.Run("Partial update org by valid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")
		retOrg, _ := proto.Clone(org).(*api.Org)
		part := &api.Org{Id: org.Id, Name: random.String(10)}
		merged := &api.Org{
			Id: org.Id, Name: part.Name, DisplayName: org.DisplayName,
			Email: org.Email,
		}
		retMerged, _ := proto.Clone(merged).(*api.Org)

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Read(gomock.Any(), org.Id).Return(retOrg, nil).Times(1)
		orger.EXPECT().Update(gomock.Any(), matcher.NewProtoMatcher(merged)).
			Return(retMerged, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{
			Org:        part,
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"name"}},
		})
		t.Logf("merged, updateOrg, err: %+v, %+v, %v", merged, updateOrg, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(merged, updateOrg) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", merged, updateOrg)
		}
	})

	t.Run("Update org with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Update nil org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid UpdateOrgRequest.Org: value is required"), err)
	})

	t.Run("Update same org with insufficient role", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_BUILDER,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{Org: org})
		t.Logf("org, updateOrg, err: %+v, %+v, %v", org, updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Update different org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		updateOrg, err := orgSvc.UpdateOrg(ctx,
			&api.UpdateOrgRequest{Org: random.Org("api-org")})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{
			Org:        org,
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"aaa"}},
		})
		t.Logf("org, updateOrg, err: %+v, %+v, %v", org, updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid field mask"), err)
	})

	t.Run("Partial update org by unknown org", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		part := &api.Org{Id: uuid.NewString(), Name: random.String(10)}

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Read(gomock.Any(), part.Id).Return(nil, dao.ErrNotFound).
			Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{
			Org:        part,
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"name"}},
		})
		t.Logf("part, updateOrg, err: %+v, %+v, %v", part, updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("Update org validation failure", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")
		org.Name = random.String(41)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{Org: org})
		t.Logf("org, updateOrg, err: %+v, %+v, %v", org, updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid "+
			"UpdateOrgRequest.Org: embedded message failed validation | "+
			"caused by: invalid Org.Name: value length must be between 5 and "+
			"40 runes, inclusive"), err)
	})

	t.Run("Update org by invalid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Update(gomock.Any(), org).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		updateOrg, err := orgSvc.UpdateOrg(ctx, &api.UpdateOrgRequest{Org: org})
		t.Logf("org, updateOrg, err: %+v, %+v, %v", org, updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestDeleteOrg(t *testing.T) {
	t.Parallel()

	t.Run("Delete org by valid ID", func(t *testing.T) {
		t.Parallel()

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		_, err := orgSvc.DeleteOrg(ctx,
			&api.DeleteOrgRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete org with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		_, err := orgSvc.DeleteOrg(ctx, &api.DeleteOrgRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Delete org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		_, err := orgSvc.DeleteOrg(ctx, &api.DeleteOrgRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("Delete org by unknown ID", func(t *testing.T) {
		t.Parallel()

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Delete(gomock.Any(), gomock.Any()).
			Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		_, err := orgSvc.DeleteOrg(ctx,
			&api.DeleteOrgRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestListOrgs(t *testing.T) {
	t.Parallel()

	t.Run("List orgs", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		orgs := []*api.Org{
			random.Org("api-org"),
			random.Org("api-org"),
			random.Org("api-org"),
		}

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().List(gomock.Any(), time.Time{}, "", int32(51)).
			Return(orgs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		listOrgs, err := orgSvc.ListOrgs(ctx, &api.ListOrgsRequest{})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listOrgs.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListOrgsResponse{Orgs: orgs, TotalSize: 3},
			listOrgs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListOrgsResponse{Orgs: orgs, TotalSize: 3}, listOrgs)
		}
	})

	t.Run("List orgs with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		orgs := []*api.Org{
			random.Org("api-org"),
			random.Org("api-org"),
			random.Org("api-org"),
		}

		next, err := session.GeneratePageToken(orgs[1].CreatedAt.AsTime(),
			orgs[1].Id)
		require.NoError(t, err)

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().List(gomock.Any(), time.Time{}, "", int32(3)).
			Return(orgs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		listOrgs, err := orgSvc.ListOrgs(ctx, &api.ListOrgsRequest{PageSize: 2})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listOrgs.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListOrgsResponse{
			Orgs: orgs[:2], NextPageToken: next, TotalSize: 3,
		}, listOrgs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListOrgsResponse{
				Orgs: orgs[:2], NextPageToken: next, TotalSize: 3,
			}, listOrgs)
		}
	})

	t.Run("List orgs with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		listOrgs, err := orgSvc.ListOrgs(ctx, &api.ListOrgsRequest{})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.Nil(t, listOrgs)
		require.Equal(t, errPerm(api.Role_SYS_ADMIN), err)
	})

	t.Run("List own org with insufficient role", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Read(gomock.Any(), org.Id).Return(org, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		listOrgs, err := orgSvc.ListOrgs(ctx, &api.ListOrgsRequest{})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListOrgsResponse{
			Orgs: []*api.Org{org}, TotalSize: 1,
		}, listOrgs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListOrgsResponse{
				Orgs: []*api.Org{org}, TotalSize: 1,
			}, listOrgs)
		}
	})

	t.Run("List orgs by unknown ID", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().Read(gomock.Any(), gomock.Any()).Return(nil,
			dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: org.Id, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		listOrgs, err := orgSvc.ListOrgs(ctx, &api.ListOrgsRequest{})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.Nil(t, listOrgs)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})

	t.Run("List orgs by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(nil)
		listOrgs, err := orgSvc.ListOrgs(ctx,
			&api.ListOrgsRequest{PageToken: badUUID})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.Nil(t, listOrgs)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid page token"), err)
	})

	t.Run("List orgs by invalid org ID", func(t *testing.T) {
		t.Parallel()

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil, int32(0), dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: "aaa", Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		listOrgs, err := orgSvc.ListOrgs(ctx, &api.ListOrgsRequest{})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.Nil(t, listOrgs)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("List orgs with generation failure", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		orgs := []*api.Org{
			random.Org("api-org"),
			random.Org("api-org"),
			random.Org("api-org"),
		}
		orgs[1].Id = badUUID

		orger := NewMockOrger(gomock.NewController(t))
		orger.EXPECT().List(gomock.Any(), time.Time{}, "", int32(3)).
			Return(orgs, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_SYS_ADMIN,
			}), testTimeout)
		defer cancel()

		orgSvc := NewOrg(orger)
		listOrgs, err := orgSvc.ListOrgs(ctx, &api.ListOrgsRequest{PageSize: 2})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listOrgs.TotalSize)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListOrgsResponse{Orgs: orgs[:2], TotalSize: 3},
			listOrgs) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.ListOrgsResponse{Orgs: orgs[:2], TotalSize: 3}, listOrgs)
		}
	})
}
