// Package session provides functions for creating, retrieving, and validating
// sessions and tokens for authentication.
package session

import (
	"context"

	"github.com/thingspect/api/go/common"
)

// Session represents user information as retrieved from an encrypted web token.
type Session struct {
	UserID string
	OrgID  string
	Role   common.Role
}

// sessionKey is the key for Session values in Contexts. It is unexported,
// clients should use NewContext and FromContext instead of using this key
// directly.
type sessionKey struct{}

// NewContext returns a new Context that carries a Session.
func NewContext(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, sessionKey{}, sess)
}

// FromContext returns the Session value stored in a Context, if any.
func FromContext(ctx context.Context) (*Session, bool) {
	sess, ok := ctx.Value(sessionKey{}).(*Session)

	return sess, ok
}
