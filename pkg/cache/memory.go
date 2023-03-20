package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ReneKroon/ttlcache/v3"
	"github.com/thingspect/atlas/pkg/consterr"
)

var errWrongType consterr.Error = "value wrong type"

// memoryCache contains methods to create and query data in memory and
// implements the Cacher interface. An additional RWMutex is used to support
// transactions.
type memoryCache struct {
	cache   *ttlcache.Cache[string, any]
	cacheMu sync.RWMutex
}

// Verify memoryCache implements Cacher.
var _ Cacher = &memoryCache{}

// NewMemory builds a new Cacher and returns it.
func NewMemory() Cacher {
	cache := ttlcache.NewCache[string, any]()
	cache.SkipTTLExtensionOnHit(true)

	return &memoryCache{cache: cache}
}

// Set sets key to value.
func (m *memoryCache) Set(_ context.Context, key string, value any) error {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	return m.cache.Set(key, value)
}

// SetTTL sets key to value with expiration.
func (m *memoryCache) SetTTL(
	_ context.Context, key string, value any, exp time.Duration,
) error {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	return m.cache.SetWithTTL(key, value, exp)
}

// Get retrieves a string value by key. If the key does not exist, the boolean
// returned is set to false.
func (m *memoryCache) Get(_ context.Context, key string) (bool, string, error) {
	m.cacheMu.RLock()
	defer m.cacheMu.RUnlock()

	iface, err := m.cache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	s, ok := iface.(string)
	if !ok {
		return false, "", errWrongType
	}

	return true, s, nil
}

// GetB retrieves a []byte value by key. If the key does not exist, the boolean
// returned is set to false.
func (m *memoryCache) GetB(_ context.Context, key string) (
	bool, []byte, error,
) {
	m.cacheMu.RLock()
	defer m.cacheMu.RUnlock()

	iface, err := m.cache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}

	b, ok := iface.([]byte)
	if !ok {
		return false, nil, errWrongType
	}

	return true, b, nil
}

// GetI retrieves an int64 value by key. If the key does not exist, the boolean
// returned is set to false.
func (m *memoryCache) GetI(_ context.Context, key string) (
	bool, int64, error,
) {
	m.cacheMu.RLock()
	defer m.cacheMu.RUnlock()

	iface, err := m.cache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	i, ok := iface.(int64)
	if !ok {
		return false, 0, errWrongType
	}

	return true, i, nil
}

// SetIfNotExist sets key to value if the key does not exist. If it is
// successful, it returns true.
func (m *memoryCache) SetIfNotExist(
	ctx context.Context, key string, value any,
) (bool, error) {
	return m.SetIfNotExistTTL(ctx, key, value, ttlcache.ItemExpireWithGlobalTTL)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If it is successful, it returns true.
func (m *memoryCache) SetIfNotExistTTL(
	_ context.Context, key string, value any, exp time.Duration,
) (bool, error) {
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

// Incr increments an int64 value at key by one. If the key does not exist, the
// value is set to 1. The incremented value is returned.
func (m *memoryCache) Incr(_ context.Context, key string) (int64, error) {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	iface, err := m.cache.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		iface = int64(0)
	} else if err != nil {
		return 0, err
	}

	i, ok := iface.(int64)
	if !ok {
		return 0, errWrongType
	}
	i++

	if err = m.cache.Set(key, i); err != nil {
		return 0, err
	}

	return i, nil
}

// Del removes the specified key. A key is ignored if it does not exist.
func (m *memoryCache) Del(_ context.Context, key string) error {
	m.cacheMu.Lock()
	defer m.cacheMu.Unlock()

	if err := m.cache.Remove(key); err != nil &&
		!errors.Is(err, ttlcache.ErrNotFound) {
		return err
	}

	return nil
}

// Close closes the Cacher, releasing any open resources.
func (m *memoryCache) Close() error {
	return m.cache.Close()
}
