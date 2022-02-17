// Package session provides functions for creating, retrieving, and validating
// sessions and tokens for authentication.
package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
)

// Session represents session metadata as retrieved from an encrypted token.
// Either UserID or KeyID will be present, but not both. TraceID will be
// generated whenever a token is validated and a session is returned.
type Session struct {
	UserID  string
	KeyID   string
	OrgID   string
	Role    api.Role
	TraceID uuid.UUID
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
