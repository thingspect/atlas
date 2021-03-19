// Package cache provides functions to create and query data in a cache.
package cache

//go:generate mockgen -source cacher.go -destination mock_cacher.go -package cache

import (
	"context"
	"time"
)

// Cacher defines the methods provided by a Cache.
type Cacher interface {
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
