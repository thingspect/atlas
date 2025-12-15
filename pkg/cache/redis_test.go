//go:build !unit

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

const testTimeout = 5 * time.Second

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
		t.Run(fmt.Sprintf("Can connect %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := NewRedis[string](test.inp)
			t.Logf("res, err: %+v, %v", res, err)
			if test.err == "" {
				require.NotNil(t, res)
				require.NoError(t, err)
			} else {
				require.Contains(t, err.Error(), test.err)
			}
		})
	}
}

func TestRedisSetGetString(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[string](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetGetString-" + random.String(10)
	val := random.String(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	res, err := redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, val, res)
	require.NoError(t, err)

	res, err = redis.Get(ctx, "testRedisSetGetString-"+random.String(10))
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.Equal(t, ErrNotFound, err)

	require.NoError(t, redis.Close())

	res, err = redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.EqualError(t, err, "redis: client is closed")
}

func TestRedisSetTTLGetBytes(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetTTLGetBytes-" + random.String(10)
	val := random.Bytes(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.SetTTL(ctx, key, val, testTimeout))

	res, err := redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, val, res)
	require.NoError(t, err)

	res, err = redis.Get(ctx, "testRedisSetTTLGetBytes-"+random.String(10))
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.Equal(t, ErrNotFound, err)

	require.NoError(t, redis.Close())

	res, err = redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.EqualError(t, err, "redis: client is closed")
}

func TestRedisSetTTLGetBytesShort(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetTTLGetBytesShort-" + random.String(10)
	val := random.Bytes(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.SetTTL(ctx, key, val, time.Millisecond))

	time.Sleep(10 * time.Millisecond)
	res, err := redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.Equal(t, ErrNotFound, err)
}

func TestRedisSetGetInt64(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[int64](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetGetInt64-" + random.String(10)
	val := int64(random.Intn(999))

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	res, err := redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, val, res)
	require.NoError(t, err)

	res, err = redis.Get(ctx, "testRedisSetGetInt64-"+random.String(10))
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.Equal(t, ErrNotFound, err)

	require.NoError(t, redis.Close())

	res, err = redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.EqualError(t, err, "redis: client is closed")
}

func TestRedisSetIfNotExistBytes(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetIfNotExistBytes-" + random.String(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.SetIfNotExist(ctx, key, random.Bytes(10)))

	require.Equal(t, ErrAlreadyExists, redis.SetIfNotExist(ctx, key,
		random.Bytes(10)))
}

func TestRedisSetIfNotExistTTLBytes(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetIfNotExistTTLBytes-" + random.String(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.SetIfNotExistTTL(ctx, key, random.Bytes(10),
		testTimeout))

	require.Equal(t, ErrAlreadyExists, redis.SetIfNotExistTTL(ctx, key,
		random.Bytes(10), testTimeout))
}

func TestRedisSetIfNotExistTTLBytesShort(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetIfNotExistTTLBytesShort-" + random.String(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.SetIfNotExistTTL(ctx, key, random.Bytes(10),
		time.Millisecond))

	time.Sleep(10 * time.Millisecond)
	require.NoError(t, redis.SetIfNotExistTTL(ctx, key, random.Bytes(10),
		testTimeout))
}

func TestRedisIncrInt64(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[int64](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisIncrInt64-" + random.String(10)
	val := int64(random.Intn(999))

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	res, err := redis.Incr(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, val+1, res)
	require.NoError(t, err)

	res, err = redis.Incr(ctx, "testRedisIncrInt64-"+random.String(10))
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, int64(1), res)
	require.NoError(t, err)
}

func TestRedisIncrString(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[string](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisIncrString-" + random.String(10)
	val := random.String(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	res, err := redis.Incr(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.EqualError(t, err, "ERR value is not an integer or out of range")

	res, err = redis.Incr(ctx, "testRedisIncrString-"+random.String(10))
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, int64(1), res)
	require.NoError(t, err)
}

func TestRedisDelString(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[string](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisDelString-" + random.String(10)
	val := random.String(10)

	ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	res, err := redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, val, res)
	require.NoError(t, err)

	err = redis.Del(ctx, key)
	t.Logf("err: %v", err)
	require.NoError(t, err)

	res, err = redis.Get(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.Equal(t, ErrNotFound, err)
}

func TestRedisClose(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[string](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	require.NoError(t, redis.Close())
}
