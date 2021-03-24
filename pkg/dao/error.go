// Package dao provides sentinel errors and helper functions for use by data
// access object packages.
package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/consterr"
)

// Sentinel errors for DAO packages.
const (
	ErrAlreadyExists consterr.Error = "object already exists"
	ErrInvalidFormat consterr.Error = "invalid format"
	ErrNotFound      consterr.Error = "object not found"
)

// DBToSentinel maps database/sql or driver errors to sentinel errors. This
// function should only be used from within DAO packages.
func DBToSentinel(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		// unique_violation
		case "23505":
			return ErrAlreadyExists
		// string_data_right_truncation
		case "22001":
			if strings.Contains(pgErr.Message, "value too long") {
				return fmt.Errorf("%w: value too long", ErrInvalidFormat)
			}

			return ErrInvalidFormat
		// invalid_text_representation
		case "22P02":
			if pgErr.File == "uuid.c" {
				return fmt.Errorf("%w: UUID", ErrInvalidFormat)
			}

			return ErrInvalidFormat
		// check_violation
		case "23514":
			return fmt.Errorf("%w: %s", ErrInvalidFormat, pgErr.ConstraintName)
		// not_null_violation
		case "23502":
			return fmt.Errorf("%w: %s", ErrInvalidFormat, pgErr.ColumnName)
		// foreign_key_violation
		case "23503":
			return fmt.Errorf("%w: %s", ErrInvalidFormat, pgErr.ConstraintName)
		default:
			alog.Errorf("DBToSentinel unmatched PgError: %#v", pgErr)

			return err
		}
	}

	alog.Errorf("DBToSentinel unmatched error: %#v", err)

	return err
}
