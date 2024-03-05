//go:build !unit

package test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateOrg(t *testing.T) {
	t.Parallel()

	t.Run("Create valid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: org})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		require.NotEqual(t, org.GetId(), createOrg.GetId())
		require.Equal(t, org.GetName(), createOrg.GetName())
		require.WithinDuration(t, time.Now(), createOrg.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createOrg.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid org with uppercase name", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")
		org.Name = strings.ToUpper(org.GetName())

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: org})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		require.NotEqual(t, org.GetId(), createOrg.GetId())
		require.Equal(t, strings.ToLower(org.GetName()), createOrg.GetName())
		require.WithinDuration(t, time.Now(), createOrg.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createOrg.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create valid org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(globalAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: random.Org("api-org")})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.Nil(t, createOrg)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, SYS_ADMIN role required")
	})

	t.Run("Create invalid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-org")
		org.Name = "api-org-" + random.String(40)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: org})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.Nil(t, createOrg)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateOrgRequest.Org: embedded message failed validation "+
			"| caused by: invalid Org.Name: value length must be between 5 "+
			"and 40 runes, inclusive")
	})
}

func TestGetOrg(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
	createOrg, err := orgCli.CreateOrg(ctx,
		&api.CreateOrgRequest{Org: random.Org("api-org")})
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Get org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		getOrg, err := orgCli.GetOrg(ctx, &api.GetOrgRequest{Id: createOrg.GetId()})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, createOrg, getOrg)
	})

	t.Run("Get org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(globalAdminGRPCConn)
		getOrg, err := orgCli.GetOrg(ctx,
			&api.GetOrgRequest{Id: createOrg.GetId()})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.Nil(t, getOrg)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, SYS_ADMIN role required")
	})

	t.Run("Get org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		getOrg, err := orgCli.GetOrg(ctx,
			&api.GetOrgRequest{Id: uuid.NewString()})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.Nil(t, getOrg)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestUpdateOrg(t *testing.T) {
	t.Parallel()

	t.Run("Update org by valid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: random.Org("api-org")})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = "api-org-" + random.String(10)
		createOrg.DisplayName = "dao-org-" + random.String(10)
		createOrg.Email = "dao-org-" + random.Email()

		updateOrg, err := orgCli.UpdateOrg(ctx,
			&api.UpdateOrgRequest{Org: createOrg})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.NoError(t, err)
		require.Equal(t, createOrg.GetName(), updateOrg.GetName())
		require.Equal(t, createOrg.GetDisplayName(), updateOrg.GetDisplayName())
		require.Equal(t, createOrg.GetEmail(), updateOrg.GetEmail())
		require.True(t, updateOrg.GetUpdatedAt().AsTime().After(
			updateOrg.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createOrg.GetCreatedAt().AsTime(),
			updateOrg.GetUpdatedAt().AsTime(), 2*time.Second)

		getOrg, err := orgCli.GetOrg(ctx,
			&api.GetOrgRequest{Id: createOrg.GetId()})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateOrg, getOrg)
	})

	t.Run("Partial update org by valid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: random.Org("api-org")})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		part := &api.Org{
			Id: createOrg.GetId(), Name: "api-org-" + random.String(10),
			DisplayName: "dao-org-" + random.String(10),
			Email:       "dao-org-" + random.Email(),
		}

		updateOrg, err := orgCli.UpdateOrg(ctx, &api.UpdateOrgRequest{
			Org: part, UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{
				"name", "display_name", "email",
			}},
		})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.NoError(t, err)
		require.Equal(t, part.GetName(), updateOrg.GetName())
		require.Equal(t, part.GetDisplayName(), updateOrg.GetDisplayName())
		require.Equal(t, part.GetEmail(), updateOrg.GetEmail())
		require.True(t, updateOrg.GetUpdatedAt().AsTime().After(
			updateOrg.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createOrg.GetCreatedAt().AsTime(),
			updateOrg.GetUpdatedAt().AsTime(), 2*time.Second)

		getOrg, err := orgCli.GetOrg(ctx,
			&api.GetOrgRequest{Id: createOrg.GetId()})
		t.Logf("getOrg, err: %+v, %v", getOrg, err)
		require.NoError(t, err)
		require.EqualExportedValues(t, updateOrg, getOrg)
	})

	t.Run("Update nil org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		updateOrg, err := orgCli.UpdateOrg(ctx, &api.UpdateOrgRequest{Org: nil})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateOrgRequest.Org: value is required")
	})

	t.Run("Update different org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(globalAdminGRPCConn)
		updateOrg, err := orgCli.UpdateOrg(ctx,
			&api.UpdateOrgRequest{Org: random.Org("api-org")})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, SYS_ADMIN role required")
	})

	t.Run("Partial update invalid field mask", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		updateOrg, err := orgCli.UpdateOrg(ctx, &api.UpdateOrgRequest{
			Org: random.Org("api-org"), UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"aaa"},
			},
		})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid field mask")
	})

	t.Run("Partial update org by unknown org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		updateOrg, err := orgCli.UpdateOrg(ctx, &api.UpdateOrgRequest{
			Org: random.Org("api-org"), UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"name"},
			},
		})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update org by unknown org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		updateOrg, err := orgCli.UpdateOrg(ctx,
			&api.UpdateOrgRequest{Org: random.Org("api-org")})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Update org validation failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: random.Org("api-org")})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = "api-org-" + random.String(40)

		updateOrg, err := orgCli.UpdateOrg(ctx,
			&api.UpdateOrgRequest{Org: createOrg})
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid UpdateOrgRequest.Org: embedded message failed validation "+
			"| caused by: invalid Org.Name: value length must be between 5 "+
			"and 40 runes, inclusive")
	})
}

