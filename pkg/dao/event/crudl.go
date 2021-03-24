package event

import (
	"context"
	"strings"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createEvent = `
INSERT INTO events (org_id, uniq_id, rule_id, created_at, trace_id)
VALUES ($1, $2, $3, $4, $5)
`

// Create creates an event in the database. Events are retrieved elsewhere in
// bulk, so only an error value is returned.
func (d *DAO) Create(ctx context.Context, event *api.Event) error {
	// Truncate timestamp to milliseconds for deduplication.
	createdAt := event.CreatedAt.AsTime().UTC().Truncate(time.Millisecond)

	_, err := d.pg.ExecContext(ctx, createEvent, event.OrgId,
		strings.ToLower(event.UniqId), event.RuleId, createdAt, event.TraceId)

	return dao.DBToSentinel(err)
}

const listEventsByUniqID = `
SELECT e.org_id, e.uniq_id, e.rule_id, e.created_at, e.trace_id
FROM events e
WHERE (e.org_id, e.uniq_id) = ($1, $2)
AND e.created_at <= $3
AND e.created_at > $4
`

const listEventsByDevID = `
SELECT e.org_id, e.uniq_id, e.rule_id, e.created_at, e.trace_id
FROM events e
INNER JOIN devices d ON (e.org_id, e.uniq_id) = (d.org_id, d.uniq_id)
WHERE (e.org_id, d.id) = ($1, $2)
AND e.created_at <= $3
AND e.created_at > $4
`

const listEventsRuleID = `
AND e.rule_id = $5
`

const listEventsOrder = `
ORDER BY e.created_at DESC
`

// List retrieves all events by org ID, UniqID or device ID, optional rule ID,
// and [end, start) times. If both uniqID and devID are provided, uniqID takes
// precedence and devID is ignored.
func (d *DAO) List(ctx context.Context, orgID, uniqID, devID, ruleID string,
	end, start time.Time) ([]*api.Event, error) {
	// Build list query.
	query := listEventsByUniqID
	args := []interface{}{orgID}

	if uniqID == "" && devID != "" {
		query = listEventsByDevID
		args = append(args, devID, end, start)
	} else {
		args = append(args, uniqID, end, start)
	}

	if ruleID != "" {
		query += listEventsRuleID
		args = append(args, ruleID)
	}

	query += listEventsOrder

	// Run list query.
	rows, err := d.pg.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("List rows.Close: %v", err)
		}
	}()

	var events []*api.Event
	for rows.Next() {
		event := &api.Event{}
		var createdAt time.Time

		if err = rows.Scan(&event.OrgId, &event.UniqId, &event.RuleId,
			&createdAt, &event.TraceId); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		event.CreatedAt = timestamppb.New(createdAt)
		events = append(events, event)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return events, nil
}

const latestEvents = `
SELECT
  e.org_id,
  e.uniq_id,
  e.rule_id,
  e.created_at,
  e.trace_id
FROM
  EVENTS e
  INNER JOIN (
    SELECT
      org_id,
      uniq_id,
      MAX(created_at) AS created_at
    FROM
      EVENTS
    WHERE
      org_id = $1
    GROUP BY
      org_id,
      uniq_id
  ) m ON (e.org_id, e.uniq_id, e.created_at) = (
    m.org_id, m.uniq_id, m.created_at)
`

const latestEventsRuleID = `
AND e.rule_id = $2
`

const latestEventsOrder = `
ORDER BY e.created_at DESC
`

// Latest retrieves the latest events for each of an organization's devices by
// org ID and optional rule ID.
func (d *DAO) Latest(ctx context.Context, orgID, ruleID string) ([]*api.Event,
	error) {
	// Build latest query.
	query := latestEvents
	args := []interface{}{orgID}

	if ruleID != "" {
		query += latestEventsRuleID
		args = append(args, ruleID)
	}

	query += latestEventsOrder

	// Run latest query.
	rows, err := d.pg.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("Latest rows.Close: %v", err)
		}
	}()

	var events []*api.Event
	for rows.Next() {
		event := &api.Event{}
		var createdAt time.Time

		if err = rows.Scan(&event.OrgId, &event.UniqId, &event.RuleId,
			&createdAt, &event.TraceId); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		event.CreatedAt = timestamppb.New(createdAt)
		events = append(events, event)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return events, nil
}
