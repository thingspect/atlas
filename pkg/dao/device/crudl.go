package device

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createDevice = `
INSERT INTO devices (org_id, uniq_id, name, status, decoder, tags, created_at,
updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
RETURNING id, token
`

// Create creates a device in the database.
func (d *DAO) Create(ctx context.Context, dev *api.Device) (
	*api.Device, error,
) {
	dev.UniqId = strings.ToLower(dev.GetUniqId())

	now := time.Now().UTC().Truncate(time.Microsecond)
	dev.CreatedAt = timestamppb.New(now)
	dev.UpdatedAt = timestamppb.New(now)

	if err := d.rw.QueryRowContext(ctx, createDevice, dev.GetOrgId(),
		dev.GetUniqId(), dev.GetName(), dev.GetStatus().String(),
		dev.GetDecoder().String(), dev.GetTags(), now).Scan(&dev.Id,
		&dev.Token); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	return dev, nil
}

const readDevice = `
SELECT id, org_id, uniq_id, name, status, token, decoder, tags, created_at,
updated_at
FROM devices
WHERE (id, org_id) = ($1, $2)
`

// Read retrieves a device by ID and org ID.
func (d *DAO) Read(ctx context.Context, devID, orgID string) (
	*api.Device, error,
) {
	dev := &api.Device{}

	if d.cache != nil {
		ok, bDev, err := d.cache.GetB(ctx, devKey(orgID, devID))
		if err != nil {
			return nil, dao.DBToSentinel(err)
		}

		if ok {
			if err := proto.Unmarshal(bDev, dev); err != nil {
				return nil, dao.DBToSentinel(err)
			}

			return dev, nil
		}
	}

	var status, decoder string
	var createdAt, updatedAt time.Time

	if err := d.ro.QueryRowContext(ctx, readDevice, devID, orgID).Scan(&dev.Id,
		&dev.OrgId, &dev.UniqId, &dev.Name, &status, &dev.Token, &decoder,
		pgtype.NewMap().SQLScanner(&dev.Tags), &createdAt,
		&updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	dev.Status = api.Status(api.Status_value[status])
	dev.Decoder = api.Decoder(api.Decoder_value[decoder])
	dev.CreatedAt = timestamppb.New(createdAt)
	dev.UpdatedAt = timestamppb.New(updatedAt)

	// Cache write errors should not prevent successful database reads.
	if d.cache != nil {
		logger := alog.FromContext(ctx)

		bDev, err := proto.Marshal(dev)
		if err != nil {
			logger.Errorf("Read proto.Marshal: %v", err)

			return dev, nil
		}

		if err = d.cache.SetTTL(ctx, devKey(orgID, devID), bDev,
			d.exp); err != nil {
			logger.Errorf("Read d.cache.SetTTL: %v", err)
		}
	}

	return dev, nil
}

const readDeviceByUniqID = `
SELECT id, org_id, uniq_id, name, status, token, decoder, tags, created_at,
updated_at
FROM devices
WHERE uniq_id = $1
`

// ReadByUniqID retrieves a device by UniqID. This method does not limit by org
// ID and should only be used in the service layer.
func (d *DAO) ReadByUniqID(ctx context.Context, uniqID string) (
	*api.Device, error,
) {
	dev := &api.Device{}

	if d.cache != nil {
		ok, bDev, err := d.cache.GetB(ctx, devKeyByUniqID(uniqID))
		if err != nil {
			return nil, dao.DBToSentinel(err)
		}

		if ok {
			if err := proto.Unmarshal(bDev, dev); err != nil {
				return nil, dao.DBToSentinel(err)
			}

			return dev, nil
		}
	}

	var status, decoder string
	var createdAt, updatedAt time.Time

	if err := d.ro.QueryRowContext(ctx, readDeviceByUniqID, uniqID).Scan(
		&dev.Id, &dev.OrgId, &dev.UniqId, &dev.Name, &status, &dev.Token,
		&decoder, pgtype.NewMap().SQLScanner(&dev.Tags), &createdAt,
		&updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	dev.Status = api.Status(api.Status_value[status])
	dev.Decoder = api.Decoder(api.Decoder_value[decoder])
	dev.CreatedAt = timestamppb.New(createdAt)
	dev.UpdatedAt = timestamppb.New(updatedAt)

	// Cache write errors should not prevent successful database reads.
	if d.cache != nil {
		logger := alog.FromContext(ctx)

		bDev, err := proto.Marshal(dev)
		if err != nil {
			logger.Errorf("ReadByUniqID proto.Marshal: %v", err)

			return dev, nil
		}

		if err = d.cache.SetTTL(ctx, devKeyByUniqID(uniqID), bDev,
			d.exp); err != nil {
			logger.Errorf("ReadByUniqID d.cache.SetTTL: %v", err)
		}
	}

	return dev, nil
}

const updateDevice = `
UPDATE devices
SET uniq_id = $1, name = $2, status = $3, token = $4, decoder = $5, tags = $6,
updated_at = $7
WHERE (id, org_id) = ($8, $9)
RETURNING created_at
`

// Update updates a device in the database. CreatedAt should not update, so it
// is safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, dev *api.Device) (
	*api.Device, error,
) {
	dev.UniqId = strings.ToLower(dev.GetUniqId())

	var createdAt time.Time
	updatedAt := time.Now().UTC().Truncate(time.Microsecond)
	dev.UpdatedAt = timestamppb.New(updatedAt)

	if err := d.rw.QueryRowContext(ctx, updateDevice, dev.GetUniqId(),
		dev.GetName(), dev.GetStatus().String(), dev.GetToken(),
		dev.GetDecoder().String(), dev.GetTags(), updatedAt, dev.GetId(),
		dev.GetOrgId()).Scan(&createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	dev.CreatedAt = timestamppb.New(createdAt)

	// Invalidate cache on update.
	if d.cache != nil {
		if err := d.cache.Del(ctx, devKey(dev.GetOrgId(),
			dev.GetId())); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("Update devKey d.cache.Del: %v", err)
		}

		if err := d.cache.Del(ctx,
			devKeyByUniqID(dev.GetUniqId())); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("Update devKeyByUniqID d.cache.Del: %v", err)
		}
	}

	return dev, nil
}

