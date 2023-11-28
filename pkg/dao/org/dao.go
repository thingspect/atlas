// Package org provides functions to query and modify organizations in the
// database.
package org

import (
	"database/sql"
)

// DAO contains functions to query and modify organizations in the database.
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
