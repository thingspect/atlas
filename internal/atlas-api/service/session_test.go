//go:build !integration

package service

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("Log in valid user", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-session")
		user := random.User("api-session", org.GetId())
		user.Role = api.Role_ADMIN
		user.Status = api.Status_ACTIVE

		userer := NewMockUserer(gomock.NewController(t))
		userer.EXPECT().ReadByEmail(gomock.Any(), user.GetEmail(), org.GetName()).
			Return(user, globalHash, nil).Times(1)

		pwtKey := make([]byte, 32)
		_, err := rand.Read(pwtKey)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		sessSvc := NewSession(userer, nil, nil, pwtKey)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.GetEmail(), OrgName: org.GetName(), Password: globalPass,
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.NoError(t, err)
		require.Greater(t, len(loginResp.GetToken()), 90)
		require.WithinDuration(t, time.Now().Add(
			session.WebTokenExp*time.Second), loginResp.GetExpiresAt().AsTime(),
			2*time.Second)
	})

	t.Run("Log in unknown user", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-session")
		user := random.User("api-session", org.GetId())
		user.Status = api.Status_ACTIVE

		userer := NewMockUserer(gomock.NewController(t))
		userer.EXPECT().ReadByEmail(gomock.Any(), user.GetEmail(), org.GetName()).
			Return(nil, nil, dao.ErrNotFound).Times(1)

		pwtKey := make([]byte, 32)
		_, err := rand.Read(pwtKey)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		sessSvc := NewSession(userer, nil, nil, pwtKey)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.GetEmail(), OrgName: org.GetName(), Password: globalPass,
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})

	t.Run("Log in wrong password", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-session")
		user := random.User("api-session", org.GetId())
		user.Status = api.Status_ACTIVE

		userer := NewMockUserer(gomock.NewController(t))
		userer.EXPECT().ReadByEmail(gomock.Any(), user.GetEmail(), org.GetName()).
			Return(user, globalHash, nil).Times(1)

		pwtKey := make([]byte, 32)
		_, err := rand.Read(pwtKey)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		sessSvc := NewSession(userer, nil, nil, pwtKey)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.GetEmail(), OrgName: org.GetName(), Password: random.String(10),
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})

	t.Run("Log in disabled user", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-session")
		user := random.User("api-session", org.GetId())
		user.Status = api.Status_DISABLED

		userer := NewMockUserer(gomock.NewController(t))
		userer.EXPECT().ReadByEmail(gomock.Any(), user.GetEmail(), org.GetName()).
			Return(user, globalHash, nil).Times(1)

		pwtKey := make([]byte, 32)
		_, err := rand.Read(pwtKey)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		sessSvc := NewSession(userer, nil, nil, pwtKey)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.GetEmail(), OrgName: org.GetName(), Password: globalPass,
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})

	t.Run("Log in contact user", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-session")
		user := random.User("api-session", org.GetId())
		user.Role = api.Role_CONTACT
		user.Status = api.Status_ACTIVE

		userer := NewMockUserer(gomock.NewController(t))
		userer.EXPECT().ReadByEmail(gomock.Any(), user.GetEmail(), org.GetName()).
			Return(user, globalHash, nil).Times(1)

		pwtKey := make([]byte, 32)
		_, err := rand.Read(pwtKey)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		sessSvc := NewSession(userer, nil, nil, pwtKey)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.GetEmail(), OrgName: org.GetName(), Password: globalPass,
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})

	t.Run("Log in wrong key", func(t *testing.T) {
		t.Parallel()

		org := random.Org("api-session")
		user := random.User("api-session", org.GetId())
		user.Role = api.Role_ADMIN
		user.Status = api.Status_ACTIVE

		userer := NewMockUserer(gomock.NewController(t))
		userer.EXPECT().ReadByEmail(gomock.Any(), user.GetEmail(), org.GetName()).
			Return(user, globalHash, nil).Times(1)

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		sessSvc := NewSession(userer, nil, nil, nil)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.GetEmail(), OrgName: org.GetName(), Password: globalPass,
		})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})
}