func TestDeleteOrg(t *testing.T) {
	t.Parallel()

	t.Run("Delete org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: random.Org("api-org")})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		_, err = orgCli.DeleteOrg(ctx, &api.DeleteOrgRequest{Id: createOrg.GetId()})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read org by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
			getOrg, err := orgCli.GetOrg(ctx,
				&api.GetOrgRequest{Id: createOrg.GetId()})
			t.Logf("getOrg, err: %+v, %v", getOrg, err)
			require.Nil(t, getOrg)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete org with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(globalAdminGRPCConn)
		_, err := orgCli.DeleteOrg(ctx,
			&api.DeleteOrgRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied "+
			"desc = permission denied, SYS_ADMIN role required")
	})

	t.Run("Delete org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		_, err := orgCli.DeleteOrg(ctx,
			&api.DeleteOrgRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListOrgs(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	orgIDs := []string{}
	orgNames := []string{}
	for range 3 {
		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		createOrg, err := orgCli.CreateOrg(ctx,
			&api.CreateOrgRequest{Org: random.Org("api-org")})
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		orgIDs = append(orgIDs, createOrg.GetId())
		orgNames = append(orgNames, createOrg.GetName())
	}

	t.Run("List orgs by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		listOrgs, err := orgCli.ListOrgs(ctx,
			&api.ListOrgsRequest{PageSize: 250})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listOrgs.GetOrgs()), 3)
		require.GreaterOrEqual(t, listOrgs.GetTotalSize(), int32(3))

		var found bool
		for _, org := range listOrgs.GetOrgs() {
			if org.GetId() == orgIDs[len(orgIDs)-1] &&
				org.GetName() == orgNames[len(orgNames)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List orgs by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		listOrgs, err := orgCli.ListOrgs(ctx, &api.ListOrgsRequest{PageSize: 2})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.Len(t, listOrgs.GetOrgs(), 2)
		require.NotEmpty(t, listOrgs.GetNextPageToken())
		require.GreaterOrEqual(t, listOrgs.GetTotalSize(), int32(3))

		nextOrgs, err := orgCli.ListOrgs(ctx, &api.ListOrgsRequest{
			PageSize: 2, PageToken: listOrgs.GetNextPageToken(),
		})
		t.Logf("nextOrgs, err: %+v, %v", nextOrgs, err)
		require.NoError(t, err)
		require.NotEmpty(t, nextOrgs.GetOrgs())
		require.GreaterOrEqual(t, nextOrgs.GetTotalSize(), int32(3))
	})

	t.Run("List orgs with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewOrgServiceClient(globalAdminGRPCConn)
		listOrgs, err := secCli.ListOrgs(ctx, &api.ListOrgsRequest{})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.Len(t, listOrgs.GetOrgs(), 1)
		require.Equal(t, int32(1), listOrgs.GetTotalSize())
	})

	t.Run("List orgs by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		orgCli := api.NewOrgServiceClient(secondarySysAdminGRPCConn)
		listOrgs, err := orgCli.ListOrgs(ctx,
			&api.ListOrgsRequest{PageToken: badUUID})
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.Nil(t, listOrgs)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}
