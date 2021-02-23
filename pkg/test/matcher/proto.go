// Package matcher provides types that implement the gomock.Matcher interface.
// This package must not be used outside of tests.
package matcher

import (
	"fmt"

	"github.com/golang/mock/gomock"

	// mock/mockgen/model imported for reflection-based mockgen with vendoring.
	_ "github.com/golang/mock/mockgen/model"
	"google.golang.org/protobuf/proto"
)

// protoMatcher implements the gomock.Matcher interface.
type protoMatcher struct {
	msg proto.Message
}

// Verify protoMatcher implements gomock.Matcher.
var _ gomock.Matcher = &protoMatcher{}

// NewProtoMatcher builds a new gomock.Matcher and returns it.
func NewProtoMatcher(msg proto.Message) gomock.Matcher {
	return &protoMatcher{
		msg: msg,
	}
}

// Matches returns whether x is a match.
func (pm *protoMatcher) Matches(x interface{}) bool {
	msg, ok := x.(proto.Message)
	if !ok {
		return false
	}

	return proto.Equal(pm.msg, msg)
}

// String describes what the matcher matches.
func (pm *protoMatcher) String() string {
	return fmt.Sprintf("is %+v", pm.msg)
}
