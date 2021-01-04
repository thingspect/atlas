// +build !integration

package interceptor

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/internal/api/session"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuth(t *testing.T) {
	t.Parallel()

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	token, _, err := session.GenerateToken(key, uuid.New().String(),
		uuid.New().String())
	t.Logf("token, err: %v, %v", token, err)
	require.NoError(t, err)

	skipPath := random.String(10)

	tests := []struct {
		inpMD         []string
		inpHandlerErr error
		inpSkipPaths  map[string]struct{}
		inpInfo       *grpc.UnaryServerInfo
	}{
		{[]string{"authorization", "Bearer " + token}, nil, nil,
			&grpc.UnaryServerInfo{FullMethod: random.String(10)}},
		{nil, nil, map[string]struct{}{skipPath: {}}, &grpc.UnaryServerInfo{
			FullMethod: skipPath}},
		{nil, status.Error(codes.Unauthenticated, "unauthorized"), nil,
			&grpc.UnaryServerInfo{FullMethod: random.String(10)}},
		{[]string{}, status.Error(codes.Unauthenticated, "unauthorized"), nil,
			&grpc.UnaryServerInfo{FullMethod: random.String(10)}},
		{[]string{"authorization", "NoBearer " + token},
			status.Error(codes.Unauthenticated, "unauthorized"), nil,
			&grpc.UnaryServerInfo{FullMethod: random.String(10)}},
		{[]string{"authorization", "Bearer ..."},
			status.Error(codes.Unauthenticated, "unauthorized"), nil,
			&grpc.UnaryServerInfo{FullMethod: random.String(10)}},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can log %+v", lTest), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(),
				2*time.Second)
			defer cancel()
			if lTest.inpMD != nil {
				ctx = metadata.NewIncomingContext(ctx,
					metadata.Pairs(lTest.inpMD...))
			}

			handler := func(ctx context.Context, req interface{}) (interface{},
				error) {
				return req, lTest.inpHandlerErr
			}

			res, err := Auth(lTest.inpSkipPaths, key)(ctx, nil, lTest.inpInfo,
				handler)
			t.Logf("res, err: %v, %v", res, err)
			require.Nil(t, res)
			require.Equal(t, lTest.inpHandlerErr, err)
		})
	}
}
