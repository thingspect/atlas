package device

import (
	"context"
	"strings"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const createDevice = `
INSERT INTO devices (org_id, uniq_id, is_disabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, token
`

// Create creates a device in the database.
func (d *DAO) Create(ctx context.Context, dev *api.Device) (*api.Device,
	error) {
	dev.UniqId = strings.ToLower(dev.UniqId)
	now := time.Now().UTC().Truncate(time.Microsecond)
	dev.CreatedAt = timestamppb.New(now)
	dev.UpdatedAt = timestamppb.New(now)

	if err := d.pg.QueryRowContext(ctx, createDevice, dev.OrgId, dev.UniqId,
		dev.IsDisabled, now, now).Scan(&dev.Id, &dev.Token); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return dev, nil
}

const readDevice = `
SELECT id, org_id, uniq_id, is_disabled, token, created_at, updated_at
FROM devices
WHERE id = $1
AND org_id = $2
`

// Read retrieves a device by ID and org ID.
func (d *DAO) Read(ctx context.Context, devID, orgID string) (*api.Device,
	error) {
	dev := &api.Device{}
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readDevice, devID, orgID).Scan(&dev.Id,
		&dev.OrgId, &dev.UniqId, &dev.IsDisabled, &dev.Token, &createdAt,
		&updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	dev.CreatedAt = timestamppb.New(createdAt)
	dev.UpdatedAt = timestamppb.New(updatedAt)
	return dev, nil
}

const readDeviceByUniqID = `
SELECT id, org_id, uniq_id, is_disabled, token, created_at, updated_at
FROM devices
WHERE uniq_id = $1
`

// ReadByUniqID retrieves a device by UniqID. This method does not limit by org
// ID and should only be used in the service layer.
func (d *DAO) ReadByUniqID(ctx context.Context, uniqID string) (*api.Device,
	error) {
	dev := &api.Device{}
	var createdAt, updatedAt time.Time

	if err := d.pg.QueryRowContext(ctx, readDeviceByUniqID, uniqID).Scan(
		&dev.Id, &dev.OrgId, &dev.UniqId, &dev.IsDisabled, &dev.Token,
		&createdAt, &updatedAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	dev.CreatedAt = timestamppb.New(createdAt)
	dev.UpdatedAt = timestamppb.New(updatedAt)
	return dev, nil
}

const updateDevice = `
UPDATE devices
SET uniq_id = $1, is_disabled = $2, token = $3, updated_at = $4
WHERE id = $5
AND org_id = $6
RETURNING created_at
`

// Update updates a device in the database. CreatedAt should not update, so it
// is safe to override it at the DAO level.
func (d *DAO) Update(ctx context.Context, dev *api.Device) (*api.Device,
	error) {
	var createdAt time.Time
	updatedAt := time.Now().UTC().Truncate(time.Microsecond)
	dev.UpdatedAt = timestamppb.New(updatedAt)

	if err := d.pg.QueryRowContext(ctx, updateDevice, dev.UniqId,
		dev.IsDisabled, dev.Token, updatedAt, dev.Id, dev.OrgId).Scan(
		&createdAt); err != nil {
		return nil, dao.DBToSentinel(err)
	}

	dev.CreatedAt = timestamppb.New(createdAt)
	return dev, nil
}

const deleteDevice = `
DELETE FROM devices
WHERE id = $1
AND org_id = $2
`

// Delete deletes a device by ID and org ID.
func (d *DAO) Delete(ctx context.Context, devID, orgID string) error {
	// Verify a device exists before attempting to delete it. Do not remap the
	// error.
	if _, err := d.Read(ctx, devID, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteDevice, devID, orgID)
	return dao.DBToSentinel(err)
}

const listDevices = `
SELECT id, org_id, uniq_id, is_disabled, token, created_at, updated_at
FROM devices
WHERE org_id = $1
`

// List retrieves all devices by org ID.
func (d *DAO) List(ctx context.Context, orgID string) ([]*api.Device, error) {
	var devs []*api.Device

	rows, err := d.pg.QueryContext(ctx, listDevices, orgID)
	if err != nil {
		return nil, dao.DBToSentinel(err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			alog.Errorf("List rows.Close: %v", err)
		}
	}()

	for rows.Next() {
		dev := &api.Device{}
		var createdAt, updatedAt time.Time

		if err = rows.Scan(&dev.Id, &dev.OrgId, &dev.UniqId, &dev.IsDisabled,
			&dev.Token, &createdAt, &updatedAt); err != nil {
			return nil, dao.DBToSentinel(err)
		}

		dev.CreatedAt = timestamppb.New(createdAt)
		dev.UpdatedAt = timestamppb.New(updatedAt)
		devs = append(devs, dev)
	}

	if err = rows.Close(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	if err = rows.Err(); err != nil {
		return nil, dao.DBToSentinel(err)
	}
	return devs, nil
}
