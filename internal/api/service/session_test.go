// +build !integration

package service

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("Log in valid user", func(t *testing.T) {
		t.Parallel()

		orgName := random.String(10)
		user := &api.User{Id: uuid.New().String(), OrgId: uuid.New().String(),
			Email: random.Email()}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().ReadByEmail(gomock.Any(), user.Email, orgName).
			Return(user, globalHash, nil).Times(1)

		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		sessSvc := NewSession(userer, key)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.Email, OrgName: orgName, Password: globalPass})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(loginResp.Token), 90)
		require.WithinDuration(t, time.Now().Add(session.TokenExp*time.Second),
			loginResp.ExpiresAt.AsTime(), 2*time.Second)
	})

	t.Run("Log in unknown user", func(t *testing.T) {
		t.Parallel()

		orgName := random.String(10)
		user := &api.User{Id: uuid.New().String(), OrgId: uuid.New().String(),
			Email: random.Email()}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().ReadByEmail(gomock.Any(), user.Email, orgName).
			Return(nil, nil, dao.ErrNotFound).Times(1)

		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		sessSvc := NewSession(userer, key)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.Email, OrgName: orgName, Password: globalPass})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})

	t.Run("Log in wrong password", func(t *testing.T) {
		t.Parallel()

		orgName := random.String(10)
		user := &api.User{Id: uuid.New().String(), OrgId: uuid.New().String(),
			Email: random.Email()}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().ReadByEmail(gomock.Any(), user.Email, orgName).
			Return(user, globalHash, nil).Times(1)

		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		sessSvc := NewSession(userer, key)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.Email, OrgName: orgName, Password: random.String(10)})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})

	t.Run("Log in disabled user", func(t *testing.T) {
		t.Parallel()

		orgName := random.String(10)
		user := &api.User{Id: uuid.New().String(), OrgId: uuid.New().String(),
			Email: random.Email(), IsDisabled: true}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().ReadByEmail(gomock.Any(), user.Email, orgName).
			Return(user, globalHash, nil).Times(1)

		key := make([]byte, 32)
		_, err := rand.Read(key)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		sessSvc := NewSession(userer, key)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.Email, OrgName: orgName, Password: globalPass})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})

	t.Run("Log in wrong key", func(t *testing.T) {
		t.Parallel()

		orgName := random.String(10)
		user := &api.User{Id: uuid.New().String(), OrgId: uuid.New().String(),
			Email: random.Email()}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userer := NewMockUserer(ctrl)
		userer.EXPECT().ReadByEmail(gomock.Any(), user.Email, orgName).
			Return(user, globalHash, nil).Times(1)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		sessSvc := NewSession(userer, nil)
		loginResp, err := sessSvc.Login(ctx, &api.LoginRequest{
			Email: user.Email, OrgName: orgName, Password: globalPass})
		t.Logf("loginResp, err: %+v, %v", loginResp, err)
		require.Nil(t, loginResp)
		require.Equal(t, status.Error(codes.Unauthenticated, "unauthorized"),
			err)
	})
}
