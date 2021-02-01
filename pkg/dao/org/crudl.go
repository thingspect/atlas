package org

import (
	"context"
	"strings"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createOrg = `
INSERT INTO orgs (name, created_at, updated_at)
VALUES ($1, $2, $3)
RETURNING id
`

// Create creates an organization in the database.
func (d *DAO) Create(ctx context.Context, org *api.Org) (*api.Org, error) {
	org.Name = strings.ToLower(org.Name)
	now := time.Now().UTC().Truncate(time.Microsecond)
	org.CreatedAt = timestamppb.New(now)
	org.UpdatedAt = timestamppb.New(now)

	if err := d.pg.QueryRowContext(ctx, createOrg, org.Name, now,
		now).Scan(&org.Id); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return org, nil
}

const readOrg = `
SELECT id, name, created_at, updated_at
FROM orgs
WHERE id = $1
`

// Read retrieves an organization by ID.
func (d *DAO) Read(ctx context.Context, orgID string) (*api.Org, error) {
	org := &api.Org{}
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readOrg, orgID).Scan(&org.Id, &org.Name,
		&createdAt, &updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	org.CreatedAt = timestamppb.New(createdAt)
	org.UpdatedAt = timestamppb.New(updatedAt)
	return org, nil
}

const updateOrg = `
UPDATE orgs
SET name = $1, updated_at = $2
WHERE id = $3
RETURNING created_at
`

// Update updates an organization in the database. CreatedAt should not
// update, so it is safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, org *api.Org) (*api.Org, error) {
	org.Name = strings.ToLower(org.Name)
	var createdAt time.Time
	updatedAt := time.Now().UTC().Truncate(time.Microsecond)
	org.UpdatedAt = timestamppb.New(updatedAt)

	if err := d.pg.QueryRowContext(ctx, updateOrg, org.Name, updatedAt,
		org.Id).Scan(&createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	org.CreatedAt = timestamppb.New(createdAt)
	return org, nil
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
func (d *DAO) List(ctx context.Context) ([]*api.Org, error) {
	var orgs []*api.Org

	rows, err := d.pg.QueryContext(ctx, listOrgs)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("List rows.Close: %v", err)
		}
	}()

	for rows.Next() {
		org := &api.Org{}
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&org.Id, &org.Name, &createdAt,
			&updatedAt); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		org.CreatedAt = timestamppb.New(createdAt)
		org.UpdatedAt = timestamppb.New(updatedAt)
		orgs = append(orgs, org)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return orgs, nil
}
