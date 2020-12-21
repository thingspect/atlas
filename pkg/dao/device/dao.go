// Package device provides functions to query and modify devices in the
// database.
package device

import (
	"database/sql"
)

// DAO contains functions to query and modify devices in the database.
type DAO struct {
	pg *sql.DB
}

// NewDAO instantiates and returns a new DAO.
func NewDAO(pg *sql.DB) *DAO {
	return &DAO{
		pg: pg,
	}
}
