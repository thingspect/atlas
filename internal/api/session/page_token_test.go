// +build !integration

package session

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/api/go/token"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGeneratePageToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpTS     time.Time
		inpID     string
		resMinLen int
		err       string
	}{
		{time.Now(), uuid.NewString(), 40, ""},
		{time.Time{}, random.String(10), 0, "invalid UUID length: 10"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can generate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := GeneratePageToken(lTest.inpTS, lTest.inpID)
			t.Logf("res, err: %v, %#v", res, err)
			require.GreaterOrEqual(t, len(res), lTest.resMinLen)
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}

func TestParsePageToken(t *testing.T) {
	t.Parallel()

	prevID := uuid.New()

	nilTSPT := &token.Page{PrevId: prevID[:]}
	bNilTSPT, err := proto.Marshal(nilTSPT)
	require.NoError(t, err)

	badUUIDPT := &token.Page{BoundTs: timestamppb.Now(), PrevId: []byte("aaa")}
	bBadUUIDPT, err := proto.Marshal(badUUIDPT)
	require.NoError(t, err)

	tests := []struct {
		inpID string
		inpTS time.Time
		inpPT string
		err   string
	}{
		{prevID.String(), time.Now().UTC(), "res", ""},
		{prevID.String(), time.Time{}, "", ""},
		{prevID.String(), time.Time{}, "...",
			"illegal base64 data at input byte 0"},
		{prevID.String(), time.Time{}, base64.RawURLEncoding.EncodeToString(
			[]byte("aaa")), "unexpected EOF"},
		{prevID.String(), time.Time{}, base64.RawURLEncoding.EncodeToString(
			bNilTSPT), ""},
		{prevID.String(), time.Time{}, base64.RawURLEncoding.EncodeToString(
			bBadUUIDPT), "invalid UUID (got 3 bytes)"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can generate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			resGen, err := GeneratePageToken(lTest.inpTS, lTest.inpID)
			t.Logf("resGen, err: %v, %#v", resGen, err)
			require.NoError(t, err)

			var resTS time.Time
			var resID string
			if lTest.inpPT == "res" {
				resTS, resID, err = ParsePageToken(resGen)
			} else {
				resTS, resID, err = ParsePageToken(lTest.inpPT)
			}
			t.Logf("resTS, resID, err: %v, %v, %#v", resTS, resID, err)
			if resID != "" {
				require.Equal(t, lTest.inpID, resID)
				require.Equal(t, lTest.inpTS, resTS)
			}
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
