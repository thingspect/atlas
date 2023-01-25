package lora

import (
	"context"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/common"
)

// CreateGateway creates a gateway by UniqID.
func (cs *Chirpstack) CreateGateway(ctx context.Context, uniqID string) error {
	gwCli := api.NewGatewayServiceClient(cs.conn)
	_, err := gwCli.Create(ctx, &api.CreateGatewayRequest{Gateway: &api.Gateway{
		GatewayId:     uniqID,
		Name:          uniqID,
		Description:   uniqID,
		Location:      &common.Location{},
		TenantId:      cs.tenantID,
		StatsInterval: 900,
	}})

	return err
}

// DeleteGateway deletes a gateway by UniqID.
func (cs *Chirpstack) DeleteGateway(ctx context.Context, uniqID string) error {
	gwCli := api.NewGatewayServiceClient(cs.conn)
	_, err := gwCli.Delete(ctx, &api.DeleteGatewayRequest{GatewayId: uniqID})

	return err
}
