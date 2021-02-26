package org

import (
	"context"
	"fmt"
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

const countOrgs = `
SELECT count(*)
FROM orgs
`

const listOrgs = `
SELECT id, name, created_at, updated_at
FROM orgs
`

const listOrgsTSAndID = `
WHERE (created_at > $1
OR (created_at = $1
AND id > $2
))
`

const listOrgsLimit = `
ORDER BY created_at ASC, id ASC
LIMIT %d
`

// List retrieves all organizations. If lboundTS and prevID are zero values,
// the first page of results is returned. Limits of 0 or less do not apply a
// limit. List returns a slice of users, a total count, and an error value.
func (d *DAO) List(ctx context.Context, lboundTS time.Time, prevID string,
	limit int32) ([]*api.Org, int32, error) {
	// Run count query.
	var count int32
	if err := d.pg.QueryRowContext(ctx, countOrgs).Scan(&count); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	// Build list query.
	query := listOrgs
	args := []interface{}{}

	if prevID != "" && !lboundTS.IsZero() {
		query += listOrgsTSAndID
		args = append(args, lboundTS, prevID)
	}

	// Ordering is applied with the limit, which will always be present for API
	// usage, whereas lboundTS and prevID will not for first pages.
	if limit > 0 {
		query += fmt.Sprintf(listOrgsLimit, limit)
	}

	// Run list query.
	rows, err := d.pg.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("List rows.Close: %v", err)
		}
	}()

	var orgs []*api.Org
	for rows.Next() {
		org := &api.Org{}
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&org.Id, &org.Name, &createdAt,
			&updatedAt); err != nil {
			return nil, 0, dao.DBToSentinel(err)
		}

		org.CreatedAt = timestamppb.New(createdAt)
		org.UpdatedAt = timestamppb.New(updatedAt)
		orgs = append(orgs, org)
	}

	if err = rows.Close(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	return orgs, count, nil
}
