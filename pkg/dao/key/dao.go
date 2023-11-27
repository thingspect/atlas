// Package key provides functions to create and query API keys in the database.
package key

import (
	"database/sql"
)

// DAO contains functions to create and query API keys in the database.
type DAO struct {
	rw *sql.DB
	ro *sql.DB
}

// NewDAO instantiates and returns a new DAO.
func NewDAO(rw *sql.DB, ro *sql.DB) *DAO {
	return &DAO{
		rw: rw,
		ro: ro,
	}
}
