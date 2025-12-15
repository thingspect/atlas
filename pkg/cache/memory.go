package cache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// memoryCache contains methods to create and query data in memory and
// implements the Cacher interface.
type memoryCache[V Cacheable] struct {
	cache *ttlcache.Cache[string, V]
}

// Verify memoryCache implements Cacher.
var _ Cacher[string] = &memoryCache[string]{}

// NewMemory builds a new Cacher and returns it.
func NewMemory[V Cacheable]() Cacher[V] {
	cache := ttlcache.New(
		ttlcache.WithDisableTouchOnHit[string, V](),
	)

	go cache.Start()

	return &memoryCache[V]{cache: cache}
}

// Set sets key to value.
func (m *memoryCache[V]) Set(ctx context.Context, key string, value V) error {
	return m.SetTTL(ctx, key, value, ttlcache.NoTTL)
}

// SetTTL sets key to value with expiration.
func (m *memoryCache[V]) SetTTL(_ context.Context, key string, value V,
	exp time.Duration,
) error {
	_ = m.cache.Set(key, value, exp)

	return nil
}

// Get retrieves a value by key. If the key does not exist, ErrNotFound is
// returned.
func (m *memoryCache[V]) Get(_ context.Context, key string) (V, error) {
	item := m.cache.Get(key)
	if item == nil {
		return *new(V), ErrNotFound
	}

	return item.Value(), nil
}

// SetIfNotExist sets key to value if the key does not exist. If the key already
// exists, ErrAlreadyExists is returned.
func (m *memoryCache[V]) SetIfNotExist(ctx context.Context, key string,
	value V,
) error {
	return m.SetIfNotExistTTL(ctx, key, value, ttlcache.NoTTL)
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If the key already exists, ErrAlreadyExists is returned.
func (m *memoryCache[V]) SetIfNotExistTTL(_ context.Context, key string,
	value V, exp time.Duration,
) error {
	if _, ok := m.cache.GetOrSet(key, value,
		ttlcache.WithTTL[string, V](exp)); ok {
		return ErrAlreadyExists
	}

	return nil
}

// Incr is not supported.
func (m *memoryCache[V]) Incr(_ context.Context, _ string) (int64, error) {
	panic("unimplemented")
}

// Del removes the specified key. A key is ignored if it does not exist.
func (m *memoryCache[V]) Del(_ context.Context, key string) error {
	m.cache.Delete(key)

	return nil
}

// Close closes the Cacher, releasing any open resources.
func (m *memoryCache[V]) Close() error {
	m.cache.Stop()

	return nil
}
