// Package tag provides functions to query tags in the database.
package tag

import (
	"database/sql"
)

// DAO contains functions to query tags in the database.
type DAO struct {
	pg *sql.DB
}

// NewDAO instantiates and returns a new DAO.
func NewDAO(pg *sql.DB) *DAO {
	return &DAO{
		pg: pg,
	}
}
