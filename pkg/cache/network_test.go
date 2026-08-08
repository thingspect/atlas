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
		{"127.0.0.1:6399", "connect: connection refused"},
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

func TestNewValkey(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	tests := []struct {
		inp string
		err string
	}{
		// Success.
		{testConfig.ValkeyHost + ":" + testConfig.ValkeyPort, ""},
		// Wrong port.
		{"127.0.0.1:6399", "connect: connection refused"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can connect %+v", test), func(t *testing.T) {
			t.Parallel()

			res, err := NewValkey[string](test.inp)
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

func TestNetworkSetGetStringClose(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[string](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[string](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkSetGetString-" + random.String(10)
	val := random.String(10)

	for _, net := range []Cacher[string]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.Set(ctx, key, val))

			res, err := net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			res, err = net.Get(ctx, "testNetworkSetGetString-"+
				random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)

			require.NoError(t, net.Close())

			res, err = net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Contains(t, err.Error(), "client is clos")
		})
	}
}

func TestNetworkSetTTLGetBytesClose(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[[]byte](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkSetTTLGetBytes-" + random.String(10)
	val := random.Bytes(10)

	for _, net := range []Cacher[[]byte]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.SetTTL(ctx, key, val, testTimeout))

			res, err := net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			res, err = net.Get(ctx, "testNetworkSetTTLGetBytes-"+
				random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)

			require.NoError(t, net.Close())

			res, err = net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Contains(t, err.Error(), "client is clos")
		})
	}
}

func TestNetworkSetTTLGetBytesShort(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[[]byte](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkSetTTLGetBytesShort-" + random.String(10)
	val := random.Bytes(10)

	for _, net := range []Cacher[[]byte]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.SetTTL(ctx, key, val, time.Millisecond))

			time.Sleep(10 * time.Millisecond)
			res, err := net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)
		})
	}
}

func TestNetworkSetGetInt64Close(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[int64](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[int64](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkSetGetInt64-" + random.String(10)
	val := int64(random.Intn(999))

	for _, net := range []Cacher[int64]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.Set(ctx, key, val))

			res, err := net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			res, err = net.Get(ctx, "testNetworkSetGetInt64-"+random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)

			require.NoError(t, net.Close())

			res, err = net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Contains(t, err.Error(), "client is clos")
		})
	}
}

func TestNetworkSetIfNotExistBytesClose(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[[]byte](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkSetIfNotExistBytes-" + random.String(10)

	for _, net := range []Cacher[[]byte]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.SetIfNotExist(ctx, key, random.Bytes(10)))

			require.Equal(t, ErrAlreadyExists, net.SetIfNotExist(ctx, key,
				random.Bytes(10)))

			require.NoError(t, net.Close())

			err := net.SetIfNotExist(ctx, key, random.Bytes(10))
			t.Logf("err: %v", err)
			require.Contains(t, err.Error(), "client is clos")
		})
	}
}

func TestNetworkSetIfNotExistTTLBytes(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[[]byte](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkSetIfNotExistTTLBytes-" + random.String(10)

	for _, net := range []Cacher[[]byte]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.SetIfNotExistTTL(ctx, key, random.Bytes(10),
				testTimeout))

			require.Equal(t, ErrAlreadyExists, net.SetIfNotExistTTL(ctx, key,
				random.Bytes(10), testTimeout))
		})
	}
}

func TestNetworkSetIfNotExistTTLBytesShort(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[[]byte](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[[]byte](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkSetIfNotExistTTLBytesShort-" + random.String(10)

	for _, net := range []Cacher[[]byte]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.SetIfNotExistTTL(ctx, key, random.Bytes(10),
				time.Millisecond))

			time.Sleep(10 * time.Millisecond)
			require.NoError(t, net.SetIfNotExistTTL(ctx, key, random.Bytes(10),
				testTimeout))
		})
	}
}

func TestNetworkIncrInt64(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[int64](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[int64](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkIncrInt64-" + random.String(10)
	val := int64(random.Intn(999))

	for _, net := range []Cacher[int64]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.Set(ctx, key, val))

			res, err := net.Incr(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val+1, res)
			require.NoError(t, err)

			res, err = net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val+1, res)
			require.NoError(t, err)

			res, err = net.Incr(ctx, "testNetworkIncrInt64-"+random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, int64(1), res)
			require.NoError(t, err)
		})
	}
}

func TestNetworkIncrString(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[string](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[string](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkIncrString-" + random.String(10)
	val := random.String(10)

	for _, net := range []Cacher[string]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.Set(ctx, key, val))

			res, err := net.Incr(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Contains(t, err.Error(),
				"value is not an integer or out of range")

			res, err = net.Incr(ctx, "testNetworkIncrString-"+random.String(10))
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, int64(1), res)
			require.NoError(t, err)
		})
	}
}

func TestNetworkDelString(t *testing.T) {
	t.Parallel()

	testConfig := config.New()

	redis, err := NewRedis[string](testConfig.RedisHost + ":6379")
	t.Logf("redis, err: %+v, %v", redis, err)
	require.NoError(t, err)

	valkey, err := NewValkey[string](testConfig.ValkeyHost + ":" +
		testConfig.ValkeyPort)
	t.Logf("valkey, err: %+v, %v", valkey, err)
	require.NoError(t, err)

	key := "testNetworkDelString-" + random.String(10)
	val := random.String(10)

	for _, net := range []Cacher[string]{redis, valkey} {
		t.Run(fmt.Sprintf("Can set get %+v", net), func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
			defer cancel()

			require.NoError(t, net.Set(ctx, key, val))

			res, err := net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Equal(t, val, res)
			require.NoError(t, err)

			err = net.Del(ctx, key)
			t.Logf("err: %v", err)
			require.NoError(t, err)

			res, err = net.Get(ctx, key)
			t.Logf("res, err: %v, %v", res, err)
			require.Empty(t, res)
			require.Equal(t, ErrNotFound, err)

			err = net.Del(ctx, "testNetworkDelString-"+random.String(10))
			t.Logf("err: %v", err)
			require.NoError(t, err)
		})
	}
}
