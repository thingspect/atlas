// Package event provides functions to create and query events in the database.
package event

import (
	"database/sql"
)

// DAO contains functions to create and query events in the database.
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
