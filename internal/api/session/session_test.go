// +build !integration

package session

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewFromContext(t *testing.T) {
	t.Parallel()

	sess := &Session{UserID: uuid.New().String(), OrgID: uuid.New().String()}
	t.Logf("sess: %+v", sess)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ctx = NewContext(ctx, sess)
	ctxSess, ok := FromContext(ctx)
	t.Logf("ctxSess, ok: %+v, %v", ctxSess, ok)
	require.True(t, ok)
	require.Equal(t, sess, ctxSess)
}
