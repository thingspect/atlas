//go:build !unit

package org

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
)

const testTimeout = 8 * time.Second

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("dao-org")
		createOrg, _ := proto.Clone(org).(*api.Org)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, createOrg)
		t.Logf("org, createOrg, err: %+v, %+v, %v", org, createOrg, err)
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

		org := random.Org("dao-org")
		org.Name = strings.ToUpper(org.GetName())
		createOrg, _ := proto.Clone(org).(*api.Org)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, createOrg)
		t.Logf("org, createOrg, err: %+v, %+v, %v", org, createOrg, err)
		require.NoError(t, err)
		require.NotEqual(t, org.GetId(), createOrg.GetId())
		require.Equal(t, strings.ToLower(org.GetName()), createOrg.GetName())
		require.WithinDuration(t, time.Now(), createOrg.GetCreatedAt().AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createOrg.GetUpdatedAt().AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid org", func(t *testing.T) {
		t.Parallel()

		org := random.Org("dao-org")
		org.Name = "dao-org-" + random.String(40)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("org, createOrg, err: %+v, %+v, %v", org, createOrg, err)
		require.Nil(t, createOrg)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-org"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Read org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readOrg, err := globalOrgDAO.Read(ctx, createOrg.GetId())
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.NoError(t, err)
		require.Equal(t, createOrg, readOrg)
	})

	t.Run("Read org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readOrg, err := globalOrgDAO.Read(ctx, uuid.NewString())
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.Nil(t, readOrg)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read org by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		readOrg, err := globalOrgDAO.Read(ctx, random.String(10))
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.Nil(t, readOrg)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("Update org by valid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-org"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = "dao-org-" + random.String(10)
		createOrg.DisplayName = "dao-org-" + random.String(10)
		createOrg.Email = "dao-org-" + random.Email()
		updateOrg, _ := proto.Clone(createOrg).(*api.Org)

		updateOrg, err = globalOrgDAO.Update(ctx, updateOrg)
		t.Logf("createOrg, updateOrg, err: %+v, %+v, %v", createOrg, updateOrg,
			err)
		require.NoError(t, err)
		require.Equal(t, createOrg.GetName(), updateOrg.GetName())
		require.Equal(t, createOrg.GetDisplayName(), updateOrg.GetDisplayName())
		require.Equal(t, createOrg.GetEmail(), updateOrg.GetEmail())
		require.True(t, updateOrg.GetUpdatedAt().AsTime().After(
			updateOrg.GetCreatedAt().AsTime()))
		require.WithinDuration(t, createOrg.GetCreatedAt().AsTime(),
			updateOrg.GetUpdatedAt().AsTime(), 2*time.Second)

		readOrg, err := globalOrgDAO.Read(ctx, createOrg.GetId())
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.NoError(t, err)
		require.Equal(t, updateOrg, readOrg)
	})

	t.Run("Update unknown org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		updateOrg, err := globalOrgDAO.Update(ctx, random.Org("dao-org"))
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update org by invalid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-org"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = "dao-org-" + random.String(40)
		updateOrg, _ := proto.Clone(createOrg).(*api.Org)

		updateOrg, err = globalOrgDAO.Update(ctx, updateOrg)
		t.Logf("createOrg, updateOrg, err: %+v, %+v, %v", createOrg, updateOrg,
			err)
		require.Nil(t, updateOrg)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("Delete org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-org"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		err = globalOrgDAO.Delete(ctx, createOrg.GetId())
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read org by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readOrg, err := globalOrgDAO.Read(ctx, createOrg.GetId())
			t.Logf("readOrg, err: %+v, %v", readOrg, err)
			require.Nil(t, readOrg)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := globalOrgDAO.Delete(ctx, uuid.NewString())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	orgIDs := []string{}
	orgNames := []string{}
	orgEmails := []string{}
	orgTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-org"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		orgIDs = append(orgIDs, createOrg.GetId())
		orgNames = append(orgNames, createOrg.GetName())
		orgEmails = append(orgEmails, createOrg.GetEmail())
		orgTSes = append(orgTSes, createOrg.GetCreatedAt().AsTime())
	}

	t.Run("List orgs", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listOrgs, listCount, err := globalOrgDAO.List(ctx, time.Time{}, "", 0)
		t.Logf("listOrgs, listCount, err: %+v, %v, %v", listOrgs, listCount,
			err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listOrgs), 3)
		require.GreaterOrEqual(t, listCount, int32(3))

		var found bool
		for _, org := range listOrgs {
			if org.GetId() == orgIDs[len(orgIDs)-1] &&
				org.GetName() == orgNames[len(orgNames)-1] &&
				org.GetEmail() == orgEmails[len(orgEmails)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List orgs with pagination", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listOrgs, listCount, err := globalOrgDAO.List(ctx, orgTSes[1],
			orgIDs[1], 1)
		t.Logf("listOrgs, listCount, err: %+v, %v, %v", listOrgs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listOrgs, 1)
		require.GreaterOrEqual(t, listCount, int32(3))
	})

	t.Run("List orgs with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listOrgs, listCount, err := globalOrgDAO.List(ctx, time.Time{}, "", 2)
		t.Logf("listOrgs, listCount, err: %+v, %v, %v", listOrgs, listCount,
			err)
		require.NoError(t, err)
		require.Len(t, listOrgs, 2)
		require.GreaterOrEqual(t, listCount, int32(3))
	})
}
