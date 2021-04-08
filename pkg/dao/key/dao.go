// Package key provides functions to create and query API keys in the database.
package key

import (
	"database/sql"
)

// DAO contains functions to create and query API keys in the database.
type DAO struct {
	pg *sql.DB
}

// NewDAO instantiates and returns a new DAO.
func NewDAO(pg *sql.DB) *DAO {
	return &DAO{
		pg: pg,
	}
}
