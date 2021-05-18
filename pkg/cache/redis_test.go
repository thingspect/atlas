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

func TestRedisSetGet(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetGet-" + random.String(10)
	val := random.String(10)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	ok, res, err := redis.Get(ctx, key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = redis.Get(ctx, "testRedisSetGet-"+random.String(10))
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	require.NoError(t, redis.Close())

	ok, res, err = redis.Get(ctx, key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.EqualError(t, err, "redis: client is closed")
}

func TestRedisSetTTLGetB(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetTTLGetB-" + random.String(10)
	val := random.Bytes(10)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	require.NoError(t, redis.SetTTL(ctx, key, val, testTimeout))

	ok, res, err := redis.GetB(ctx, key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = redis.GetB(ctx, "testRedisSetTTLGetB-"+random.String(10))
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	require.NoError(t, redis.Close())

	ok, res, err = redis.GetB(ctx, key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.EqualError(t, err, "redis: client is closed")
}

func TestRedisSetTTLGetBShort(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetTTLGetBShort-" + random.String(10)
	val := random.Bytes(10)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	require.NoError(t, redis.SetTTL(ctx, key, val, time.Millisecond))

	time.Sleep(100 * time.Millisecond)
	ok, res, err := redis.GetB(ctx, key)
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)
}

func TestRedisSetGetI(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisSetGetI-" + random.String(10)
	val := int64(random.Intn(999))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	ok, res, err := redis.GetI(ctx, key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.True(t, ok)
	require.Equal(t, val, res)
	require.NoError(t, err)

	ok, res, err = redis.GetI(ctx, "testRedisSetGetI-"+random.String(10))
	t.Logf("ok, res, err: %v, %x, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.NoError(t, err)

	require.NoError(t, redis.Close())

	ok, res, err = redis.GetI(ctx, key)
	t.Logf("ok, res, err: %v, %v, %v", ok, res, err)
	require.False(t, ok)
	require.Empty(t, res)
	require.EqualError(t, err, "redis: client is closed")
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

func TestRedisIncr(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	key := "testRedisIncr-" + random.String(10)
	val := int64(random.Intn(999))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	require.NoError(t, redis.Set(ctx, key, val))

	res, err := redis.Incr(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, val+1, res)
	require.NoError(t, err)

	res, err = redis.Incr(ctx, "testRedisIncr-"+random.String(10))
	t.Logf("res, err: %v, %v", res, err)
	require.Equal(t, int64(1), res)
	require.NoError(t, err)

	require.NoError(t, redis.Set(ctx, key, "testRedisIncr-"+random.String(10)))

	res, err = redis.Incr(ctx, key)
	t.Logf("res, err: %v, %v", res, err)
	require.Empty(t, res)
	require.EqualError(t, err, "ERR value is not an integer or out of range")
}

func TestRedisClose(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis(testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	require.NoError(t, redis.Close())
}
