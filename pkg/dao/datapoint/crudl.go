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

// Create creates a data point in the database. Data points are retrieved in
// bulk, so only an error value is returned.
func (d *DAO) Create(ctx context.Context, point *common.DataPoint,
	orgID string) error {
	// Truncate timestamp to milliseconds for deduplication.
	createdAt := point.Ts.AsTime().UTC().Truncate(time.Millisecond)

	var intVal *int64
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

const listDataPoints = `
SELECT uniq_id, attr, int_val, fl64_val, str_val, bool_val, bytes_val,
created_at, trace_id
FROM data_points
WHERE (org_id, uniq_id) = ($1, $2)
AND created_at >= $3
AND created_at < $4
`

// List retrieves all data points by org ID, UniqID, and [start, end) times.
func (d *DAO) List(ctx context.Context, orgID, uniqID string, start,
	end time.Time) ([]*common.DataPoint, error) {
	var points []*common.DataPoint

	rows, err := d.pg.QueryContext(ctx, listDataPoints, orgID, uniqID,
		start.Truncate(time.Millisecond), end.Truncate(time.Millisecond))
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			alog.Errorf("List rows.Close: %v", err)
		}
	}()

	for rows.Next() {
		point := &common.DataPoint{}
		var intVal *int64
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
