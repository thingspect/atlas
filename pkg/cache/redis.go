package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
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

// SetIfNotExist sets key to value if the key does not exist. If it is
// successful, it returns true.
func (r *redisCache) SetIfNotExist(ctx context.Context, key string,
	value interface{}) (bool, error) {
	return r.SetIfNotExistTTL(ctx, key, value, 0)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If it is successful, it returns true.
func (r *redisCache) SetIfNotExistTTL(ctx context.Context, key string,
	value interface{}, exp time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, exp).Result()
}

// Close closes the Cacher, releasing any open resources.
func (r *redisCache) Close() error {
	return r.client.Close()
}
