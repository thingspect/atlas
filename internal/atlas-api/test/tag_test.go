//go:build !unit

package test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/proto/go/api"
)

func TestListTags(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var tCount int
	var lastTag string
	for i := 0; i < 3; i++ {
		devCli := api.NewDeviceServiceClient(globalAdminGRPCConn)
		createDev, err := devCli.CreateDevice(ctx, &api.CreateDeviceRequest{
			Device: random.Device("api-tag", uuid.NewString()),
		})
		t.Logf("createDev, err: %+v, %v", createDev, err)
		require.NoError(t, err)

		tCount += len(createDev.GetTags())
		lastTag = createDev.GetTags()[len(createDev.GetTags())-1]
	}

	for i := 0; i < 3; i++ {
		user := random.User("api-user", uuid.NewString())
		user.Role = api.Role_BUILDER

		userCli := api.NewUserServiceClient(globalAdminGRPCConn)
		createUser, err := userCli.CreateUser(ctx, &api.CreateUserRequest{
			User: user,
		})
		t.Logf("createUser, err: %+v, %v", createUser, err)
		require.NoError(t, err)

		// Don't include role-based tags.
		tCount += len(createUser.GetTags()) - 1
		lastTag = createUser.GetTags()[len(createUser.GetTags())-1]
	}

	t.Run("List tags by valid org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		tagCli := api.NewTagServiceClient(globalAdminGRPCConn)
		listTags, err := tagCli.ListTags(ctx, &api.ListTagsRequest{})
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listTags.GetTags()), tCount)

		var found bool
		for _, tag := range listTags.GetTags() {
			if tag == lastTag {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Lists are isolated by org ID", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		secCli := api.NewTagServiceClient(secondaryAdminGRPCConn)
		listTags, err := secCli.ListTags(ctx, &api.ListTagsRequest{})
		t.Logf("listTags, err: %+v, %v", listTags, err)
		require.NoError(t, err)
		require.NotEmpty(t, listTags.GetTags())
	})
}
