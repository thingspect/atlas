// +build !unit

package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	org := org.Org{Name: "api-session-" + random.String(10)}
	createOrg, err := globalOrgDAO.Create(ctx, org)
	t.Logf("createOrg, err: %+v, %v", createOrg, err)
	require.NoError(t, err)

	user := &api.User{OrgId: createOrg.ID, Email: "api-session-" +
		random.Email(), Status: api.Status_ACTIVE}
	createUser, err := globalUserDAO.Create(ctx, user)
	t.Logf("createUser, err: %+v, %v", createUser, err)
	require.NoError(t, err)

	err = globalUserDAO.UpdatePassword(ctx, createUser.Id, createOrg.ID,
		globalHash)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	disUser := &api.User{OrgId: createOrg.ID, Email: "api-session-" +
		random.Email(), Status: api.Status_DISABLED}
	createDisUser, err := globalUserDAO.Create(ctx, disUser)
	t.Logf("createDisUser, err: %+v, %v", createDisUser, err)
	require.NoError(t, err)

	t.Run("Log in valid user", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
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
}
