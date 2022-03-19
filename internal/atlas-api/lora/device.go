package lora

import (
	"context"

	as "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
)

// CreateDevice creates a device by UniqID and application key.
func (cs *Chirpstack) CreateDevice(
	ctx context.Context, uniqID, appKey string,
) error {
	devCli := as.NewDeviceServiceClient(cs.conn)
	if _, err := devCli.Create(ctx, &as.CreateDeviceRequest{Device: &as.Device{
		DevEui:          uniqID,
		Name:            uniqID,
		ApplicationId:   cs.appID,
		Description:     uniqID,
		DeviceProfileId: cs.devProfID,
	}}); err != nil {
		return err
	}

	if _, err := devCli.CreateKeys(ctx, &as.CreateDeviceKeysRequest{
		DeviceKeys: &as.DeviceKeys{
			DevEui: uniqID,
			NwkKey: appKey,
		},
	}); err != nil {
		// Perform a best-effort rollback in the event of application key
		// failure, but return original error.
		_ = cs.DeleteDevice(ctx, uniqID)

		return err
	}

	return nil
}

// DeleteDevice deletes a device by UniqID.
func (cs *Chirpstack) DeleteDevice(ctx context.Context, uniqID string) error {
	devCli := as.NewDeviceServiceClient(cs.conn)
	_, err := devCli.Delete(ctx, &as.DeleteDeviceRequest{DevEui: uniqID})

	return err
}
