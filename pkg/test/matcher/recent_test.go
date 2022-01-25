//go:build !integration

//go:generate mockgen -source recent_test.go -destination mock_recenter_test.go -package matcher

package matcher

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type recenter interface {
	f(t time.Time) error
}

func runRecent(r recenter, t time.Time) error {
	return r.f(t)
}

func TestRecentMatcher(t *testing.T) {
	t.Parallel()

	recenter := NewMockrecenter(gomock.NewController(t))
	recenter.EXPECT().f(NewRecentMatcher(2 * time.Second)).Return(nil).Times(1)

	require.NoError(t, runRecent(recenter, time.Now()))
}
