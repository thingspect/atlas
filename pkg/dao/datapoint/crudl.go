package datapoint

import (
	"context"
	"strings"
	"time"

	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createDataPoint = `
INSERT INTO data_points (org_id, uniq_id, attr, int_val, fl64_val, str_val,
bool_val, bytes_val, created_at, trace_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

// Create creates a data point in the database. Data points are retrieved
// elsewhere in bulk, so only an error value is returned.
func (d *DAO) Create(ctx context.Context, point *common.DataPoint,
	orgID string) error {
	// Truncate timestamp to milliseconds for deduplication.
	createdAt := point.Ts.AsTime().UTC().Truncate(time.Millisecond)

	var intVal *int32
	var fl64Val *float64
	var strVal *string
	var boolVal *bool
	var bytesVal []byte

	switch v := point.ValOneof.(type) {
	case *common.DataPoint_IntVal:
		intVal = &v.IntVal
	case *common.DataPoint_Fl64Val:
		fl64Val = &v.Fl64Val
	case *common.DataPoint_StrVal:
		strVal = &v.StrVal
	case *common.DataPoint_BoolVal:
		boolVal = &v.BoolVal
	case *common.DataPoint_BytesVal:
		bytesVal = v.BytesVal
	}

	_, err := d.pg.ExecContext(ctx, createDataPoint, orgID,
		strings.ToLower(point.UniqId), point.Attr, intVal, fl64Val, strVal,
		boolVal, bytesVal, createdAt, point.TraceId)
	return dao.DBToSentinel(err)
}

const listDataPointsByUniqID = `
SELECT d.uniq_id, d.attr, d.int_val, d.fl64_val, d.str_val, d.bool_val,
d.bytes_val, d.created_at, d.trace_id
FROM data_points d
WHERE (d.org_id, d.uniq_id) = ($1, $2)
AND d.created_at <= $3
AND d.created_at > $4
`

const listDataPointsByDevID = `
SELECT d.uniq_id, d.attr, d.int_val, d.fl64_val, d.str_val, d.bool_val,
d.bytes_val, d.created_at, d.trace_id
FROM data_points d
INNER JOIN devices de ON (
  d.org_id = de.org_id
  AND d.uniq_id = de.uniq_id
)
WHERE (d.org_id, de.id) = ($1, $2)
AND d.created_at <= $3
AND d.created_at > $4
`

const listDataPointsAttr = `
AND d.attr = $5
`

const listDataPointsOrder = `
ORDER BY d.created_at DESC
`

// List retrieves all data points by org ID, UniqID or device ID, optional
// attribute, and [end, start) times. If both uniqID and devID are provided,
// uniqID takes precedence and devID is ignored.
func (d *DAO) List(ctx context.Context, orgID, uniqID, devID, attr string, end,
	start time.Time) ([]*common.DataPoint, error) {
	// Build list query.
	query := listDataPointsByUniqID
	args := []interface{}{orgID}

	if uniqID == "" && devID != "" {
		query = listDataPointsByDevID
		args = append(args, devID, end, start)
	} else {
		args = append(args, uniqID, end, start)
	}

	if attr != "" {
		query += listDataPointsAttr
		args = append(args, attr)
	}

	query += listDataPointsOrder

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

	var points []*common.DataPoint
	for rows.Next() {
		point := &common.DataPoint{}
		var intVal *int32
		var fl64Val *float64
		var strVal *string
		var boolVal *bool
		var bytesVal []byte
		var createdAt time.Time

		if err = rows.Scan(&point.UniqId, &point.Attr, &intVal, &fl64Val,
			&strVal, &boolVal, &bytesVal, &createdAt,
			&point.TraceId); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		switch {
		case intVal != nil:
			point.ValOneof = &common.DataPoint_IntVal{IntVal: *intVal}
		case fl64Val != nil:
			point.ValOneof = &common.DataPoint_Fl64Val{Fl64Val: *fl64Val}
		case strVal != nil:
			point.ValOneof = &common.DataPoint_StrVal{StrVal: *strVal}
		case boolVal != nil:
			point.ValOneof = &common.DataPoint_BoolVal{BoolVal: *boolVal}
		case bytesVal != nil:
			point.ValOneof = &common.DataPoint_BytesVal{BytesVal: bytesVal}
		}

		point.Ts = timestamppb.New(createdAt)
		points = append(points, point)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return points, nil
}

const latestDataPointsByUniqID = `
SELECT
  d.uniq_id,
  d.attr,
  d.int_val,
  d.fl64_val,
  d.str_val,
  d.bool_val,
  d.bytes_val,
  d.created_at,
  d.trace_id
FROM
  data_points d
  INNER JOIN (
    SELECT
      org_id,
      uniq_id,
      attr,
      MAX(created_at) AS created_at
    FROM
      data_points
    WHERE
      (org_id, uniq_id) = ($1, $2)
    GROUP BY
      org_id,
      uniq_id,
      attr
  ) m ON (d.org_id, d.uniq_id, d.attr, d.created_at) = (
    m.org_id,
    m.uniq_id,
    m.attr,
    m.created_at
  )
ORDER BY
  d.attr ASC
`

const latestDataPointsByDevID = `
SELECT
  d.uniq_id,
  d.attr,
  d.int_val,
  d.fl64_val,
  d.str_val,
  d.bool_val,
  d.bytes_val,
  d.created_at,
  d.trace_id
FROM
  data_points d
  INNER JOIN (
    SELECT
      id.org_id,
      id.uniq_id,
      id.attr,
      MAX(id.created_at) AS created_at
    FROM
      data_points id
      INNER JOIN devices de ON (
        id.org_id = de.org_id
        AND id.uniq_id = de.uniq_id
      )
    WHERE
      (id.org_id, id.uniq_id, de.id) = ($1, de.uniq_id, $2)
    GROUP BY
      id.org_id,
      id.uniq_id,
      id.attr
  ) m ON (d.org_id, d.uniq_id, d.attr, d.created_at) = (
    m.org_id,
    m.uniq_id,
    m.attr,
    m.created_at
  )
ORDER BY
  d.attr ASC
`

// Latest retrieves the latest data point for each of a device's attributes by
// org ID and UniqID or device ID. If both uniqID and devID are provided, uniqID
// takes precedence and devID is ignored.
func (d *DAO) Latest(ctx context.Context, orgID, uniqID,
	devID string) ([]*common.DataPoint, error) {
	// Build latest query.
	query := latestDataPointsByUniqID
	args := []interface{}{orgID}

	if uniqID == "" && devID != "" {
		query = latestDataPointsByDevID
		args = append(args, devID)
	} else {
		args = append(args, uniqID)
	}

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

	var points []*common.DataPoint
	for rows.Next() {
		point := &common.DataPoint{}
		var intVal *int32
		var fl64Val *float64
		var strVal *string
		var boolVal *bool
		var bytesVal []byte
		var createdAt time.Time

		if err = rows.Scan(&point.UniqId, &point.Attr, &intVal, &fl64Val,
			&strVal, &boolVal, &bytesVal, &createdAt,
			&point.TraceId); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		switch {
		case intVal != nil:
			point.ValOneof = &common.DataPoint_IntVal{IntVal: *intVal}
		case fl64Val != nil:
			point.ValOneof = &common.DataPoint_Fl64Val{Fl64Val: *fl64Val}
		case strVal != nil:
			point.ValOneof = &common.DataPoint_StrVal{StrVal: *strVal}
		case boolVal != nil:
			point.ValOneof = &common.DataPoint_BoolVal{BoolVal: *boolVal}
		case bytesVal != nil:
			point.ValOneof = &common.DataPoint_BytesVal{BytesVal: bytesVal}
		}

		point.Ts = timestamppb.New(createdAt)
		points = append(points, point)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return points, nil
}
