package session

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/api/go/token"
	"github.com/thingspect/atlas/internal/api/crypto"
	"github.com/thingspect/atlas/pkg/consterr"
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
func GenerateWebToken(key []byte, user *api.User) (string,
	*timestamppb.Timestamp, error) {
	// Convert user.Id and user.OrgId to bytes.
	userUUID, err := uuid.Parse(user.Id)
	if err != nil {
		return "", nil, err
	}

	orgUUID, err := uuid.Parse(user.OrgId)
	if err != nil {
		return "", nil, err
	}

	// Calculate expiration. Set exp.Nanos to zero for compactness.
	exp := timestamppb.Now()
	exp.Seconds += WebTokenExp
	exp.Nanos = 0

	// Build unencrypted PWT.
	pwt := &token.Web{
		UserId:    userUUID[:],
		OrgId:     orgUUID[:],
		Role:      user.Role,
		ExpiresAt: exp,
	}

	bPWT, err := proto.Marshal(pwt)
	if err != nil {
		return "", nil, err
	}

	// Encrypt and encode PWT.
	ePWT, err := crypto.Encrypt(key, bPWT)
	if err != nil {
		return "", nil, err
	}

	return base64.RawStdEncoding.EncodeToString(ePWT), exp, nil
}

// ValidateWebToken validates an encrypted protobuf web token in raw (no
// padding) base64 format. A nil error as part of the return indicates success.
func ValidateWebToken(key []byte, ciphertoken string) (*Session, error) {
	// Decode and decrypt PWT.
	ePWT, err := base64.RawStdEncoding.DecodeString(ciphertoken)
	if err != nil {
		return nil, err
	}

	bPWT, err := crypto.Decrypt(key, ePWT)
	if err != nil {
		return nil, err
	}

	// Unmarshal PWT.
	pwt := &token.Web{}
	if err := proto.Unmarshal(bPWT, pwt); err != nil {
		return nil, err
	}

	// Validate expiration.
	if pwt.ExpiresAt == nil || pwt.ExpiresAt.AsTime().Before(time.Now()) {
		return nil, errWebTokenExp
	}

	// Build Session to return. UUIDs have been decrypted and are safe to copy.
	var userUUID uuid.UUID
	copy(userUUID[:], pwt.UserId)

	var orgUUID uuid.UUID
	copy(orgUUID[:], pwt.OrgId)

	return &Session{
		UserID: userUUID.String(),
		OrgID:  orgUUID.String(),
		Role:   pwt.Role,
	}, nil
}
