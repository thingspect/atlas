//go:build !integration

package dao

import (
	"database/sql"
	"fmt"
	"io"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestDBToSentinel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		inp error
		res error
	}{
		{nil, nil},
		{sql.ErrNoRows, ErrNotFound},
		{&pgconn.PgError{Code: "23505"}, ErrAlreadyExists},
		{
			&pgconn.PgError{Code: "22001", Message: "value too long"},
			fmt.Errorf("%w: value too long", ErrInvalidFormat),
		},
		{&pgconn.PgError{Code: "22001"}, ErrInvalidFormat},
		{&pgconn.PgError{Code: "22P02"}, ErrInvalidFormat},
		{
			&pgconn.PgError{Code: "22P02", Routine: "string_to_uuid"},
			fmt.Errorf("%w: UUID", ErrInvalidFormat),
		},
		{
			&pgconn.PgError{Code: "23514", ConstraintName: "constraint_name"},
			fmt.Errorf("%w: constraint_name", ErrInvalidFormat),
		},
		{
			&pgconn.PgError{Code: "23502", ColumnName: "column_name"},
			fmt.Errorf("%w: column_name", ErrInvalidFormat),
		},
		{
			&pgconn.PgError{Code: "23503", ConstraintName: "constraint_name"},
			fmt.Errorf("%w: constraint_name", ErrInvalidFormat),
		},
		{&pgconn.PgError{Code: "1"}, &pgconn.PgError{Code: "1"}},
		{io.EOF, io.EOF},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Can map %+v", test), func(t *testing.T) {
			t.Parallel()

			res := DBToSentinel(test.inp)
			t.Logf("res: %#v", res)
			require.Equal(t, test.res, res)
		})
	}
}
