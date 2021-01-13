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

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		org := Org{Name: "dao-org-" + random.String(10)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		require.Equal(t, org.Name, createOrg.Name)
		require.WithinDuration(t, time.Now(), createOrg.CreatedAt,
			2*time.Second)
		require.WithinDuration(t, time.Now(), createOrg.UpdatedAt,
			2*time.Second)
	})

	t.Run("Create valid org with uppercase name", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		org := Org{Name: strings.ToUpper("dao-org-" + random.String(10))}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		require.Equal(t, strings.ToLower(org.Name), createOrg.Name)
		require.WithinDuration(t, time.Now(), createOrg.CreatedAt,
			2*time.Second)
		require.WithinDuration(t, time.Now(), createOrg.UpdatedAt,
			2*time.Second)
	})

	t.Run("Create invalid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		org := Org{Name: "dao-org-" + random.String(40)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.Nil(t, createOrg)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := Org{Name: "dao-org-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Read org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readOrg, err := globalOrgDAO.Read(ctx, createOrg.ID)
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.NoError(t, err)
		require.Equal(t, createOrg, readOrg)
	})

	t.Run("Read org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readOrg, err := globalOrgDAO.Read(ctx, uuid.New().String())
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.Nil(t, readOrg)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Read org by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := Org{Name: "dao-org-" + random.String(10)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = "dao-org-" + random.String(10)

		updateOrg, err := globalOrgDAO.Update(ctx, *createOrg)
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.NoError(t, err)
		require.Equal(t, createOrg.Name, updateOrg.Name)
		require.Equal(t, createOrg.CreatedAt, updateOrg.CreatedAt)
		require.True(t, updateOrg.UpdatedAt.After(updateOrg.CreatedAt))
		require.WithinDuration(t, createOrg.CreatedAt, updateOrg.UpdatedAt,
			2*time.Second)
	})

	t.Run("Update unknown org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		unknownOrg := Org{ID: uuid.New().String()}
		updateOrg, err := globalOrgDAO.Update(ctx, unknownOrg)
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, dao.ErrNotFound, err)
	})

	t.Run("Update org by invalid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := Org{Name: "dao-org-" + random.String(10)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = "dao-org-" + random.String(40)

		updateOrg, err := globalOrgDAO.Update(ctx, *createOrg)
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.ErrorIs(t, err, dao.ErrInvalidFormat)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("Delete org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := Org{Name: "dao-org-" + random.String(10)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		err = globalOrgDAO.Delete(ctx, createOrg.ID)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read org by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			readOrg, err := globalOrgDAO.Read(ctx, createOrg.ID)
			t.Logf("readOrg, err: %+v, %v", readOrg, err)
			require.Nil(t, readOrg)
			require.Equal(t, dao.ErrNotFound, err)
		})
	})

	t.Run("Delete org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalOrgDAO.Delete(ctx, uuid.New().String())
		t.Logf("err: %v", err)
		require.Equal(t, dao.ErrNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	var lastOrgID string
	for i := 0; i < 3; i++ {
		org := Org{Name: "dao-org-" + random.String(10)}
		createOrg, err := globalOrgDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		lastOrgID = createOrg.ID
	}

	t.Run("List orgs", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listOrgs, err := globalOrgDAO.List(ctx)
		t.Logf("listOrgs, err: %+v, %v", listOrgs, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listOrgs), 3)

		var found bool
		for _, org := range listOrgs {
			if org.ID == lastOrgID {
				found = true
			}
		}
		require.True(t, found)
	})
}
