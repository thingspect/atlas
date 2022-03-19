package rule

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createRule = `
INSERT INTO rules (org_id, name, status, device_tag, attr, expr, created_at,
updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
RETURNING id
`

// Create creates a rule in the database.
func (d *DAO) Create(ctx context.Context, rule *api.Rule) (*api.Rule, error) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	rule.CreatedAt = timestamppb.New(now)
	rule.UpdatedAt = timestamppb.New(now)

	if err := d.pg.QueryRowContext(ctx, createRule, rule.OrgId, rule.Name,
		rule.Status.String(), rule.DeviceTag, rule.Attr, rule.Expr, now).Scan(
		&rule.Id); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return rule, nil
}

const readRule = `
SELECT id, org_id, name, status, device_tag, attr, expr, created_at, updated_at
FROM rules
WHERE (id, org_id) = ($1, $2)
`

// Read retrieves a rule by ID and org ID.
func (d *DAO) Read(ctx context.Context, ruleID, orgID string) (
	*api.Rule, error,
) {
	rule := &api.Rule{}
	var status string
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readRule, ruleID, orgID).Scan(&rule.Id,
		&rule.OrgId, &rule.Name, &status, &rule.DeviceTag, &rule.Attr,
		&rule.Expr, &createdAt, &updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	rule.Status = api.Status(api.Status_value[status])
	rule.CreatedAt = timestamppb.New(createdAt)
	rule.UpdatedAt = timestamppb.New(updatedAt)

	return rule, nil
}

const updateRule = `
UPDATE rules
SET name = $1, status = $2, device_tag = $3, attr = $4, expr = $5,
updated_at = $6
WHERE (id, org_id) = ($7, $8)
RETURNING created_at
`

// Update updates a rule in the database. CreatedAt should not update, so it
// is safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, rule *api.Rule) (*api.Rule, error) {
	var createdAt time.Time
	updatedAt := time.Now().UTC().Truncate(time.Microsecond)
	rule.UpdatedAt = timestamppb.New(updatedAt)

	if err := d.pg.QueryRowContext(ctx, updateRule, rule.Name,
		rule.Status.String(), rule.DeviceTag, rule.Attr, rule.Expr, updatedAt,
		rule.Id, rule.OrgId).Scan(&createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	rule.CreatedAt = timestamppb.New(createdAt)

	return rule, nil
}

const deleteRule = `
DELETE FROM rules
WHERE (id, org_id) = ($1, $2)
`

// Delete deletes a rule by ID and org ID.
func (d *DAO) Delete(ctx context.Context, ruleID, orgID string) error {
	// Verify a rule exists before attempting to delete it. Do not remap the
	// error.
	if _, err := d.Read(ctx, ruleID, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteRule, ruleID, orgID)

	return dao.DBToSentinel(err)
}

const countRules = `
SELECT count(*)
FROM rules
WHERE org_id = $1
`

const listRules = `
SELECT id, org_id, name, status, device_tag, attr, expr, created_at, updated_at
FROM rules
WHERE org_id = $1
`

const listRulesTSAndID = `
AND (created_at > $2
OR (created_at = $2
AND id > $3
))
`

const listRulesLimit = `
ORDER BY created_at ASC, id ASC
LIMIT %d
`

// List retrieves all rules by org ID with pagination. If lBoundTS and prevID
// are zero values, the first page of results is returned. Limits of 0 or less
// do not apply a limit. List returns a slice of rules, a total count, and an
// error value.
func (d *DAO) List(
	ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
	limit int32,
) ([]*api.Rule, int32, error) {
	// Run count query.
	var count int32
	if err := d.pg.QueryRowContext(ctx, countRules, orgID).Scan(
		&count); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	// Build list query.
	query := listRules
	args := []interface{}{orgID}

	if prevID != "" && !lBoundTS.IsZero() {
		query += listRulesTSAndID
		args = append(args, lBoundTS, prevID)
	}

	// Ordering is applied with the limit, which will always be present for API
	// usage, whereas lBoundTS and prevID will not for first pages.
	if limit > 0 {
		query += fmt.Sprintf(listRulesLimit, limit)
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

	var rules []*api.Rule
	for rows.Next() {
		rule := &api.Rule{}
		var status string
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&rule.Id, &rule.OrgId, &rule.Name, &status,
			&rule.DeviceTag, &rule.Attr, &rule.Expr, &createdAt,
			&updatedAt); err != nil {
			return nil, 0, dao.DBToSentinel(err)
		}

		rule.Status = api.Status(api.Status_value[status])
		rule.CreatedAt = timestamppb.New(createdAt)
		rule.UpdatedAt = timestamppb.New(updatedAt)
		rules = append(rules, rule)
	}

	if err = rows.Close(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	return rules, count, nil
}

const listByTags = `
SELECT id, org_id, name, status, device_tag, attr, expr, created_at, updated_at
FROM rules
WHERE (org_id, status, attr) = ($1, 'ACTIVE', $2)
AND device_tag = ANY ($3::varchar(255)[])
ORDER BY created_at
`

// ListByTags retrieves all active rules by org ID, attribute, and any matching
// device tags.
func (d *DAO) ListByTags(
	ctx context.Context, orgID string, attr string, deviceTags []string,
) ([]*api.Rule, error) {
	var tags pgtype.VarcharArray
	if err := tags.Set(deviceTags); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	rows, err := d.pg.QueryContext(ctx, listByTags, orgID, attr, tags)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("ListByTags rows.Close: %v", err)
		}
	}()

	var rules []*api.Rule
	for rows.Next() {
		rule := &api.Rule{}
		var status string
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&rule.Id, &rule.OrgId, &rule.Name, &status,
			&rule.DeviceTag, &rule.Attr, &rule.Expr, &createdAt,
			&updatedAt); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		rule.Status = api.Status(api.Status_value[status])
		rule.CreatedAt = timestamppb.New(createdAt)
		rule.UpdatedAt = timestamppb.New(updatedAt)
		rules = append(rules, rule)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return rules, nil
}
