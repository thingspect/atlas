package alert

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createAlert = `
INSERT INTO alerts (org_id, uniq_id, alarm_id, user_id, status, error,
created_at, trace_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

// Create creates an alert in the database. Alerts are retrieved elsewhere in
// bulk, so only an error value is returned.
func (d *DAO) Create(ctx context.Context, alert *api.Alert) error {
	now := time.Now().UTC().Truncate(time.Microsecond)
	alert.CreatedAt = timestamppb.New(now)

	_, err := d.rw.ExecContext(ctx, createAlert, alert.GetOrgId(),
		strings.ToLower(alert.GetUniqId()), alert.GetAlarmId(),
		alert.GetUserId(), alert.GetStatus().String(), alert.GetError(), now,
		alert.GetTraceId())

	return dao.DBToSentinel(err)
}

const listAlerts = `
SELECT e.org_id, e.uniq_id, e.alarm_id, e.user_id, e.status, e.error,
e.created_at, e.trace_id
FROM alerts e
WHERE e.org_id = $1
AND e.created_at <= $2
AND e.created_at > $3
`

const listAlertsByUniqID = `
SELECT e.org_id, e.uniq_id, e.alarm_id, e.user_id, e.status, e.error,
e.created_at, e.trace_id
FROM alerts e
WHERE (e.org_id, e.uniq_id) = ($1, $2)
AND e.created_at <= $3
AND e.created_at > $4
`

const listAlertsByDevID = `
SELECT e.org_id, e.uniq_id, e.alarm_id, e.user_id, e.status, e.error,
e.created_at, e.trace_id
FROM alerts e
INNER JOIN devices d ON (e.org_id, e.uniq_id) = (d.org_id, d.uniq_id)
WHERE (e.org_id, d.id) = ($1, $2)
AND e.created_at <= $3
AND e.created_at > $4
`

const listAlertsAlarmID = `
AND e.alarm_id = $%d
`

const listAlertsUserID = `
AND e.user_id = $%d
`

const listAlertsOrder = `
ORDER BY e.created_at DESC
`

// List retrieves all alerts by org ID, [end, start) times, and any of the
// following: UniqID, device ID, alarm ID, and user ID. If both uniqID and devID
// are provided, uniqID takes precedence and devID is ignored.
func (d *DAO) List(
	ctx context.Context, orgID, uniqID, devID, alarmID, userID string, end,
	start time.Time,
) ([]*api.Alert, error) {
	// Build list query.
	var query string
	args := []any{orgID}
	var devIncr int

	switch {
	case uniqID == "" && devID != "":
		query = listAlertsByDevID
		args = append(args, devID, end, start)
		devIncr++
	case uniqID != "":
		query = listAlertsByUniqID
		args = append(args, uniqID, end, start)
		devIncr++
	default:
		query = listAlerts
		args = append(args, end, start)
	}

	if alarmID != "" {
		query += fmt.Sprintf(listAlertsAlarmID, 4+devIncr)
		args = append(args, alarmID)
		devIncr++
	}

	if userID != "" {
		query += fmt.Sprintf(listAlertsUserID, 4+devIncr)
		args = append(args, userID)
	}

	query += listAlertsOrder

	// Run list query.
	rows, err := d.ro.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("List rows.Close: %v", err)
		}
	}()

	var alerts []*api.Alert
	for rows.Next() {
		alert := &api.Alert{}
		var status string
		var createdAt time.Time

		if err = rows.Scan(&alert.OrgId, &alert.UniqId, &alert.AlarmId,
			&alert.UserId, &status, &alert.Error, &createdAt,
			&alert.TraceId); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		alert.Status = api.AlertStatus(api.AlertStatus_value[status])
		alert.CreatedAt = timestamppb.New(createdAt)
		alerts = append(alerts, alert)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return alerts, nil
}
