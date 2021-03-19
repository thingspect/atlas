package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
)

// memoryCache contains methods to create and query data in memory and
// implements the Cacher interface. An additional RWMutex is used to support
// transactions.
type memoryCache struct {
	cache   *ttlcache.Cache
	cacheMu sync.RWMutex
}

// Verify memoryCache implements Cacher.
var _ Cacher = &memoryCache{}

// NewMemory builds a new Cacher and returns it and an error value.
func NewMemory() (Cacher, error) {
	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)

	return &memoryCache{
		cache: cache,
	}, nil
}

// SetIfNotExist sets key to value if the key does not exist. If it is
// successful, it returns true.
func (m *memoryCache) SetIfNotExist(ctx context.Context, key string,
	value interface{}) (bool, error) {
	return m.SetIfNotExistTTL(ctx, key, value, ttlcache.ItemExpireWithGlobalTTL)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If it is successful, it returns true.
func (m *memoryCache) SetIfNotExistTTL(ctx context.Context, key string,
	value interface{}, exp time.Duration) (bool, error) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	if _, err := m.cache.Get(key); !errors.Is(err, ttlcache.ErrNotFound) {
		return false, err
	}

	if err := m.cache.SetWithTTL(key, value, exp); err != nil {
		return false, err
	}

	return true, nil
}

// Close closes the Cacher, releasing any open resources.
func (m *memoryCache) Close() error {
	return m.cache.Close()
}
