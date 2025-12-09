// Package cache provides functions to create and query data in a cache.
package cache

//go:generate mockgen -source cacher.go -destination mock_cacher.go -package cache

import (
	"context"
	"time"

	"github.com/thingspect/atlas/pkg/consterr"
)

// Sentinel errors for Cacher.
const (
	ErrAlreadyExists consterr.Error = "object already exists"
	ErrNotFound      consterr.Error = "object not found"
)

// Cacher defines the methods provided by a Cache.
type Cacher[V any] interface {
	// Set sets key to value.
	Set(ctx context.Context, key string, value V) error
	// SetTTL sets key to value with expiration.
	SetTTL(ctx context.Context, key string, value V, exp time.Duration) error
	// Get retrieves a value by key. If the key does not exist, ErrNotFound is
	// returned.
	Get(ctx context.Context, key string) (V, error)
	// SetIfNotExist sets key to value if the key does not exist. If the key
	// already exists, ErrAlreadyExists is returned.
	SetIfNotExist(ctx context.Context, key string, value V) error
	// SetIfNotExistTTL sets key to value, with expiration, if the key does not
	// exist. If the key already exists, ErrAlreadyExists is returned.
	SetIfNotExistTTL(ctx context.Context, key string, value V,
		exp time.Duration) error
	// Incr increments an int64 value at key by one. If the key does not exist,
	// the value is set to 1. The incremented value is returned.
	//
	// Incr is not supported by memory caches.
	Incr(ctx context.Context, key string) (int64, error)
	// Del removes the specified key. A key is ignored if it does not exist.
	Del(ctx context.Context, key string) error
	// Close closes the Cacher, releasing any open resources.
	Close() error
}
