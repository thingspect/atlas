//go:build !integration

//go:generate mockgen -source proto_test.go -destination mock_protoer_test.go -package matcher

package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/api/go/token"
	"go.uber.org/mock/gomock"
)

type protoer interface {
	f(vIn *token.Web) error
}

func runProto(p protoer, vIn *token.Web) error {
	return p.f(vIn)
}

func TestProtoMatcher(t *testing.T) {
	t.Parallel()

	protoer := NewMockprotoer(gomock.NewController(t))
	protoer.EXPECT().f(NewProtoMatcher(&token.Web{})).Return(nil).Times(1)

	require.NoError(t, runProto(protoer, &token.Web{}))
}
