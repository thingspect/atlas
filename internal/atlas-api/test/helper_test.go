// +build !unit

package test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	iapi "github.com/thingspect/atlas/internal/atlas-api/api"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
)

// credential provides token-based credentials for gRPC.
type credential struct {
	token string
}

// GetRequestMetadata returns authentication metadata and implements the
// PerRPCCredentials interface.
func (c *credential) GetRequestMetadata(ctx context.Context,
	uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + c.token,
	}, nil
}

// RequireTransportSecurity implements the PerRPCCredentials interface.
// Transport security is not required due to use behind TLS termination.
func (c *credential) RequireTransportSecurity() bool {
	return false
}

func authGRPCConn(role common.Role) (string, *grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 14*time.Second)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-helper"))
	if err != nil {
		return "", nil, err
	}

	user := random.User("api-helper", createOrg.Id)
	user.Role = role
	user.Status = common.Status_ACTIVE
	createUser, err := globalUserDAO.Create(ctx, user)
	if err != nil {
		return "", nil, err
	}

	if err = globalUserDAO.UpdatePassword(ctx, createUser.Id, createOrg.Id,
		globalHash); err != nil {
		return "", nil, err
	}

	sessCli := api.NewSessionServiceClient(globalNoAuthGRPCConn)
	login, err := sessCli.Login(ctx, &api.LoginRequest{
		Email: createUser.Email, OrgName: createOrg.Name, Password: globalPass,
	})
	if err != nil {
		return "", nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&credential{token: login.Token}),
	}
	authConn, err := grpc.Dial(iapi.GRPCHost+iapi.GRPCPort, opts...)
	if err != nil {
		return "", nil, err
	}

	return createOrg.Id, authConn, nil
}

func keyGRPCConn(conn *grpc.ClientConn, role common.Role) (*grpc.ClientConn,
	error) {
	key := random.Key("api-key", uuid.NewString())
	key.Role = role

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sessCli := api.NewSessionServiceClient(conn)
	createKey, err := sessCli.CreateKey(ctx, &api.CreateKeyRequest{Key: key})
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&credential{token: createKey.Token}),
	}
	keyConn, err := grpc.Dial(iapi.GRPCHost+iapi.GRPCPort, opts...)
	if err != nil {
		return nil, err
	}

	return keyConn, nil
}
