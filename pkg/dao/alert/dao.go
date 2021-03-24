// Package alert provides functions to create and query alerts in the database.
package alert

import (
	"database/sql"
)

// DAO contains functions to create and query alerts in the database.
type DAO struct {
	pg *sql.DB
}

// NewDAO instantiates and returns a new DAO.
func NewDAO(pg *sql.DB) *DAO {
	return &DAO{
		pg: pg,
	}
}
