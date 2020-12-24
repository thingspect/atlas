package org

import (
	"context"
	"time"

	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
)

const createOrg = `
INSERT INTO orgs (name, created_at, updated_at)
VALUES ($1, $2, $3)
RETURNING id
`

// Create creates an organization in the database.
func (d *DAO) Create(ctx context.Context, org Org) (*Org, error) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	org.CreatedAt = now
	org.UpdatedAt = now

	if err := d.pg.QueryRowContext(ctx, createOrg, org.Name, org.CreatedAt,
		org.UpdatedAt).Scan(&org.ID); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return &org, nil
}

const readOrg = `
SELECT id, name, created_at, updated_at
FROM orgs
WHERE id = $1
`

// Read retrieves an organization by ID.
func (d *DAO) Read(ctx context.Context, orgID string) (*Org, error) {
	var org Org

	if err := d.pg.QueryRowContext(ctx, readOrg, orgID).Scan(&org.ID, &org.Name,
		&org.CreatedAt, &org.UpdatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	org.CreatedAt = org.CreatedAt.In(time.UTC)
	org.UpdatedAt = org.UpdatedAt.In(time.UTC)
	return &org, nil
}

const updateOrg = `
UPDATE orgs
SET name = $1, updated_at = $2
WHERE id = $3
RETURNING created_at
`

// Update updates an organization in the database. CreatedAt should not
// update, so it is safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, org Org) (*Org, error) {
	org.UpdatedAt = time.Now().UTC().Truncate(time.Microsecond)

	if err := d.pg.QueryRowContext(ctx, updateOrg, org.Name, org.UpdatedAt,
		org.ID).Scan(&org.CreatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	org.CreatedAt = org.CreatedAt.In(time.UTC)
	return &org, nil
}

const deleteOrg = `
DELETE FROM orgs
WHERE id = $1
`

// Delete deletes an organization by ID.
func (d *DAO) Delete(ctx context.Context, orgID string) error {
	// Verify an org exists before attempting to delete it. Do not remap the
	// error.
	if _, err := d.Read(ctx, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteOrg, orgID)
	return dao.DBToSentinel(err)
}

const listOrgs = `
SELECT id, name, created_at, updated_at
FROM orgs
`

// List retrieves all organizations.
func (d *DAO) List(ctx context.Context) ([]*Org, error) {
	var orgs []*Org

	rows, err := d.pg.QueryContext(ctx, listOrgs)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			alog.Errorf("List rows.Close: %v", err)
		}
	}()

	for rows.Next() {
		var org Org
		if err = rows.Scan(&org.ID, &org.Name, &org.CreatedAt,
			&org.UpdatedAt); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		org.CreatedAt = org.CreatedAt.In(time.UTC)
		org.UpdatedAt = org.UpdatedAt.In(time.UTC)
		orgs = append(orgs, &org)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return orgs, nil
}
