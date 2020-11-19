// Package postgres provides a wrapper for setting up, configuring, and
// verifying a database/sql connection to PostgreSQL.
package postgres

import (
	"database/sql"
	"time"

	// pgx/stdlib imported for database/sql.
	_ "github.com/jackc/pgx/stdlib"
)

// New creates a new database/sql DB using the pgx driver.
func New(uri string) (*sql.DB, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	// For the specifics of sql.DB tuning, see:
	// https://www.alexedwards.net/blog/configuring-sqldb
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
