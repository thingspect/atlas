package session

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/atlas/api/go/pwt"
	"github.com/thingspect/atlas/pkg/crypto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TokenExp represents the lifetime of a token in seconds.
const TokenExp = 10 * 60

var ErrTokenExp = errors.New("crypto: token expired")

// GenerateToken generates a protobuf web token in raw (no padding) base64
// format. It returns the token, expiration time, and an error value.
func GenerateToken(key []byte, userID, orgID string) (string,
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
	exp.Seconds += TokenExp
	exp.Nanos = 0

	// Build unencrypted PWT.
	token := &pwt.Claim{
		UserId:    userUUID[:],
		OrgId:     orgUUID[:],
		ExpiresAt: exp,
	}

	bToken, err := proto.Marshal(token)
	if err != nil {
		return "", nil, err
	}

	// Encrypt and encode PWT.
	eToken, err := crypto.Encrypt(key, bToken)
	if err != nil {
		return "", nil, err
	}
	return base64.RawStdEncoding.EncodeToString(eToken), exp, nil
}

// ValidateToken validates a protobuf web token in raw (no padding) base64
// format. A nil error as part of the return indicates success.
func ValidateToken(key []byte, ciphertoken string) (*Session, error) {
	// Decode and decrypt PWT.
	eToken, err := base64.RawStdEncoding.DecodeString(ciphertoken)
	if err != nil {
		return nil, err
	}

	bToken, err := crypto.Decrypt(key, eToken)
	if err != nil {
		return nil, err
	}

	// Unmarshal PWT.
	token := &pwt.Claim{}
	if err := proto.Unmarshal(bToken, token); err != nil {
		return nil, err
	}

	// Validate expiration.
	if token.ExpiresAt == nil || token.ExpiresAt.AsTime().Before(time.Now()) {
		return nil, ErrTokenExp
	}

	// Build Session to return.
	var userUUID uuid.UUID
	copy(userUUID[:], token.UserId)

	var orgUUID uuid.UUID
	copy(orgUUID[:], token.OrgId)

	return &Session{
		UserID: userUUID.String(),
		OrgID:  orgUUID.String(),
	}, nil
}
