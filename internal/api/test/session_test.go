// +build !unit

package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	createOrg, err := globalOrgDAO.Create(ctx, random.Org("api-session"))
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	user := random.User("api-session", createOrg.Id)
	user.Role = common.Role_ADMIN
	user.Status = common.Status_ACTIVE
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	err = globalUserDAO.UpdatePassword(ctx, createUser.Id, createOrg.Id,
		globalHash)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	disUser := random.User("api-session", createOrg.Id)
	disUser.Status = common.Status_DISABLED
	createDisUser, err := globalUserDAO.Create(ctx, disUser)
	t.Logf("createDisUser, err: %+v, %v", createDisUser, err)
	require.NoError(t, err)

	err = globalUserDAO.UpdatePassword(ctx, createDisUser.Id, createOrg.Id,
		globalHash)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	contUser := random.User("api-session", createOrg.Id)
	contUser.Role = common.Role_CONTACT
	contUser.Status = common.Status_ACTIVE
	createContUser, err := globalUserDAO.Create(ctx, contUser)
	t.Logf("createContUser, err: %+v, %v", createContUser, err)
	require.NoError(t, err)

	err = globalUserDAO.UpdatePassword(ctx, createContUser.Id, createOrg.Id,
		globalHash)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	t.Run("Log in valid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalNoAuthGRPCConn)
		loginResp, err := sessCli.Login(ctx, &api.LoginRequest{
			Email: createUser.Email, OrgName: createOrg.Name,
			Password: globalPass})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(loginResp.Token), 90)
		require.WithinDuration(t, time.Now().Add(
			session.WebTokenExp*time.Second), loginResp.ExpiresAt.AsTime(),
			2*time.Second)
	})

	t.Run("Log in unknown user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalNoAuthGRPCConn)
		loginResp, err := sessCli.Login(ctx, &api.LoginRequest{
			Email: random.Email(), OrgName: random.String(10),
			Password: random.String(10)})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+
			"unauthorized")
	})

	t.Run("Log in wrong password", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalNoAuthGRPCConn)
		loginResp, err := sessCli.Login(ctx, &api.LoginRequest{
			Email: createUser.Email, OrgName: createOrg.Name,
			Password: random.String(10)})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+
			"unauthorized")
	})

	t.Run("Log in disabled user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalNoAuthGRPCConn)
		loginResp, err := sessCli.Login(ctx, &api.LoginRequest{
			Email: createDisUser.Email, OrgName: createOrg.Name,
			Password: globalPass})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+
			"unauthorized")
	})

	t.Run("Log in contact user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalNoAuthGRPCConn)
		loginResp, err := sessCli.Login(ctx, &api.LoginRequest{
			Email: createContUser.Email, OrgName: createOrg.Name,
			Password: globalPass})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+
			"unauthorized")
	})
}
