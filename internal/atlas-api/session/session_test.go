//go:build !integration

package session

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNewUserFromContext(t *testing.T) {
	t.Parallel()

	user := random.User("session", uuid.NewString())
	sess := &Session{
		UserID: user.GetId(), OrgID: user.GetOrgId(), Role: user.GetRole(),
		TraceID: uuid.New(),
	}
	t.Logf("sess: %+v", sess)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ctx = NewContext(ctx, sess)
	ctxSess, ok := FromContext(ctx)
	t.Logf("ctxSess, ok: %+v, %v", ctxSess, ok)
	require.True(t, ok)
	require.Equal(t, sess, ctxSess)
}

func TestNewKeyFromContext(t *testing.T) {
	t.Parallel()

	user := random.User("session", uuid.NewString())
	sess := &Session{
		KeyID: uuid.NewString(), OrgID: user.GetOrgId(), Role: user.GetRole(),
		TraceID: uuid.New(),
	}
	t.Logf("sess: %+v", sess)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ctx = NewContext(ctx, sess)
	ctxSess, ok := FromContext(ctx)
	t.Logf("ctxSess, ok: %+v, %v", ctxSess, ok)
	require.True(t, ok)
	require.Equal(t, sess, ctxSess)
}
