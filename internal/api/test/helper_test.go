// +build !unit

package test

import (
	"context"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
)

type credential struct {
	token string
}

func (c *credential) GetRequestMetadata(ctx context.Context,
	uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + c.token,
	}, nil
}

func (c *credential) RequireTransportSecurity() bool {
	return false
}

func authGRPCConn(grpcAddr string) (string, *grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	org := org.Org{Name: "api-helper-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	if err != nil {
		return "", nil, err
	}

	user := &api.User{OrgId: createOrg.ID, Email: "api-helper-" +
		random.Email()}
	createUser, err := globalUserDAO.Create(ctx, user, globalHash)
	if err != nil {
		return "", nil, err
	}

	sessCli := api.NewSessionServiceClient(globalNoAuthGRPCConn)
	loginResp, err := sessCli.Login(ctx, &api.LoginRequest{
		Email: createUser.Email, OrgName: createOrg.Name,
		Password: globalPass})
	if err != nil {
		return "", nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&credential{token: loginResp.Token}),
	}
	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return "", nil, err
	}
	return createOrg.ID, conn, nil
}
