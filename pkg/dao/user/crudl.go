package user

import (
	"context"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createUser = `
INSERT INTO users (org_id, email, password_hash, is_disabled, created_at,
updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`

// Create creates a user in the database.
func (d *DAO) Create(ctx context.Context, user *api.User,
	passHash []byte) (*api.User, error) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	user.CreatedAt = timestamppb.New(now)
	user.UpdatedAt = timestamppb.New(now)

	if err := d.pg.QueryRowContext(ctx, createUser, user.OrgId, user.Email,
		passHash, user.IsDisabled, now, now).Scan(&user.Id); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return user, nil
}

const readUser = `
SELECT id, org_id, email, password_hash, is_disabled, created_at, updated_at
FROM users
WHERE id = $1
AND org_id = $2
`

// Read retrieves a user and password hash by ID and org ID.
func (d *DAO) Read(ctx context.Context, userID, orgID string) (*api.User,
	[]byte, error) {
	user := &api.User{}
	var passHash []byte
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readUser, userID, orgID).Scan(&user.Id,
		&user.OrgId, &user.Email, &passHash, &user.IsDisabled, &createdAt,
		&updatedAt); err != nil {
		return nil, nil, dao.DBToSentinel(err)
	}

	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return user, passHash, nil
}

const readUserByEmail = `
SELECT u.id, u.org_id, u.email, u.password_hash, u.is_disabled, u.created_at,
u.updated_at
FROM users u
INNER JOIN orgs o ON u.org_id = o.id
WHERE u.email = $1
AND o.name = $2
`

// ReadByEmail retrieves a user and password hash by email and org name.
func (d *DAO) ReadByEmail(ctx context.Context, email,
	orgName string) (*api.User, []byte, error) {
	user := &api.User{}
	var passHash []byte
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readUserByEmail, email, orgName).Scan(
		&user.Id, &user.OrgId, &user.Email, &passHash, &user.IsDisabled,
		&createdAt, &updatedAt); err != nil {
		return nil, nil, dao.DBToSentinel(err)
	}

	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return user, passHash, nil
}

const updateUser = `
UPDATE users
SET email = $1, password_hash = $2, is_disabled = $3, updated_at = $4
WHERE id = $5
AND org_id = $6
RETURNING created_at
`

// Update updates a user in the database. CreatedAt should not update, so it is
// safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, user *api.User,
	passHash []byte) (*api.User, error) {
	var createdAt time.Time
	updatedAt := time.Now().UTC().Truncate(time.Microsecond)
	user.UpdatedAt = timestamppb.New(updatedAt)

	if err := d.pg.QueryRowContext(ctx, updateUser, user.Email, passHash,
		user.IsDisabled, updatedAt, user.Id, user.OrgId).Scan(
		&createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	user.CreatedAt = timestamppb.New(createdAt)
	return user, nil
}

const deleteUser = `
DELETE FROM users
WHERE id = $1
AND org_id = $2
`

// Delete deletes a user by ID and org ID.
func (d *DAO) Delete(ctx context.Context, userID, orgID string) error {
	// Verify a user exists before attempting to delete it. Do not remap the
	// error.
	if _, _, err := d.Read(ctx, userID, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteUser, userID, orgID)
	return dao.DBToSentinel(err)
}

const listUsers = `
SELECT id, org_id, email, is_disabled, created_at, updated_at
FROM users
WHERE org_id = $1
`

// List retrieves all users by org ID. Password hashes are omitted from the
// results.
func (d *DAO) List(ctx context.Context, orgID string) ([]*api.User, error) {
	var users []*api.User

	rows, err := d.pg.QueryContext(ctx, listUsers, orgID)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			alog.Errorf("List rows.Close: %v", err)
		}
	}()

	for rows.Next() {
		user := &api.User{}
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&user.Id, &user.OrgId, &user.Email, &user.IsDisabled,
			&createdAt, &updatedAt); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(updatedAt)
		users = append(users, user)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return users, nil
}