func TestCreateKey(t *testing.T) {
	t.Parallel()

	t.Run("Create valid key", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = api.Role_ADMIN
		retKey, _ := proto.Clone(key).(*api.Key)

		keyer := NewMockKeyer(gomock.NewController(t))
		keyer.EXPECT().Create(gomock.Any(), key).Return(retKey, nil).Times(1)

		pwtKey := make([]byte, 32)
		_, err := rand.Read(pwtKey)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: key.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, nil, pwtKey)
		createKey, err := keySvc.CreateKey(ctx, &api.CreateKeyRequest{Key: key})
		t.Logf("key, createKey, err: %+v, %+v, %v", key, createKey, err)
		require.NoError(t, err)

		// Normalize token.
		resp := &api.CreateKeyResponse{Key: key, Token: createKey.GetToken()}

		require.EqualExportedValues(t, resp, createKey)
	})

	t.Run("Create key with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, nil, nil)
		createKey, err := keySvc.CreateKey(ctx, &api.CreateKeyRequest{})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.Nil(t, createKey)
		require.Equal(t, errPerm(api.Role_ADMIN), err)
	})

	t.Run("Create key with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_BUILDER,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, nil, nil)
		createKey, err := keySvc.CreateKey(ctx, &api.CreateKeyRequest{})
		t.Logf("createKey, err: %+v, %v", createKey, err)
		require.Nil(t, createKey)
		require.Equal(t, errPerm(api.Role_ADMIN), err)
	})

	t.Run("Create sysadmin key as non-sysadmin", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = api.Role_SYS_ADMIN

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, nil, nil)
		createKey, err := keySvc.CreateKey(ctx, &api.CreateKeyRequest{Key: key})
		t.Logf("key, createKey, err: %+v, %+v, %v", key, createKey, err)
		require.Nil(t, createKey)
		require.Equal(t, status.Error(codes.PermissionDenied, "permission "+
			"denied, role modification not allowed"), err)
	})

	t.Run("Create invalid key", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = api.Role_BUILDER

		keyer := NewMockKeyer(gomock.NewController(t))
		keyer.EXPECT().Create(gomock.Any(), key).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: key.GetOrgId(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, nil, nil)
		createKey, err := keySvc.CreateKey(ctx, &api.CreateKeyRequest{Key: key})
		t.Logf("key, createKey, err: %+v, %+v, %v", key, createKey, err)
		require.Nil(t, createKey)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("Create invalid token", func(t *testing.T) {
		t.Parallel()

		key := random.Key("api-key", uuid.NewString())
		key.Role = api.Role_BUILDER
		retKey, _ := proto.Clone(key).(*api.Key)

		keyer := NewMockKeyer(gomock.NewController(t))
		keyer.EXPECT().Create(gomock.Any(), key).Return(retKey, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, nil, nil)
		createKey, err := keySvc.CreateKey(ctx, &api.CreateKeyRequest{Key: key})
		t.Logf("key, createKey, err: %+v, %+v, %v", key, createKey, err)
		require.Nil(t, createKey)
		require.Equal(t, status.Error(codes.Unknown,
			"crypto: incorrect key length"), err)
	})
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()

	t.Run("Delete key by valid ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		keyer := NewMockKeyer(ctrl)
		keyer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)
		cacher := cache.NewMockCacher(ctrl)
		cacher.EXPECT().Set(gomock.Any(), gomock.Any(), "").Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, cacher, nil)
		_, err := keySvc.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.NoError(t, err)
	})

	t.Run("Delete key with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, nil, nil)
		_, err := keySvc.DeleteKey(ctx, &api.DeleteKeyRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_ADMIN), err)
	})

	t.Run("Delete key with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_BUILDER,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, nil, nil)
		_, err := keySvc.DeleteKey(ctx, &api.DeleteKeyRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_ADMIN), err)
	})

	t.Run("Delete key with cacher error", func(t *testing.T) {
		t.Parallel()

		cacher := cache.NewMockCacher(gomock.NewController(t))
		cacher.EXPECT().Set(gomock.Any(), gomock.Any(), "").
			Return(dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, cacher, nil)
		_, err := keySvc.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("Delete key by unknown ID", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		keyer := NewMockKeyer(ctrl)
		keyer.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(dao.ErrNotFound).Times(1)
		cacher := cache.NewMockCacher(ctrl)
		cacher.EXPECT().Set(gomock.Any(), gomock.Any(), "").Return(nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, cacher, nil)
		_, err := keySvc.DeleteKey(ctx,
			&api.DeleteKeyRequest{Id: uuid.NewString()})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.NotFound, "object not found"), err)
	})
}

