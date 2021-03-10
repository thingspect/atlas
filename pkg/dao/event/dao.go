// Package event provides functions to create and query events in the database.
package event

import (
	"database/sql"
)

// DAO contains functions to create and query events in the database.
type DAO struct {
	pg *sql.DB
}

// NewDAO instantiates and returns a new DAO.
func NewDAO(pg *sql.DB) *DAO {
	return &DAO{
		pg: pg,
	}
}
