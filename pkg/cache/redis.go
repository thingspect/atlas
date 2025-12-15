package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache contains methods to create and query data in Redis and implements
// the Cacher interface.
type redisCache[V Cacheable] struct {
	client redis.UniversalClient
}

// Verify redisCache implements Cacher.
var _ Cacher[string] = &redisCache[string]{}

// NewRedis builds and verifies a new Cacher and returns it and an error value.
func NewRedis[V Cacheable](redisAddr string) (Cacher[V], error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{redisAddr},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisCache[V]{
		client: client,
	}, nil
}

// Set sets key to value.
func (r *redisCache[V]) Set(ctx context.Context, key string, value V) error {
	return r.SetTTL(ctx, key, value, 0)
}

// SetTTL sets key to value with expiration.
func (r *redisCache[V]) SetTTL(ctx context.Context, key string, value V,
	exp time.Duration,
) error {
	// Switching on - and converting back to - type parameters is not supported.
	// Use JSON bytes as a workaround. This has the added benefit of supporting
	// primitives as understood by Redis.
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, b, exp).Err()
}

// Get retrieves a value by key. If the key does not exist, ErrNotFound is
// returned.
func (r *redisCache[V]) Get(ctx context.Context, key string) (V, error) {
	b, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return *new(V), ErrNotFound
	}
	if err != nil {
		return *new(V), err
	}

	var item V
	if err = json.Unmarshal(b, &item); err != nil {
		return *new(V), err
	}

	return item, nil
}

// SetIfNotExist sets key to value if the key does not exist. If the key already
// exists, ErrAlreadyExists is returned.
func (r *redisCache[V]) SetIfNotExist(ctx context.Context, key string,
	value V,
) error {
	return r.SetIfNotExistTTL(ctx, key, value, 0)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If the key already exists, ErrAlreadyExists is returned.
func (r *redisCache[V]) SetIfNotExistTTL(ctx context.Context, key string,
	value V, exp time.Duration,
) error {
	// Switching on - and converting back to - type parameters is not supported.
	// Use JSON bytes as a workaround. This has the added benefit of supporting
	// primitives as understood by Redis.
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ok, err := r.client.SetNX(ctx, key, b, exp).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrAlreadyExists
	}

	return nil
}

// Incr increments an int64 value at key by one. If the key does not exist, the
// value is set to 1. The incremented value is returned.
func (r *redisCache[V]) Incr(ctx context.Context, key string) (int64, error) {
	i, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return i, nil
}

// Del removes the specified key. A key is ignored if it does not exist.
func (r *redisCache[V]) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close closes the Cacher, releasing any open resources.
func (r *redisCache[V]) Close() error {
	return r.client.Close()
}
