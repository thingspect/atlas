// Package lora provides functions to create and modify LoRaWAN resources.
package lora

//go:generate mockgen -source loraer.go -destination mock_loraer.go -package lora

import "context"

// Loraer defines the methods provided by a Lora.
type Loraer interface {
	// CreateGateway creates a gateway by UniqID.
	CreateGateway(ctx context.Context, uniqID string) error
	// DeleteGateway deletes a gateway by UniqID.
	DeleteGateway(ctx context.Context, uniqID string) error
	// CreateDevice creates a device by UniqID and application key.
	CreateDevice(ctx context.Context, uniqID, appKey string) error
	// DeleteDevice deletes a device by UniqID.
	DeleteDevice(ctx context.Context, uniqID string) error
}
