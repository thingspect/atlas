package session

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/atlas/api/go/token"
	"github.com/thingspect/atlas/pkg/crypto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// WebTokenExp represents the lifetime of a web token in seconds.
const WebTokenExp = 10 * 60

var errWebTokenExp = errors.New("crypto: token expired")

// GenerateWebToken generates an encrypted protobuf web token in raw (no
// padding) base64 format. It returns the token, expiration time, and an error
// value.
func GenerateWebToken(key []byte, userID, orgID string) (string,
	*timestamppb.Timestamp, error) {
	// Convert userID and orgID to byte slices.
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", nil, err
	}

	orgUUID, err := uuid.Parse(orgID)
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
	}, nil
}
