package session

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/thingspect/atlas/api/go/token"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GeneratePageToken generates a page token in raw (no padding), URL-safe base64
// format. It returns the token and an error value. The token is not currently
// encrypted due to the inclusion of only known IDs and timestamps.
func GeneratePageToken(boundTS time.Time, prevID string) (string, error) {
	// Convert lastID to []byte.
	lastUUID, err := uuid.Parse(prevID)
	if err != nil {
		return "", err
	}

	// Build page token.
	pt := &token.Page{
		BoundTs: timestamppb.New(boundTS),
		PrevId:  lastUUID[:],
	}

	bPT, err := proto.Marshal(pt)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bPT), nil
}

// ParsePageToken parses a page token in raw (no padding), URL-safe base64
// format, or an empty string. An ID, timestamp, and error value are returned.
func ParsePageToken(pToken string) (time.Time, string, error) {
	// In the case of an empty token, start from the beginning.
	if pToken == "" {
		return time.Time{}, "", nil
	}

	// Decode token.
	bPT, err := base64.RawURLEncoding.DecodeString(pToken)
	if err != nil {
		return time.Time{}, "", err
	}

	// Unmarshal page token. A nil error with missing timestamp is treated as an
	// empty token.
	pt := &token.Page{}
	if err := proto.Unmarshal(bPT, pt); err != nil || pt.BoundTs == nil {
		return time.Time{}, "", err
	}

	lastUUID, err := uuid.FromBytes(pt.PrevId)
	if err != nil {
		return time.Time{}, "", err
	}

	return pt.BoundTs.AsTime().UTC(), lastUUID.String(), nil
}
