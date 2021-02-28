package lora

import (
	"context"

	as "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
	"github.com/brocaar/chirpstack-api/go/v3/common"
)

// CreateGateway creates a gateway by UniqID.
func (cs *Chirpstack) CreateGateway(ctx context.Context, uniqID string) error {
	gwCli := as.NewGatewayServiceClient(cs.conn)
	_, err := gwCli.Create(ctx, &as.CreateGatewayRequest{Gateway: &as.Gateway{
		Id:              uniqID,
		Name:            uniqID,
		Description:     uniqID,
		Location:        &common.Location{},
		OrganizationId:  cs.orgID,
		NetworkServerId: cs.nsID,
	}})

	return err
}

// DeleteGateway deletes a gateway by UniqID.
func (cs *Chirpstack) DeleteGateway(ctx context.Context, uniqID string) error {
	gwCli := as.NewGatewayServiceClient(cs.conn)
	_, err := gwCli.Delete(ctx, &as.DeleteGatewayRequest{Id: uniqID})

	return err
}