func TestListKeys(t *testing.T) {
	t.Parallel()

	t.Run("List keys by valid org ID", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		keys := []*api.Key{
			random.Key("api-key", uuid.NewString()),
			random.Key("api-key", uuid.NewString()),
			random.Key("api-key", uuid.NewString()),
		}

		keyer := NewMockKeyer(gomock.NewController(t))
		keyer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(51)).
			Return(keys, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, nil, nil)
		listKeys, err := keySvc.ListKeys(ctx, &api.ListKeysRequest{})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listKeys.GetTotalSize())
		require.EqualExportedValues(t,
			&api.ListKeysResponse{Keys: keys, TotalSize: 3}, listKeys)
	})

	t.Run("List keys by valid org ID with next page", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		keys := []*api.Key{
			random.Key("api-key", uuid.NewString()),
			random.Key("api-key", uuid.NewString()),
			random.Key("api-key", uuid.NewString()),
		}

		next, err := session.GeneratePageToken(keys[1].GetCreatedAt().AsTime(),
			keys[1].GetId())
		require.NoError(t, err)

		keyer := NewMockKeyer(gomock.NewController(t))
		keyer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(keys, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, nil, nil)
		listKeys, err := keySvc.ListKeys(ctx, &api.ListKeysRequest{PageSize: 2})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listKeys.GetTotalSize())
		require.EqualExportedValues(t, &api.ListKeysResponse{
			Keys: keys[:2], NextPageToken: next, TotalSize: 3,
		}, listKeys)
	})

	t.Run("List keys with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, nil, nil)
		listKeys, err := keySvc.ListKeys(ctx, &api.ListKeysRequest{})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.Nil(t, listKeys)
		require.Equal(t, errPerm(api.Role_ADMIN), err)
	})

	t.Run("List keys by invalid page token", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, nil, nil, nil)
		listKeys, err := keySvc.ListKeys(ctx,
			&api.ListKeysRequest{PageToken: badUUID})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.Nil(t, listKeys)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"invalid page token"), err)
	})

	t.Run("List keys by invalid org ID", func(t *testing.T) {
		t.Parallel()

		keyer := NewMockKeyer(gomock.NewController(t))
		keyer.EXPECT().List(gomock.Any(), "aaa", gomock.Any(), gomock.Any(),
			gomock.Any()).Return(nil, int32(0),
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, nil, nil)
		listKeys, err := keySvc.ListKeys(ctx, &api.ListKeysRequest{})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.Nil(t, listKeys)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})

	t.Run("List keys with generation failure", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()

		keys := []*api.Key{
			random.Key("api-key", uuid.NewString()),
			random.Key("api-key", uuid.NewString()),
			random.Key("api-key", uuid.NewString()),
		}
		keys[1].Id = badUUID

		keyer := NewMockKeyer(gomock.NewController(t))
		keyer.EXPECT().List(gomock.Any(), orgID, time.Time{}, "", int32(3)).
			Return(keys, int32(3), nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			t.Context(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		keySvc := NewSession(nil, keyer, nil, nil)
		listKeys, err := keySvc.ListKeys(ctx, &api.ListKeysRequest{PageSize: 2})
		t.Logf("listKeys, err: %+v, %v", listKeys, err)
		require.NoError(t, err)
		require.Equal(t, int32(3), listKeys.GetTotalSize())
		require.EqualExportedValues(t,
			&api.ListKeysResponse{Keys: keys[:2], TotalSize: 3}, listKeys)
	})
}
