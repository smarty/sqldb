package sqldb

import (
	"context"
	"database/sql"
	"errors"
)

var ErrParameterCountMismatch = errors.New("the number of parameters supplied does not match the statement")

type (
	logger interface {
		Printf(string, ...any)
	}

	// Pool is a common subset of methods implemented by *sql.DB and *sql.Tx.
	// The name is a nod to that fact that a *sql.DB implements a Pool of connections.
	// TODO: perhaps un-export?
	Pool interface {
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	}

	// Handle is a high level approach to common database operations, where each operation implements either
	// the Query or Script interface.
	Handle interface {
		Execute(ctx context.Context, script Script) error
		Populate(ctx context.Context, query Query) error
		PopulateRow(ctx context.Context, query Query) error
	}

	// Script represents SQL statements that aren't expected to provide rows as a result.
	// TODO: perhaps un-export?
	Script interface {
		// Statements returns a string containing 1 or more SQL statements, separated by `;`.
		// This means that the ';' character should NOT be used within any of the statements.
		Statements() string

		// Parameters returns a slice of the parameters to be interleaved across all sql returns by Statements().
		Parameters() []any
	}

	// RowsAffected provides an (optional) hook for a type implemented Script to receive
	// the number of rows affected by executing a statement provided by a Script. It is
	// called for each statement that doesn't result in an error.
	RowsAffected interface {
		RowsAffected(uint64)
	}

	// Query represents a SQL statement that is expected to provide rows as a result.
	// Rows are fed to the Scan method.
	// TODO: perhaps un-export?
	Query interface {
		// Statement returns a string containing a single SQL query.
		Statement() string

		// Parameters returns a slice of the parameters to be used in the query.
		Parameters() []any

		// Scan specifies a callback to be provided the *sql.Rows for each retrieved record.
		Scan(Scanner) error
	}

	// Scanner is implemented by *sql.Row and *sql.Rows.
	Scanner interface {
		Scan(...any) error
	}
)
