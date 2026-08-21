//go:build !integration

package session

import (
	"context"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestNewUserFromContext(t *testing.T) {
	t.Parallel()

	user := random.User("session", uuid.NewV7().String())
	sess := &Session{
		UserID: user.GetId(), OrgID: user.GetOrgId(), Role: user.GetRole(),
		TraceID: uuid.NewV7(),
	}
	t.Logf("sess: %+v", sess)

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	ctx = NewContext(ctx, sess)
	ctxSess, ok := FromContext(ctx)
	t.Logf("ctxSess, ok: %+v, %v", ctxSess, ok)
	require.True(t, ok)
	require.Equal(t, sess, ctxSess)
}

func TestNewKeyFromContext(t *testing.T) {
	t.Parallel()

	user := random.User("session", uuid.NewV7().String())
	sess := &Session{
		KeyID: uuid.NewV7().String(), OrgID: user.GetOrgId(), Role: user.GetRole(),
		TraceID: uuid.NewV7(),
	}
	t.Logf("sess: %+v", sess)

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	ctx = NewContext(ctx, sess)
	ctxSess, ok := FromContext(ctx)
	t.Logf("ctxSess, ok: %+v, %v", ctxSess, ok)
	require.True(t, ok)
	require.Equal(t, sess, ctxSess)
}