const deleteDevice = `
DELETE FROM devices
WHERE (id, org_id) = ($1, $2)
`

// Delete deletes a device by ID and org ID.
func (d *DAO) Delete(ctx context.Context, devID, orgID string) error {
	// Verify a device exists before attempting to delete it. Do not remap the
	// error.
	dev, err := d.Read(ctx, devID, orgID)
	if err != nil {
		return err
	}

	_, err = d.rw.ExecContext(ctx, deleteDevice, devID, orgID)

	// Invalidate cache on delete.
	if d.cache != nil {
		if err := d.cache.Del(ctx, devKey(orgID, devID)); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("Delete devKey d.cache.Del: %v", err)
		}

		if err := d.cache.Del(ctx,
			devKeyByUniqID(dev.GetUniqId())); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("Delete devKeyByUniqID d.cache.Del: %v", err)
		}
	}

	return dao.DBToSentinel(err)
}

const countDevices = `
SELECT count(*)
FROM devices
WHERE org_id = $1
`

const countDevicesTag = `
AND $2 = ANY (tags)
`

const listDevices = `
SELECT id, org_id, uniq_id, name, status, token, decoder, tags, created_at,
updated_at
FROM devices
WHERE org_id = $1
`

const listDevicesTSAndID = `
AND (created_at > $%d
OR (created_at = $%d
AND id > $%d
))
`

const listDevicesTag = `
AND $%d = ANY (tags)
`

const listDevicesLimit = `
ORDER BY created_at ASC, id ASC
LIMIT %d
`

// List retrieves all devices by org ID with pagination and optional tag filter.
// If lBoundTS and prevID are zero values, the first page of results is
// returned. Limits of 0 or less do not apply a limit. List returns a slice of
// devices, a total count, and an error value.
func (d *DAO) List(
	ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
	limit int32, tag string,
) ([]*api.Device, int32, error) {
	// Build count query.
	cQuery := countDevices
	cArgs := []interface{}{orgID}

	if tag != "" {
		cQuery += countDevicesTag
		cArgs = append(cArgs, tag)
	}

	// Run count query.
	var count int32
	if err := d.ro.QueryRowContext(ctx, cQuery, cArgs...).Scan(
		&count); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	// Build list query.
	lQuery := listDevices
	lArgs := []interface{}{orgID}

	if prevID != "" && !lBoundTS.IsZero() {
		lQuery += fmt.Sprintf(listDevicesTSAndID, 2, 2, 3)
		lArgs = append(lArgs, lBoundTS, prevID)

		if tag != "" {
			lQuery += fmt.Sprintf(listDevicesTag, 4)
			lArgs = append(lArgs, tag)
		}
	} else if tag != "" {
		lQuery += fmt.Sprintf(listDevicesTag, 2)
		lArgs = append(lArgs, tag)
	}

	// Ordering is applied with the limit, which will always be present for API
	// usage, whereas lBoundTS and prevID will not for first pages.
	if limit > 0 {
		lQuery += fmt.Sprintf(listDevicesLimit, limit)
	}

	// Run list query.
	rows, err := d.ro.QueryContext(ctx, lQuery, lArgs...)
	if err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("List rows.Close: %v", err)
		}
	}()

	var devs []*api.Device
	pgtmap := pgtype.NewMap()
	for rows.Next() {
		dev := &api.Device{}
		var status, decoder string
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&dev.Id, &dev.OrgId, &dev.UniqId, &dev.Name, &status,
			&dev.Token, &decoder, pgtmap.SQLScanner(&dev.Tags), &createdAt,
			&updatedAt); err != nil {
			return nil, 0, dao.DBToSentinel(err)
		}

		dev.Status = api.Status(api.Status_value[status])
		dev.Decoder = api.Decoder(api.Decoder_value[decoder])
		dev.CreatedAt = timestamppb.New(createdAt)
		dev.UpdatedAt = timestamppb.New(updatedAt)
		devs = append(devs, dev)
	}

	if err = rows.Close(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, dao.DBToSentinel(err)
	}

	return devs, count, nil
}
