package key

import (
	"context"
	"fmt"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createKey = `
INSERT INTO keys (org_id, name, role, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id
`

// Create creates an API key in the database.
func (d *DAO) Create(ctx context.Context, key *api.Key) (*api.Key, error) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	key.CreatedAt = timestamppb.New(now)

	if err := d.pg.QueryRowContext(ctx, createKey, key.GetOrgId(), key.GetName(),
		key.GetRole().String(), now).Scan(&key.Id); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return key, nil
}

const readKey = `
SELECT id, org_id, name, role, created_at
FROM keys
WHERE (id, org_id) = ($1, $2)
`

// read retrieves an API key by ID and org ID.
func (d *DAO) read(ctx context.Context, keyID, orgID string) (*api.Key, error) {
	key := &api.Key{}
	var role string
	var createdAt time.Time

	if err := d.pg.QueryRowContext(ctx, readKey, keyID, orgID).Scan(&key.Id,
		&key.OrgId, &key.Name, &role, &createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	key.Role = api.Role(api.Role_value[role])
	key.CreatedAt = timestamppb.New(createdAt)

	return key, nil
}

const deleteKey = `
DELETE FROM keys
WHERE (id, org_id) = ($1, $2)
`

// Delete deletes an API key by ID and org ID.
func (d *DAO) Delete(ctx context.Context, keyID, orgID string) error {
	// Verify a key exists before attempting to delete it. Do not remap the
	// error.
	if _, err := d.read(ctx, keyID, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteKey, keyID, orgID)

	return dao.DBToSentinel(err)
}

const countKeys = `
SELECT count(*)
FROM keys
WHERE org_id = $1
`

const listKeys = `
SELECT id, org_id, name, role, created_at
FROM keys
WHERE org_id = $1
`

const listKeysTSAndID = `
AND (created_at > $2
OR (created_at = $2
AND id > $3
))
`

const listKeysLimit = `
ORDER BY created_at ASC, id ASC
LIMIT %d
`

// List retrieves all API keys by org ID with pagination. If lBoundTS and prevID
// are zero values, the first page of results is returned. Limits of 0 or less
// do not apply a limit. List returns a slice of keys, a total count, and an
// error value.
func (d *DAO) List(
	ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
	limit int32,
) ([]*api.Key, int32, error) {
	// Run count query.
	var count int32
	if err := d.pg.QueryRowContext(ctx, countKeys, orgID).Scan(
		&count); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	// Build list query.
	query := listKeys
	args := []interface{}{orgID}

	if prevID != "" && !lBoundTS.IsZero() {
		query += listKeysTSAndID
		args = append(args, lBoundTS, prevID)
	}

	// Ordering is applied with the limit, which will always be present for API
	// usage, whereas lBoundTS and prevID will not for first pages.
	if limit > 0 {
		query += fmt.Sprintf(listKeysLimit, limit)
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

	var keys []*api.Key
	for rows.Next() {
		key := &api.Key{}
		var role string
		var createdAt time.Time

		if err = rows.Scan(&key.Id, &key.OrgId, &key.Name, &role,
			&createdAt); err != nil {
			return nil, 0, dao.DBToSentinel(err)
		}

		key.Role = api.Role(api.Role_value[role])
		key.CreatedAt = timestamppb.New(createdAt)
		keys = append(keys, key)
	}

	if err = rows.Close(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	return keys, count, nil
}
