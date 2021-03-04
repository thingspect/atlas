package dao

import (
	"context"
	"database/sql"
	"time"

	// pgx/stdlib imported for database/sql.
	_ "github.com/jackc/pgx/v4/stdlib"
)

// NewPgDB builds, configures, and verifies a new database/sql DB using the pgx
// driver.
func NewPgDB(uri string) (*sql.DB, error) {
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
