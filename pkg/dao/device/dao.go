// Package device provides functions to query and modify devices in the
// database.
package device

import (
	"database/sql"
	"time"

	"github.com/thingspect/atlas/pkg/cache"
)

// DAO contains functions to query and modify devices in the database.
type DAO struct {
	rw    *sql.DB
	ro    *sql.DB
	cache cache.Cacher
	exp   time.Duration
}

// NewDAO instantiates and returns a new DAO with device read caching. Cache
// can be set to nil to disable caching.
func NewDAO(rw *sql.DB, ro *sql.DB, cache cache.Cacher, exp time.Duration) *DAO {
	return &DAO{
		rw:    rw,
		ro:    ro,
		cache: cache,
		exp:   exp,
	}
}
