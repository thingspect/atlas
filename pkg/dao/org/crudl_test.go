// +build !unit

package org

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		org := random.Org("dao-org")
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		require.Equal(t, org.Name, createOrg.Name)
		require.WithinDuration(t, time.Now(), createOrg.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createOrg.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create valid org with uppercase name", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		org := random.Org("dao-org")
		org.Name = strings.ToUpper(org.Name)
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		require.Equal(t, strings.ToLower(org.Name), createOrg.Name)
		require.WithinDuration(t, time.Now(), createOrg.CreatedAt.AsTime(),
			2*time.Second)
		require.WithinDuration(t, time.Now(), createOrg.UpdatedAt.AsTime(),
			2*time.Second)
	})

	t.Run("Create invalid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		org := random.Org("dao-org")
		org.Name = "dao-org-" + random.String(40)
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
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

		readOrg, err := globalOrgDAO.Read(ctx, createOrg.Id)
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

		updateOrg, err := globalOrgDAO.Update(ctx, createOrg)
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.NoError(t, err)
		require.Equal(t, createOrg.Name, updateOrg.Name)
		require.Equal(t, createOrg.CreatedAt, updateOrg.CreatedAt)
		require.True(t, updateOrg.UpdatedAt.AsTime().After(
			updateOrg.CreatedAt.AsTime()))
		require.WithinDuration(t, createOrg.CreatedAt.AsTime(),
			updateOrg.UpdatedAt.AsTime(), 2*time.Second)
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

		updateOrg, err := globalOrgDAO.Update(ctx, createOrg)
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
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

		err = globalOrgDAO.Delete(ctx, createOrg.Id)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read org by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			readOrg, err := globalOrgDAO.Read(ctx, createOrg.Id)
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
	orgTSes := []time.Time{}
	for i := 0; i < 3; i++ {
		createOrg, err := globalOrgDAO.Create(ctx, random.Org("dao-org"))
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		orgIDs = append(orgIDs, createOrg.Id)
		orgNames = append(orgNames, createOrg.Name)
		orgTSes = append(orgTSes, createOrg.CreatedAt.AsTime())
	}

	t.Run("List orgs", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listOrgs, listCount, err := globalOrgDAO.List(ctx, time.Time{}, "", 0)
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listOrgs), 3)
		require.GreaterOrEqual(t, listCount, int32(3))

		var found bool
		for _, org := range listOrgs {
			if org.Id == orgIDs[len(orgIDs)-1] &&
				org.Name == orgNames[len(orgIDs)-1] {
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
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.Len(t, listOrgs, 1)
		require.GreaterOrEqual(t, listCount, int32(3))
	})

	t.Run("List orgs with limit", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		listOrgs, listCount, err := globalOrgDAO.List(ctx, time.Time{}, "", 2)
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.Len(t, listOrgs, 2)
		require.GreaterOrEqual(t, listCount, int32(3))
	})
}
