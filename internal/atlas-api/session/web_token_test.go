//go:build !integration

package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/internal/atlas-api/crypto"
	"github.com/thingspect/atlas/pkg/test/random"
	"github.com/thingspect/atlas/proto/go/token"
	"github.com/thingspect/proto/go/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGenerateWebToken(t *testing.T) {
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
		{
			key, uuid.NewString(), uuid.NewString(), 90, "",
		},
		{
			key, random.String(10), uuid.NewString(), 0,
			"invalid UUID length: 10",
		},
		{
			key, uuid.NewString(), random.String(10), 0,
			"invalid UUID length: 10",
		},
		{
			[]byte{},
			uuid.NewString(), uuid.NewString(), 0,
			crypto.ErrKeyLength.Error(),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can generate %+v", test), func(t *testing.T) {
			t.Parallel()

			res, exp, err := GenerateWebToken(test.inpKey, &api.User{
				Id: test.inpUserID, OrgId: test.inpOrgID,
				Role: api.Role_BUILDER,
			})
			t.Logf("res, exp, err: %v, %+v, %#v", res, exp, err)
			require.GreaterOrEqual(t, len(res), test.resMinLen)
			if exp != nil {
				require.WithinDuration(t, time.Now().Add(
					WebTokenExp*time.Second), exp.AsTime(), 2*time.Second)
			}
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}

func TestGenerateKeyToken(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	tests := []struct {
		inpKey    []byte
		inpKeyID  string
		inpOrgID  string
		resMinLen int
		err       string
	}{
		{
			key, uuid.NewString(), uuid.NewString(), 80, "",
		},
		{
			key, random.String(10), uuid.NewString(), 0,
			"invalid UUID length: 10",
		},
		{
			key, uuid.NewString(), random.String(10), 0,
			"invalid UUID length: 10",
		},
		{
			[]byte{},
			uuid.NewString(), uuid.NewString(), 0,
			crypto.ErrKeyLength.Error(),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can generate %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := GenerateKeyToken(test.inpKey, test.inpKeyID,
				test.inpOrgID, api.Role_BUILDER)
			t.Logf("res, err: %v, %#v", res, err)
			require.GreaterOrEqual(t, len(res), test.resMinLen)
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}

func TestValidateWebToken(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	badCipher, err := crypto.Encrypt(key, []byte("aaa"))
	require.NoError(t, err)

	oldToken := &token.Web{ExpiresAt: timestamppb.New(time.Now().Add(-2 *
		WebTokenExp * time.Second))}
	bOldToken, err := proto.Marshal(oldToken)
	require.NoError(t, err)
	eOldToken, err := crypto.Encrypt(key, bOldToken)
	require.NoError(t, err)

	tests := []struct {
		inpKey         []byte
		inpCiphertoken string
		err            string
	}{
		{
			key, "", "",
		},
		{
			key, "...", "illegal base64 data at input byte 0",
		},
		{
			key, random.String(10), "crypto: malformed ciphertext",
		},
		{
			key, base64.RawStdEncoding.EncodeToString(badCipher),
			"cannot parse invalid wire-format data",
		},
		{
			key, base64.RawStdEncoding.EncodeToString(eOldToken),
			errWebTokenExp.Error(),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can validate %+v", test), func(t *testing.T) {
			t.Parallel()

			user := random.User("webtoken", uuid.NewString())
			resGen, exp, err := GenerateWebToken(test.inpKey, user)
			t.Logf("resGen, exp, err: %v, %v, %v", resGen, exp, err)
			require.NoError(t, err)

			var resVal *Session
			if test.inpCiphertoken == "" {
				resVal, err = ValidateWebToken(test.inpKey, resGen)
			} else {
				resVal, err = ValidateWebToken(test.inpKey,
					test.inpCiphertoken)
			}
			t.Logf("resVal, err: %+v, %v", resVal, err)
			if resVal != nil {
				require.Equal(t, user.GetId(), resVal.UserID)
				require.Empty(t, resVal.KeyID)
				require.Equal(t, user.GetOrgId(), resVal.OrgID)
				require.Equal(t, user.GetRole(), resVal.Role)
				require.NotEmpty(t, resVal.TraceID)
			}
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}

func TestValidateKeyToken(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	badCipher, err := crypto.Encrypt(key, []byte("aaa"))
	require.NoError(t, err)

	tests := []struct {
		inpKey         []byte
		inpCiphertoken string
		err            string
	}{
		{
			key, "", "",
		},
		{
			key, "...", "illegal base64 data at input byte 0",
		},
		{
			key, random.String(10), "crypto: malformed ciphertext",
		},
		{
			key, base64.RawStdEncoding.EncodeToString(badCipher),
			"cannot parse invalid wire-format data",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can validate %+v", test), func(t *testing.T) {
			t.Parallel()

			keyID := uuid.NewString()
			user := random.User("keytoken", uuid.NewString())

			resGen, err := GenerateKeyToken(test.inpKey, keyID, user.GetOrgId(),
				user.GetRole())
			t.Logf("resGen, err: %v, %v", resGen, err)
			require.NoError(t, err)

			var resVal *Session
			if test.inpCiphertoken == "" {
				resVal, err = ValidateWebToken(test.inpKey, resGen)
			} else {
				resVal, err = ValidateWebToken(test.inpKey,
					test.inpCiphertoken)
			}
			t.Logf("resVal, err: %+v, %v", resVal, err)
			if resVal != nil {
				require.Empty(t, resVal.UserID)
				require.Equal(t, keyID, resVal.KeyID)
				require.Equal(t, user.GetOrgId(), resVal.OrgID)
				require.Equal(t, user.GetRole(), resVal.Role)
				require.NotEmpty(t, resVal.TraceID)
			}
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}
