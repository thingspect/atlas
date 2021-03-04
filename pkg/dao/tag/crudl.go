package tag

import (
	"context"

	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
)

const listTags = `
SELECT unnest(tags) AS tag
FROM devices
WHERE org_id = $1
UNION
SELECT unnest(tags) AS tag
FROM users
WHERE org_id = $1
ORDER BY tag
`

// List retrieves all device and user tags by org ID.
func (d *DAO) List(ctx context.Context, orgID string) ([]string, error) {
	rows, err := d.pg.QueryContext(ctx, listTags, orgID)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("List rows.Close: %v", err)
		}
	}()

	var tags []string
	for rows.Next() {
		var tag string

		if err = rows.Scan(&tag); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		tags = append(tags, tag)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return tags, nil
}
