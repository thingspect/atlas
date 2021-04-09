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
	// SetIfNotExist sets key to value if the key does not exist. If it is
	// successful, it returns true.
	SetIfNotExist(ctx context.Context, key string, value interface{}) (bool,
		error)
	// SetIfNotExistTTL sets key to value, with expiration, if the key does not
	// exist. If it is successful, it returns true.
	SetIfNotExistTTL(ctx context.Context, key string, value interface{},
		exp time.Duration) (bool, error)
	// Close closes the Cacher, releasing any open resources.
	Close() error
}
