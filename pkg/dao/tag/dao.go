// Package tag provides functions to query tags in the database.
package tag

import (
	"database/sql"
)

// DAO contains functions to query tags in the database.
type DAO struct {
	ro *sql.DB
}

// NewDAO instantiates and returns a new DAO.
func NewDAO(ro *sql.DB) *DAO {
	return &DAO{
		ro: ro,
	}
}
