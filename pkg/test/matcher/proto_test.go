// +build !integration

//go:generate mockgen -source proto_test.go -destination mock_protoer_test.go -package matcher

package matcher

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/api/go/message"
)

type protoer interface {
	f(vIn *message.ValidatorIn) error
}

func run(p protoer, vIn *message.ValidatorIn) error {
	return p.f(vIn)
}

func TestProtoMatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	protoer := NewMockprotoer(ctrl)
	protoer.EXPECT().f(NewProtoMatcher(&message.ValidatorIn{})).Return(
		nil).Times(1)

	require.NoError(t, run(protoer, &message.ValidatorIn{}))
}
