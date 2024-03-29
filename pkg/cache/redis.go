package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache contains methods to create and query data in Redis and implements
// the Cacher interface.
type redisCache struct {
	client redis.UniversalClient
}

// Verify redisCache implements Cacher.
var _ Cacher = &redisCache{}

// NewRedis builds and verifies a new Cacher and returns it and an error value.
func NewRedis(redisAddr string) (Cacher, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{redisAddr},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisCache{
		client: client,
	}, nil
}

// Set sets key to value.
func (r *redisCache) Set(ctx context.Context, key string, value any) error {
	return r.SetTTL(ctx, key, value, 0)
}

// SetTTL sets key to value with expiration.
func (r *redisCache) SetTTL(
	ctx context.Context, key string, value any, exp time.Duration,
) error {
	return r.client.Set(ctx, key, value, exp).Err()
}

// Get retrieves a string value by key. If the key does not exist, the boolean
// returned is set to false.
func (r *redisCache) Get(ctx context.Context, key string) (
	bool, string, error,
) {
	s, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, s, nil
}

// GetB retrieves a []byte value by key. If the key does not exist, the boolean
// returned is set to false.
func (r *redisCache) GetB(ctx context.Context, key string) (
	bool, []byte, error,
) {
	b, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}

	return true, b, nil
}

// GetI retrieves an int64 value by key. If the key does not exist, the boolean
// returned is set to false.
func (r *redisCache) GetI(ctx context.Context, key string) (
	bool, int64, error,
) {
	i, err := r.client.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	return true, i, nil
}

// SetIfNotExist sets key to value if the key does not exist. If it is
// successful, it returns true.
func (r *redisCache) SetIfNotExist(ctx context.Context, key string, value any) (
	bool, error,
) {
	return r.SetIfNotExistTTL(ctx, key, value, 0)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If it is successful, it returns true.
func (r *redisCache) SetIfNotExistTTL(
	ctx context.Context, key string, value any, exp time.Duration,
) (bool, error) {
	return r.client.SetNX(ctx, key, value, exp).Result()
}

// Incr increments an int64 value at key by one. If the key does not exist, the
// value is set to 1. The incremented value is returned.
func (r *redisCache) Incr(ctx context.Context, key string) (int64, error) {
	i, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return i, nil
}

// Del removes the specified key. A key is ignored if it does not exist.
func (r *redisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close closes the Cacher, releasing any open resources.
func (r *redisCache) Close() error {
	return r.client.Close()
}
