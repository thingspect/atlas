// +build !integration

package dao

import (
	"database/sql"
	"fmt"
	"io"
	"testing"

	"github.com/jackc/pgconn"
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
		{&pgconn.PgError{Code: "22001"}, ErrInvalidFormat},
		{&pgconn.PgError{Code: "23514"}, ErrInvalidFormat},
		{&pgconn.PgError{Code: "22P02"}, ErrInvalidFormat},
		{&pgconn.PgError{Code: "23503"}, &pgconn.PgError{Code: "23503"}},
		{io.EOF, io.EOF},
	}

	for _, test := range tests {
		lTest := test

		t.Run(fmt.Sprintf("Can map %+v", lTest), func(t *testing.T) {
			t.Parallel()

			res := DBToSentinel(lTest.inp)
			t.Logf("res: %#v", res)
			require.Equal(t, lTest.res, res)
		})
	}
}
