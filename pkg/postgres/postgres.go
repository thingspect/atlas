// Package postgres provides a wrapper for setting up, configuring, and
// verifying a database/sql connection to PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"time"

	// pgx/stdlib imported for database/sql.
	_ "github.com/jackc/pgx/v4/stdlib"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
