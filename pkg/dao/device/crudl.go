package device

import (
	"context"
	"strings"
	"time"

	"github.com/thingspect/atlas/pkg/alog"
)

const createDevice = `
INSERT INTO devices (org_id, uniq_id, is_disabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, token
`

// Create creates a device in the database.
func (d *DAO) Create(ctx context.Context, dev Device) (*Device, error) {
	dev.UniqID = strings.ToLower(dev.UniqID)
	dev.CreatedAt = time.Now().UTC().Truncate(time.Microsecond)
	dev.UpdatedAt = time.Now().UTC().Truncate(time.Microsecond)

	if err := d.pg.QueryRowContext(ctx, createDevice, dev.OrgID, dev.UniqID,
		dev.Disabled, dev.CreatedAt, dev.UpdatedAt).Scan(&dev.ID,
		&dev.Token); err != nil {
		return nil, err
	}
	return &dev, nil
}

const readDevice = `
SELECT id, org_id, uniq_id, is_disabled, token, created_at, updated_at
FROM devices
WHERE id = $1
AND org_id = $2
`

// Read retrieves a device by ID and org ID.
func (d *DAO) Read(ctx context.Context, devID, orgID string) (*Device, error) {
	var dev Device

	if err := d.pg.QueryRowContext(ctx, readDevice, devID, orgID).Scan(&dev.ID,
		&dev.OrgID, &dev.UniqID, &dev.Disabled, &dev.Token, &dev.CreatedAt,
		&dev.UpdatedAt); err != nil {
		return nil, err
	}

	dev.CreatedAt = dev.CreatedAt.In(time.UTC)
	dev.UpdatedAt = dev.UpdatedAt.In(time.UTC)
	return &dev, nil
}

const readDeviceByUniqID = `
SELECT id, org_id, uniq_id, is_disabled, token, created_at, updated_at
FROM devices
WHERE uniq_id = $1
`

// ReadByUniqID retrieves a device by UniqID. This method does not limit by org
// ID and should only be used in the service layer.
func (d *DAO) ReadByUniqID(ctx context.Context, uniqID string) (*Device,
	error) {
	var dev Device

	if err := d.pg.QueryRowContext(ctx, readDeviceByUniqID, uniqID).Scan(
		&dev.ID, &dev.OrgID, &dev.UniqID, &dev.Disabled, &dev.Token,
		&dev.CreatedAt, &dev.UpdatedAt); err != nil {
		return nil, err
	}

	dev.CreatedAt = dev.CreatedAt.In(time.UTC)
	dev.UpdatedAt = dev.UpdatedAt.In(time.UTC)
	return &dev, nil
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
func (d *DAO) Update(ctx context.Context, dev Device) (*Device, error) {
	dev.UpdatedAt = time.Now().UTC().Truncate(time.Microsecond)

	if err := d.pg.QueryRowContext(ctx, updateDevice, dev.UniqID, dev.Disabled,
		dev.Token, dev.UpdatedAt, dev.ID, dev.OrgID).Scan(
		&dev.CreatedAt); err != nil {
		return nil, err
	}

	dev.CreatedAt = dev.CreatedAt.In(time.UTC)
	return &dev, nil
}

const deleteDevice = `
DELETE FROM devices
WHERE id = $1
AND org_id = $2
`

// Delete deletes a device by ID and org ID.
func (d *DAO) Delete(ctx context.Context, devID, orgID string) error {
	// Verify a org exists before attempting to delete it.
	if _, err := d.Read(ctx, devID, orgID); err != nil {
		return err
	}

	_, err := d.pg.ExecContext(ctx, deleteDevice, devID, orgID)
	return err
}

const listDevices = `
SELECT id, org_id, uniq_id, is_disabled, token, created_at, updated_at
FROM devices
WHERE org_id = $1
`

// List retrieves all devices by org ID.
func (d *DAO) List(ctx context.Context, orgID string) ([]*Device, error) {
	var devs []*Device

	rows, err := d.pg.QueryContext(ctx, listDevices, orgID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			alog.Errorf("List rows.Close: %v", err)
		}
	}()

	for rows.Next() {
		var dev Device
		if err = rows.Scan(&dev.ID, &dev.OrgID, &dev.UniqID, &dev.Disabled,
			&dev.Token, &dev.CreatedAt, &dev.UpdatedAt); err != nil {
			return nil, err
		}

		dev.CreatedAt = dev.CreatedAt.In(time.UTC)
		dev.UpdatedAt = dev.UpdatedAt.In(time.UTC)
		devs = append(devs, &dev)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return devs, nil
}
