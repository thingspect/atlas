package user

import (
	"context"
	"fmt"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createUser = `
INSERT INTO users (org_id, email, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

// Create creates a user in the database.
func (d *DAO) Create(ctx context.Context, user *api.User) (*api.User, error) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	user.CreatedAt = timestamppb.New(now)
	user.UpdatedAt = timestamppb.New(now)

	if err := d.pg.QueryRowContext(ctx, createUser, user.OrgId, user.Email,
		user.Status.String(), now, now).Scan(&user.Id); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return user, nil
}

const readUser = `
SELECT id, org_id, email, status, created_at, updated_at
FROM users
WHERE (id, org_id) = ($1, $2)
`

// Read retrieves a user by ID and org ID.
func (d *DAO) Read(ctx context.Context, userID, orgID string) (*api.User,
	error) {
	user := &api.User{}
	var status string
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readUser, userID, orgID).Scan(&user.Id,
		&user.OrgId, &user.Email, &status, &createdAt, &updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	user.Status = common.Status(common.Status_value[status])
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return user, nil
}

const readUserByEmail = `
SELECT u.id, u.org_id, u.email, u.password_hash, u.status, u.created_at,
u.updated_at
FROM users u
INNER JOIN orgs o ON u.org_id = o.id
WHERE (u.email, o.name) = ($1, $2)
`

// ReadByEmail retrieves a user and password hash by email and org name.
func (d *DAO) ReadByEmail(ctx context.Context, email,
	orgName string) (*api.User, []byte, error) {
	user := &api.User{}
	var passHash []byte
	var status string
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readUserByEmail, email, orgName).Scan(
		&user.Id, &user.OrgId, &user.Email, &passHash, &status, &createdAt,
		&updatedAt); err != nil {
		return nil, nil, dao.DBToSentinel(err)
	}

	user.Status = common.Status(common.Status_value[status])
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return user, passHash, nil
}

const updateUser = `
UPDATE users
SET email = $1, status = $2, updated_at = $3
WHERE (id, org_id) = ($4, $5)
RETURNING created_at
`

// Update updates a user in the database. CreatedAt should not update, so it is
// safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, user *api.User) (*api.User, error) {
	var createdAt time.Time
	updatedAt := time.Now().UTC().Truncate(time.Microsecond)
	user.UpdatedAt = timestamppb.New(updatedAt)

	if err := d.pg.QueryRowContext(ctx, updateUser, user.Email,
		user.Status.String(), updatedAt, user.Id, user.OrgId).Scan(
		&createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	user.CreatedAt = timestamppb.New(createdAt)
	return user, nil
}

//#nosec G101 // false positive for hardcoded credentials
const updateUserPassword = `
UPDATE users
SET password_hash = $1, updated_at = $2
WHERE (id, org_id) = ($3, $4)
`

// UpdatePassword updates a user's password by ID and org ID.
func (d *DAO) UpdatePassword(ctx context.Context, userID, orgID string,
	passHash []byte) error {
	// Verify a user exists before attempting to update it. Do not remap the
	// error.
	if _, err := d.Read(ctx, userID, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, updateUserPassword, passHash,
		time.Now().UTC().Truncate(time.Microsecond), userID, orgID)
	return dao.DBToSentinel(err)
}

const deleteUser = `
DELETE FROM users
WHERE (id, org_id) = ($1, $2)
`

// Delete deletes a user by ID and org ID.
func (d *DAO) Delete(ctx context.Context, userID, orgID string) error {
	// Verify a user exists before attempting to delete it. Do not remap the
	// error.
	if _, err := d.Read(ctx, userID, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteUser, userID, orgID)
	return dao.DBToSentinel(err)
}

const countUsers = `
SELECT count(*)
FROM users
WHERE org_id = $1
`

const listUsers = `
SELECT id, org_id, email, status, created_at, updated_at
FROM users
WHERE org_id = $1
`

const listUsersTSAndID = `
AND (created_at > $2
OR (created_at = $2
AND id > $3
))
`

const listUsersLimit = `
ORDER BY created_at ASC, id ASC
LIMIT %d
`

// List retrieves all users by org ID. If lboundTS and prevID are zero values,
// the first page of results is returned. Limits of 0 or less do not apply a
// limit. List returns a slice of users, a total count, and an error value.
func (d *DAO) List(ctx context.Context, orgID string, lboundTS time.Time,
	prevID string, limit int32) ([]*api.User, int32, error) {
	// Run count query.
	var count int32
	if err := d.pg.QueryRowContext(ctx, countUsers, orgID).Scan(
		&count); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	// Build list query.
	query := listUsers
	args := []interface{}{orgID}

	if prevID != "" && !lboundTS.IsZero() {
		query += listUsersTSAndID
		args = append(args, lboundTS, prevID)
	}

	// Ordering is applied with the limit, which will always be present for API
	// usage, whereas lboundTS and prevID will not for first pages.
	if limit > 0 {
		query += fmt.Sprintf(listUsersLimit, limit)
	}

	// Run list query.
	var users []*api.User
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

	for rows.Next() {
		user := &api.User{}
		var status string
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&user.Id, &user.OrgId, &user.Email, &status,
			&createdAt, &updatedAt); err != nil {
			return nil, 0, dao.DBToSentinel(err)
		}

		user.Status = common.Status(common.Status_value[status])
		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(updatedAt)
		users = append(users, user)
	}

	if err = rows.Close(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	return users, count, nil
}
