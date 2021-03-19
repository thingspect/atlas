// +build !unit

package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thingspect/atlas/pkg/test/config"
	"github.com/thingspect/atlas/pkg/test/random"
)

const testTimeout = 2 * time.Second

func TestNewRedis(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	tests := []struct {
		inp string
		err string
	}{
		// Success.
		{testConfig.RedisHost + ":6379", ""},
		// Wrong port.
		{"127.0.0.1:6380", "connect: connection refused"},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can connect %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res, err := NewRedis(lTest.inp)
			t.Logf("res, err: %+v, %v", res, err)
			if lTest.err == "" {
				require.NotNil(t, res)
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), lTest.err)
			}
		})
	}
}

func TestRedisSetIfNotExist(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetIfNotExist-" + random.String(10)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	ok, err := redis.SetIfNotExist(ctx, key, random.Bytes(10))
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	ok, err = redis.SetIfNotExist(ctx, key, random.Bytes(10))
	t.Logf("ok, err: %v, %v", ok, err)
	require.False(t, ok)
	require.NoError(t, err)
}

func TestRedisSetIfNotExistTTL(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetIfNotExistTTL-" + random.String(10)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	ok, err := redis.SetIfNotExistTTL(ctx, key, random.Bytes(10), testTimeout)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	ok, err = redis.SetIfNotExistTTL(ctx, key, random.Bytes(10), testTimeout)
	t.Logf("ok, err: %v, %v", ok, err)
	require.False(t, ok)
	require.NoError(t, err)
}

func TestRedisSetIfNotExistTTLShort(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetIfNotExistTTLShort-" + random.String(10)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	ok, err := redis.SetIfNotExistTTL(ctx, key, random.Bytes(10),
		time.Millisecond)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	ok, err = redis.SetIfNotExistTTL(ctx, key, random.Bytes(10), testTimeout)
	t.Logf("ok, err: %v, %v", ok, err)
	require.True(t, ok)
	require.NoError(t, err)
}

func TestRedisClose(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	require.NoError(t, redis.Close())
}
