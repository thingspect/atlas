// +build !integration

package session

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/random"
)

func TestGeneratePageToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inpTS     time.Time
		inpID     string
		resMinLen int
		err       string
	}{
		{time.Now(), uuid.New().String(), 40, ""},
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

	tests := []struct {
		inpID string
		inpTS time.Time
		inpPT string
		err   string
	}{
		{uuid.New().String(), time.Now().UTC(), "res", ""},
		{uuid.New().String(), time.Time{}, "", ""},
		{uuid.New().String(), time.Time{}, "...",
			"illegal base64 data at input byte 0"},
		{uuid.New().String(), time.Time{}, base64.RawURLEncoding.EncodeToString(
			[]byte("aaa")), "unexpected EOF"},
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
