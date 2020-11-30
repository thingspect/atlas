// +build !unit

package org

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Create valid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		org := Org{Name: random.String(10)}
		createOrg, err := globalDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		require.Equal(t, org.Name, createOrg.Name)
	})

	t.Run("Create invalid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		org := Org{Name: random.String(41)}
		createOrg, err := globalDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.Nil(t, createOrg)
		require.EqualError(t, err, "ERROR: value too long for type character "+
			"varying(40) (SQLSTATE 22001)")
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	org := Org{Name: random.String(10)}
	createOrg, err := globalDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	t.Run("Read org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readOrg, err := globalDAO.Read(ctx, createOrg.ID)
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.NoError(t, err)
		require.Equal(t, createOrg, readOrg)
	})

	t.Run("Read org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readOrg, err := globalDAO.Read(ctx, uuid.New().String())
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.Nil(t, readOrg)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Read org by invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		readOrg, err := globalDAO.Read(ctx, random.String(10))
		t.Logf("readOrg, err: %+v, %v", readOrg, err)
		require.Nil(t, readOrg)
		require.Contains(t, err.Error(),
			"ERROR: invalid input syntax for type uuid")
	})
}

func TestUpdateOrg(t *testing.T) {
	t.Parallel()

	t.Run("Update org by valid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := Org{Name: random.String(10)}
		createOrg, err := globalDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = random.String(10)

		updateOrg, err := globalDAO.Update(ctx, *createOrg)
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
		updateOrg, err := globalDAO.Update(ctx, unknownOrg)
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("Update org by invalid org", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := Org{Name: random.String(10)}
		createOrg, err := globalDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		// Update org fields.
		createOrg.Name = random.String(41)

		updateOrg, err := globalDAO.Update(ctx, *createOrg)
		t.Logf("updateOrg, err: %+v, %v", updateOrg, err)
		require.Nil(t, updateOrg)
		require.EqualError(t, err, "ERROR: value too long for type character "+
			"varying(40) (SQLSTATE 22001)")
	})
}

func TestDeleteOrg(t *testing.T) {
	t.Parallel()

	t.Run("Delete org by valid ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		org := Org{Name: random.String(10)}
		createOrg, err := globalDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)

		err = globalDAO.Delete(ctx, createOrg.ID)
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Read org by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()

			readOrg, err := globalDAO.Read(ctx, createOrg.ID)
			t.Logf("readOrg, err: %+v, %v", readOrg, err)
			require.Equal(t, sql.ErrNoRows, err)
		})
	})

	t.Run("Delete org by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := globalDAO.Delete(ctx, uuid.New().String())
		t.Logf("err: %v", err)
		require.Equal(t, sql.ErrNoRows, err)
	})
}

func TestListOrgs(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	var lastOrgID string
	for i := 0; i < 3; i++ {
		org := Org{Name: random.String(10)}
		createOrg, err := globalDAO.Create(ctx, org)
		t.Logf("createOrg, err: %+v, %v", createOrg, err)
		require.NoError(t, err)
		lastOrgID = createOrg.ID
	}

	t.Run("List orgs", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		listOrgs, err := globalDAO.List(ctx)
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
