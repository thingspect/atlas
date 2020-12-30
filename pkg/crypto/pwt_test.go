// +build !integration

package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/api/go/pwt"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGenerateToken(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	tests := []struct {
		inpKey    []byte
		inpUserID string
		inpOrgID  string
		resMinLen int
		err       string
	}{
		{key, uuid.New().String(), uuid.New().String(), 90, ""},
		{key, random.String(10), uuid.New().String(), 0,
			"invalid UUID length: 10"},
		{key, uuid.New().String(), random.String(10), 0,
			"invalid UUID length: 10"},
		{[]byte{}, uuid.New().String(), uuid.New().String(), 0,
			ErrKeyLength.Error()},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can generate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, exp, err := GenerateToken(lTest.inpKey, lTest.inpUserID,
				lTest.inpOrgID)
			t.Logf("res, exp, err: %v, %+v, %#v", res, exp, err)
			require.GreaterOrEqual(t, len(res), lTest.resMinLen)
			if exp != nil {
				require.WithinDuration(t, time.Now().Add(tokenExp*time.Second),
					exp.AsTime(), 2*time.Second)
			}
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	badCipher, err := Encrypt(key, []byte("aaa"))
	require.NoError(t, err)

	oldToken := &pwt.Claim{ExpiresAt: timestamppb.New(time.Now().Add(-2 *
		tokenExp * time.Second))}
	bOldToken, err := proto.Marshal(oldToken)
	require.NoError(t, err)
	eOldToken, err := Encrypt(key, bOldToken)
	require.NoError(t, err)

	tests := []struct {
		inpKey         []byte
		inpCiphertoken string
		err            string
	}{
		{key, "", ""},
		{key, "...", "illegal base64 data at input byte 0"},
		{key, random.String(10), "crypto: malformed ciphertext"},
		{key, base64.RawStdEncoding.EncodeToString(badCipher),
			"unexpected EOF"},
		{key, base64.RawStdEncoding.EncodeToString(eOldToken),
			ErrTokenExp.Error()},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can validate %+v", lTest), func(t *testing.T) {
			t.Parallel()

			resGen, exp, err := GenerateToken(lTest.inpKey, uuid.New().String(),
				uuid.New().String())
			t.Logf("resGen, exp, err: %v, %v, %v", resGen, exp, err)
			require.NoError(t, err)

			var resVal *pwt.Claim
			if lTest.inpCiphertoken == "" {
				resVal, err = ValidateToken(lTest.inpKey, resGen)
			} else {
				resVal, err = ValidateToken(lTest.inpKey, lTest.inpCiphertoken)
			}
			t.Logf("resVal, err: %+v, %v", resVal, err)
			if resVal != nil {
				require.WithinDuration(t, time.Now().Add(tokenExp*time.Second),
					resVal.ExpiresAt.AsTime(), 2*time.Second)
			}
			if lTest.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, lTest.err)
			}
		})
	}
}
