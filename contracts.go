package sqldb

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrParameterCountMismatch           = errors.New("the number of parameters supplied does not match the statement")
	ErrOptimisticConcurrencyCheckFailed = errors.New("optimistic concurrency check failed")
)

type (
	logger interface {
		Printf(string, ...any)
	}

	// Pool is a common subset of methods implemented by *sql.DB and *sql.Tx.
	// The name is a nod to that fact that a *sql.DB implements a Pool of connections.
	Pool interface {
		PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	}

	// Handle is a high level approach to common database operations, where each operation implements either
	// the Query or Script interface.
	Handle interface {
		Execute(context.Context, ...Script) error
		Populate(context.Context, ...Query) error
		PopulateRow(context.Context, ...Query) error
	}

	// Script represents SQL statements that aren't expected to provide rows as a result.
	Script interface {
		// Statements returns a string containing 1 or more SQL statements, separated by `;`.
		// This means that the ';' character should NOT be used within any of the statements.
		Statements() string

		// Parameters returns a slice of the parameters to be interleaved across all SQL
		// returned from Statements().
		Parameters() []any
	}

	// RowsAffected provides an (optional) hook for a type implementing Script to receive
	// the number of rows affected by executing a statement provided by a Script. It is
	// called for each statement that doesn't result in an error.
	RowsAffected interface {
		RowsAffected(uint64)
	}

	// OptimisticConcurrencyCheck provides an (optional) hook for a type implementing
	// Script to verify whether the total count of all rows affected matches the returned
	// value. If not, an error wrapped with ErrOptimisticConcurrencyCheckFailed will be
	// returned by Handle.Execute().
	OptimisticConcurrencyCheck interface {
		ExpectedRowsAffected() uint64
	}

	// Query represents a SQL statement that is expected to provide rows as a result.
	// Rows are provided to the Scan method.
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
