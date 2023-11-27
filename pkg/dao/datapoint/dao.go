// Package datapoint provides functions to create and query data points in the
// database.
package datapoint

import (
	"database/sql"
)

// DAO contains functions to create and query data points in the database.
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
