// Package cache provides functions to create and query data in a cache.
package cache

//go:generate mockgen -source cacher.go -destination mock_cacher.go -package cache

import (
	"context"
	"time"
)

// Cacher defines the methods provided by a Cache.
type Cacher interface {
	// Set sets key to value.
	Set(ctx context.Context, key string, value interface{}) error
	// SetTTL sets key to value with expiration.
	SetTTL(ctx context.Context, key string, value interface{},
		exp time.Duration) error
	// Get retrieves a string value by key. If the key does not exist, the
	// boolean returned is set to false.
	Get(ctx context.Context, key string) (bool, string, error)
	// Get retrieves a []byte value by key. If the key does not exist, the
	// boolean returned is set to false.
	GetB(ctx context.Context, key string) (bool, []byte, error)
	// GetI retrieves an int64 value by key. If the key does not exist, the
	// boolean returned is set to false.
	GetI(ctx context.Context, key string) (bool, int64, error)
	// SetIfNotExist sets key to value if the key does not exist. If it is
	// successful, it returns true.
	SetIfNotExist(ctx context.Context, key string, value interface{}) (bool,
		error)
	// SetIfNotExistTTL sets key to value, with expiration, if the key does not
	// exist. If it is successful, it returns true.
	SetIfNotExistTTL(ctx context.Context, key string, value interface{},
		exp time.Duration) (bool, error)
	// Incr increments an int64 value at key by one. If the key does not exist,
	// the value is set to 1. The incremented value is returned.
	Incr(ctx context.Context, key string) (int64, error)
	// Del removes the specified key. A key is ignored if it does not exist.
	Del(ctx context.Context, key string) error
	// Close closes the Cacher, releasing any open resources.
	Close() error
}
