package session

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/atlas/internal/atlas-api/auth"
	"github.com/thingspect/atlas/pkg/consterr"
	"github.com/thingspect/atlas/proto/go/token"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// WebTokenExp represents the lifetime of a web token in seconds.
	WebTokenExp = 10 * 60

	//#nosec G101 // false positive for hardcoded credentials
	errWebTokenExp consterr.Error = "crypto: token expired"
)

// GenerateWebToken generates an encrypted protobuf web token in raw (no
// padding) base64 format. It returns the token, expiration time, and an error
// value.
func GenerateWebToken(pwtKey []byte, user *api.User) (
	string, *timestamppb.Timestamp, error,
) {
	// Convert user.Id and user.OrgId to bytes.
	userUUID, err := uuid.Parse(user.GetId())
	if err != nil {
		return "", nil, err
	}

	orgUUID, err := uuid.Parse(user.GetOrgId())
	if err != nil {
		return "", nil, err
	}

	// Calculate expiration. Set exp.Nanos to zero for compactness.
	exp := timestamppb.Now()
	exp.Seconds += WebTokenExp
	exp.Nanos = 0

	// Build unencrypted PWT.
	pwt := &token.Web{
		IdOneof:   &token.Web_UserId{UserId: userUUID[:]},
		OrgId:     orgUUID[:],
		Role:      user.GetRole(),
		ExpiresAt: exp,
	}

	bPWT, err := proto.Marshal(pwt)
	if err != nil {
		return "", nil, err
	}

	// Encrypt and encode PWT.
	ePWT, err := auth.Encrypt(pwtKey, bPWT)
	if err != nil {
		return "", nil, err
	}

	return base64.RawStdEncoding.EncodeToString(ePWT), exp, nil
}

// GenerateKeyToken generates an encrypted protobuf API key token in raw (no
// padding) base64 format. It returns the token and an error value.
func GenerateKeyToken(pwtKey []byte, keyID, orgID string, role api.Role) (
	string, error,
) {
	// Convert keyID and orgID to bytes.
	keyUUID, err := uuid.Parse(keyID)
	if err != nil {
		return "", err
	}

	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		return "", err
	}

	// Build unencrypted PWT.
	pwt := &token.Web{
		IdOneof: &token.Web_KeyId{KeyId: keyUUID[:]},
		OrgId:   orgUUID[:],
		Role:    role,
	}

	bPWT, err := proto.Marshal(pwt)
	if err != nil {
		return "", err
	}

	// Encrypt and encode PWT.
	ePWT, err := auth.Encrypt(pwtKey, bPWT)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(ePWT), nil
}

// ValidateWebToken validates an encrypted protobuf web or API key token.
func ValidateWebToken(pwtKey []byte, ciphertoken string) (*Session, error) {
	// Decode and decrypt PWT.
	ePWT, err := base64.RawStdEncoding.DecodeString(ciphertoken)
	if err != nil {
		return nil, err
	}

	bPWT, err := auth.Decrypt(pwtKey, ePWT)
	if err != nil {
		return nil, err
	}

	// Unmarshal PWT.
	pwt := &token.Web{}
	if err := proto.Unmarshal(bPWT, pwt); err != nil {
		return nil, err
	}

	// Validate expiration, if present.
	if pwt.GetExpiresAt() != nil && pwt.GetExpiresAt().AsTime().Before(time.Now()) {
		return nil, errWebTokenExp
	}

	// Build Session with new TraceID. UUIDs have been authenticated and are
	// safe to copy.
	sess := &Session{
		Role:    pwt.GetRole(),
		TraceID: uuid.New(),
	}

	var idUUID uuid.UUID
	switch id := pwt.GetIdOneof().(type) {
	case *token.Web_UserId:
		_ = copy(idUUID[:], id.UserId)
		sess.UserID = idUUID.String()
	case *token.Web_KeyId:
		_ = copy(idUUID[:], id.KeyId)
		sess.KeyID = idUUID.String()
	}

	var orgUUID uuid.UUID
	_ = copy(orgUUID[:], pwt.GetOrgId())
	sess.OrgID = orgUUID.String()

	return sess, nil
}
