package cache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/thingspect/atlas/pkg/consterr"
)

// ErrWrongType is returned when a stored value does not match the type of a
// called method.
const ErrWrongType consterr.Error = "value wrong type"

// memoryCache contains methods to create and query data in memory and
// implements the Cacher interface.
type memoryCache struct {
	cache *ttlcache.Cache[string, any]
}

// Verify memoryCache implements Cacher.
var _ Cacher = &memoryCache{}

// NewMemory builds a new Cacher and returns it.
func NewMemory() Cacher {
	cache := ttlcache.New(
		ttlcache.WithDisableTouchOnHit[string, any](),
	)

	go cache.Start()

	return &memoryCache{cache: cache}
}

// Set sets key to value.
func (m *memoryCache) Set(ctx context.Context, key string, value any) error {
	return m.SetTTL(ctx, key, value, ttlcache.NoTTL)
}

// SetTTL sets key to value with expiration.
func (m *memoryCache) SetTTL(
	_ context.Context, key string, value any, exp time.Duration,
) error {
	_ = m.cache.Set(key, value, exp)

	return nil
}

// Get retrieves a string value by key. If the key does not exist, the boolean
// returned is set to false.
func (m *memoryCache) Get(_ context.Context, key string) (bool, string, error) {
	iface := m.cache.Get(key)
	if iface == nil {
		return false, "", nil
	}

	s, ok := iface.Value().(string)
	if !ok {
		return false, "", ErrWrongType
	}

	return true, s, nil
}

// GetB retrieves a []byte value by key. If the key does not exist, the boolean
// returned is set to false.
func (m *memoryCache) GetB(_ context.Context, key string) (
	bool, []byte, error,
) {
	iface := m.cache.Get(key)
	if iface == nil {
		return false, nil, nil
	}

	b, ok := iface.Value().([]byte)
	if !ok {
		return false, nil, ErrWrongType
	}

	return true, b, nil
}

// GetI retrieves an int64 value by key. If the key does not exist, the boolean
// returned is set to false.
func (m *memoryCache) GetI(_ context.Context, key string) (bool, int64, error) {
	iface := m.cache.Get(key)
	if iface == nil {
		return false, 0, nil
	}

	i, ok := iface.Value().(int64)
	if !ok {
		return false, 0, ErrWrongType
	}

	return true, i, nil
}

// SetIfNotExist sets key to value if the key does not exist. If it is
// successful, it returns true.
func (m *memoryCache) SetIfNotExist(
	ctx context.Context, key string, value any,
) (bool, error) {
	return m.SetIfNotExistTTL(ctx, key, value, ttlcache.NoTTL)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If it is successful, it returns true.
func (m *memoryCache) SetIfNotExistTTL(
	_ context.Context, key string, value any, exp time.Duration,
) (bool, error) {
	if _, ok := m.cache.GetOrSet(key, value,
		ttlcache.WithTTL[string, any](exp)); ok {
		return false, nil
	}

	return true, nil
}

// Incr is not supported.
func (m *memoryCache) Incr(_ context.Context, _ string) (int64, error) {
	panic("unimplemented")
}

// Del removes the specified key. A key is ignored if it does not exist.
func (m *memoryCache) Del(_ context.Context, key string) error {
	m.cache.Delete(key)

	return nil
}

// Close closes the Cacher, releasing any open resources.
func (m *memoryCache) Close() error {
	m.cache.Stop()

	return nil
}
