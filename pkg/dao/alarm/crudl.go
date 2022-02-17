package alarm

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

const createAlarm = `
INSERT INTO alarms (org_id, rule_id, name, status, type, user_tags,
subject_template, body_template, repeat_interval, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)
RETURNING id
`

// Create creates an alarm in the database.
func (d *DAO) Create(ctx context.Context, alarm *api.Alarm) (*api.Alarm,
	error) {
	var tags pgtype.VarcharArray
	if err := tags.Set(alarm.UserTags); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	alarm.CreatedAt = timestamppb.New(now)
	alarm.UpdatedAt = timestamppb.New(now)

	if err := d.pg.QueryRowContext(ctx, createAlarm, alarm.OrgId, alarm.RuleId,
		alarm.Name, alarm.Status.String(), alarm.Type.String(), tags,
		alarm.SubjectTemplate, alarm.BodyTemplate, alarm.RepeatInterval,
		now).Scan(&alarm.Id); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return alarm, nil
}

const readAlarm = `
SELECT id, org_id, rule_id, name, status, type, user_tags, subject_template,
body_template, repeat_interval, created_at, updated_at
FROM alarms
WHERE (id, org_id, rule_id) = ($1, $2, $3)
`

// Read retrieves an alarm by ID, org ID, and rule ID.
func (d *DAO) Read(ctx context.Context, alarmID, orgID,
	ruleID string) (*api.Alarm, error) {
	alarm := &api.Alarm{}
	var status, alarmType string
	var tags pgtype.VarcharArray
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readAlarm, alarmID, orgID, ruleID).Scan(
		&alarm.Id, &alarm.OrgId, &alarm.RuleId, &alarm.Name, &status,
		&alarmType, &tags, &alarm.SubjectTemplate, &alarm.BodyTemplate,
		&alarm.RepeatInterval, &createdAt, &updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	alarm.Status = api.Status(api.Status_value[status])
	alarm.Type = api.AlarmType(api.AlarmType_value[alarmType])
	if err := tags.AssignTo(&alarm.UserTags); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	alarm.CreatedAt = timestamppb.New(createdAt)
	alarm.UpdatedAt = timestamppb.New(updatedAt)

	return alarm, nil
}

const updateAlarm = `
UPDATE alarms
SET name = $1, status = $2, type = $3, user_tags = $4, subject_template = $5,
body_template = $6, repeat_interval = $7, updated_at = $8
WHERE (id, org_id, rule_id) = ($9, $10, $11)
RETURNING created_at
`

// Update updates an alarm in the database. CreatedAt should not update, so it
// is safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, alarm *api.Alarm) (*api.Alarm,
	error) {
	var tags pgtype.VarcharArray
	if err := tags.Set(alarm.UserTags); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	var createdAt time.Time
	updatedAt := time.Now().UTC().Truncate(time.Microsecond)
	alarm.UpdatedAt = timestamppb.New(updatedAt)

	if err := d.pg.QueryRowContext(ctx, updateAlarm, alarm.Name,
		alarm.Status.String(), alarm.Type.String(), tags, alarm.SubjectTemplate,
		alarm.BodyTemplate, alarm.RepeatInterval, updatedAt, alarm.Id,
		alarm.OrgId, alarm.RuleId).Scan(&createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	alarm.CreatedAt = timestamppb.New(createdAt)

	return alarm, nil
}

const deleteAlarm = `
DELETE FROM alarms
WHERE (id, org_id, rule_id) = ($1, $2, $3)
`

// Delete deletes an alarm by ID, org ID, and rule ID.
func (d *DAO) Delete(ctx context.Context, alarmID, orgID, ruleID string) error {
	// Verify an alarm exists before attempting to delete it. Do not remap the
	// error.
	if _, err := d.Read(ctx, alarmID, orgID, ruleID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteAlarm, alarmID, orgID, ruleID)

	return dao.DBToSentinel(err)
}

const countAlarms = `
SELECT count(*)
FROM alarms
WHERE org_id = $1
`

const countAlarmsRule = `
AND rule_id = $2
`

const listAlarms = `
SELECT id, org_id, rule_id, name, status, type, user_tags, subject_template,
body_template, repeat_interval, created_at, updated_at
FROM alarms
WHERE org_id = $1
`

const listAlarmsTSAndID = `
AND (created_at > $%d
OR (created_at = $%d
AND id > $%d
))
`

const listAlarmsRule = `
AND rule_id = $%d
`

const listAlarmsLimit = `
ORDER BY created_at ASC, id ASC
LIMIT %d
`

// List retrieves all alarms by org ID with pagination and optional rule filter.
// If lBoundTS and prevID are zero values, the first page of results is
// returned. Limits of 0 or less do not apply a limit. List returns a slice of
// alarms, a total count, and an error value.
func (d *DAO) List(ctx context.Context, orgID string, lBoundTS time.Time,
	prevID string, limit int32, ruleID string) ([]*api.Alarm, int32, error) {
	// Build count query.
	cQuery := countAlarms
	cArgs := []interface{}{orgID}

	if ruleID != "" {
		cQuery += countAlarmsRule
		cArgs = append(cArgs, ruleID)
	}

	// Run count query.
	var count int32
	if err := d.pg.QueryRowContext(ctx, cQuery, cArgs...).Scan(
		&count); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	// Build list query.
	lQuery := listAlarms
	lArgs := []interface{}{orgID}

	if prevID != "" && !lBoundTS.IsZero() {
		lQuery += fmt.Sprintf(listAlarmsTSAndID, 2, 2, 3)
		lArgs = append(lArgs, lBoundTS, prevID)

		if ruleID != "" {
			lQuery += fmt.Sprintf(listAlarmsRule, 4)
			lArgs = append(lArgs, ruleID)
		}
	} else if ruleID != "" {
		lQuery += fmt.Sprintf(listAlarmsRule, 2)
		lArgs = append(lArgs, ruleID)
	}

	// Ordering is applied with the limit, which will always be present for API
	// usage, whereas lBoundTS and prevID will not for first pages.
	if limit > 0 {
		lQuery += fmt.Sprintf(listAlarmsLimit, limit)
	}

	// Run list query.
	rows, err := d.pg.QueryContext(ctx, lQuery, lArgs...)
	if err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("List rows.Close: %v", err)
		}
	}()

	var alarms []*api.Alarm
	for rows.Next() {
		alarm := &api.Alarm{}
		var status, alarmType string
		var tags pgtype.VarcharArray
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&alarm.Id, &alarm.OrgId, &alarm.RuleId, &alarm.Name,
			&status, &alarmType, &tags, &alarm.SubjectTemplate,
			&alarm.BodyTemplate, &alarm.RepeatInterval, &createdAt,
			&updatedAt); err != nil {
			return nil, 0, dao.DBToSentinel(err)
		}

		alarm.Status = api.Status(api.Status_value[status])
		alarm.Type = api.AlarmType(api.AlarmType_value[alarmType])
		if err := tags.AssignTo(&alarm.UserTags); err != nil {
			return nil, 0, dao.DBToSentinel(err)
		}
		alarm.CreatedAt = timestamppb.New(createdAt)
		alarm.UpdatedAt = timestamppb.New(updatedAt)
		alarms = append(alarms, alarm)
	}

	if err = rows.Close(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	return alarms, count, nil
}
