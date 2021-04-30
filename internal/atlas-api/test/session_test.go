// +build !unit

package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	iapi "github.com/thingspect/atlas/internal/atlas-api/api"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
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
			Password: globalPass,
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.NoError(t, err)
		require.Greater(t, len(loginResp.Token), 90)
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
			Password: random.String(10),
		})
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
			Password: random.String(10),
		})
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
			Password: globalPass,
		})
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
			Password: globalPass,
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+
			"unauthorized")
	})
}

func TestCreateKey(t *testing.T) {
	t.Parallel()

	t.Run("Create valid key", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = common.Role_BUILDER

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: key})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)
		require.NotEqual(t, key.Id, createKey.Key.Id)
		require.WithinDuration(t, time.Now(), createKey.Key.CreatedAt.AsTime(),
			2*time.Second)
		require.NotEmpty(t, createKey.Token)
	})

	t.Run("Create valid key with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(secondaryViewerGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: random.Key("api-key", uuid.NewString())})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.Nil(t, createKey)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Create sysadmin key as non-sysadmin", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = common.Role_SYS_ADMIN

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: key})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.Nil(t, createKey)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, role modification not allowed")
	})

	t.Run("Create invalid key", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Name = "api-key-" + random.String(80)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: key})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.Nil(t, createKey)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid CreateKeyRequest.Key: embedded message failed validation "+
			"| caused by: invalid Key.Name: value length must be between 5 "+
			"and 80 runes, inclusive")
	})
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()

	t.Run("Delete key by valid ID", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = common.Role_BUILDER

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: key})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)

		_, err = sessCli.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: createKey.Key.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		t.Run("Delete key by deleted ID", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				testTimeout)
			defer cancel()

			sessCli := api.NewSessionServiceClient(globalAdminKeyGRPCConn)
			_, err := sessCli.DeleteKey(ctx,
				&api.DeleteKeyRequest{Id: createKey.Key.Id})
			t.Logf("err: %v", err)
			require.EqualError(t, err, "rpc error: code = NotFound desc = "+
				"object not found")
		})
	})

	t.Run("Delete key with invalid key", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = common.Role_ADMIN

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: key})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)

		opts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithInsecure(),
			grpc.WithPerRPCCredentials(&credential{token: createKey.Token}),
		}
		keyConn, err := grpc.Dial(iapi.GRPCHost+iapi.GRPCPort, opts...)
		require.NoError(t, err)

		sessCli = api.NewSessionServiceClient(keyConn)
		_, err = sessCli.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: createKey.Key.Id})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		_, err = sessCli.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: createKey.Key.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = Unauthenticated desc = "+
			"unauthorized")
	})

	t.Run("Delete key with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(secondaryViewerGRPCConn)
		_, err := sessCli.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = PermissionDenied "+
			"desc = permission denied, ADMIN role required")
	})

	t.Run("Delete key by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		_, err := sessCli.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})

	t.Run("Deletes are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = common.Role_BUILDER

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: key})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)

		secCli := api.NewSessionServiceClient(secondaryAdminGRPCConn)
		_, err = secCli.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: createKey.Key.Id})
		t.Logf("err: %v", err)
		require.EqualError(t, err, "rpc error: code = NotFound desc = object "+
			"not found")
	})
}

func TestListKeys(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	keyIDs := []string{}
	keyNames := []string{}
	keyRoles := []common.Role{}
	for i := 0; i < 3; i++ {
		key := random.Key("api-key", uuid.NewString())
		key.Role = common.Role_BUILDER

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		createKey, err := sessCli.CreateKey(ctx,
			&api.CreateKeyRequest{Key: key})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.NoError(t, err)

		keyIDs = append(keyIDs, createKey.Key.Id)
		keyNames = append(keyNames, createKey.Key.Name)
		keyRoles = append(keyRoles, createKey.Key.Role)
	}

	t.Run("List keys by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		listKeys, err := sessCli.ListKeys(ctx, &api.ListKeysRequest{})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listKeys.Keys), 3)
		require.GreaterOrEqual(t, listKeys.TotalSize, int32(3))

		var found bool
		for _, key := range listKeys.Keys {
			if key.Id == keyIDs[len(keyIDs)-1] &&
				key.Name == keyNames[len(keyNames)-1] &&
				key.Role == keyRoles[len(keyRoles)-1] {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("List keys by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminKeyGRPCConn)
		listKeys, err := sessCli.ListKeys(ctx,
			&api.ListKeysRequest{PageSize: 2})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.NoError(t, err)
		require.Len(t, listKeys.Keys, 2)
		require.NotEmpty(t, listKeys.NextPageToken)
		require.GreaterOrEqual(t, listKeys.TotalSize, int32(3))

		nextKeys, err := sessCli.ListKeys(ctx, &api.ListKeysRequest{
			PageSize: 2, PageToken: listKeys.NextPageToken,
		})
		t.Logf("nextKeys, err: %+v, %v", nextKeys, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(nextKeys.Keys), 1)
		require.GreaterOrEqual(t, nextKeys.TotalSize, int32(3))
	})

	t.Run("List keys with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewSessionServiceClient(secondaryViewerGRPCConn)
		listKeys, err := secCli.ListKeys(ctx, &api.ListKeysRequest{})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.Nil(t, listKeys)
		require.EqualError(t, err, "rpc error: code = PermissionDenied desc = "+
			"permission denied, ADMIN role required")
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewSessionServiceClient(secondaryAdminGRPCConn)
		listKeys, err := secCli.ListKeys(ctx, &api.ListKeysRequest{})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.NoError(t, err)
		require.Len(t, listKeys.Keys, 1)
		require.Equal(t, int32(1), listKeys.TotalSize)
	})

	t.Run("List keys by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		sessCli := api.NewSessionServiceClient(globalAdminGRPCConn)
		listKeys, err := sessCli.ListKeys(ctx,
			&api.ListKeysRequest{PageToken: badUUID})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.Nil(t, listKeys)
		require.EqualError(t, err, "rpc error: code = InvalidArgument desc = "+
			"invalid page token")
	})
}
