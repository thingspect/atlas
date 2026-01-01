package cache

import (
	"context"
	"time"

	"github.com/maypok86/otter/v2"
)

// heapCache contains methods to create and query data in memory and
// implements the Cacher interface.
type heapCache[V Cacheable] struct {
	cache *otter.Cache[string, V]
}

// Verify heapCache implements Cacher.
var _ Cacher[string] = &heapCache[string]{}

// NewHeap builds a new Cacher and returns it.
func NewHeap[V Cacheable]() Cacher[V] {
	cache := otter.Must(&otter.Options[string, V]{
		ExpiryCalculator: otter.ExpiryCreating[string, V](1 << 62),
	})

	return &heapCache[V]{cache: cache}
}

// Set sets key to value.
func (h *heapCache[V]) Set(_ context.Context, key string, value V) error {
	_, _ = h.cache.Set(key, value)

	return nil
}

// SetTTL sets key to value with expiration.
func (h *heapCache[V]) SetTTL(_ context.Context, key string, value V, exp time.Duration) error {
	_, _ = h.cache.Set(key, value)
	h.cache.SetExpiresAfter(key, exp)

	return nil
}

// Get retrieves a value by key. If the key does not exist, ErrNotFound is
// returned.
func (h *heapCache[V]) Get(_ context.Context, key string) (V, error) {
	item, ok := h.cache.GetIfPresent(key)
	if !ok {
		return *new(V), ErrNotFound
	}

	return item, nil
}

// SetIfNotExist sets key to value if the key does not exist. If the key already
// exists, ErrAlreadyExists is returned.
func (h *heapCache[V]) SetIfNotExist(_ context.Context, key string, value V) error {
	_, ok := h.cache.SetIfAbsent(key, value)
	if !ok {
		return ErrAlreadyExists
	}

	return nil
}

// SetIfNotExistTTL sets key to value, with expiration, if the key does not
// exist. If the key already exists, ErrAlreadyExists is returned.
func (h *heapCache[V]) SetIfNotExistTTL(_ context.Context, key string, value V, exp time.Duration) error {
	_, ok := h.cache.SetIfAbsent(key, value)
	if !ok {
		return ErrAlreadyExists
	}

	h.cache.SetExpiresAfter(key, exp)

	return nil
}

// Incr is not supported.
func (h *heapCache[V]) Incr(_ context.Context, _ string) (int64, error) {
	panic("unimplemented")
}

// Del removes the specified key. A key is ignored if it does not exist.
func (h *heapCache[V]) Del(_ context.Context, key string) error {
	_, _ = h.cache.Invalidate(key)

	return nil
}

// Close closes the Cacher, releasing any open resources.
func (h *heapCache[V]) Close() error {
	h.cache.InvalidateAll()
	h.cache.CleanUp()
	_ = h.cache.StopAllGoroutines()

	return nil
}
