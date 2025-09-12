package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

var ErrArgumentCountMismatch = errors.New("the number of arguments supplied does not match the statement")

type Binder func(Scanner) error

type Scanner interface {
	Scan(...any) error
}

// DBTx is either a *sql.DB or a *sql.Tx
type DBTx interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// BindAll receives the *sql.Rows + error from the QueryContext method of either
// a *sql.DB, a *sql.Tx, or a *sql.Stmt, as well as a binder callback, to be called
// for each record, which gives the caller the opportunity to scan and aggregate values.
func BindAll(rows *sql.Rows, err error, binder Binder) error {
	if err != nil {
		return normalize(err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		err = binder(rows)
		if err != nil {
			return normalize(err)
		}
	}
	return nil
}

// ExecuteStatements receives a *sql.DB or *sql.Tx as well as one or more SQL statements (separated by ';')
// and executes each one with the arguments corresponding to that statement.
func ExecuteStatements(ctx context.Context, db DBTx, statements string, args ...any) (uint64, error) {
	placeholderCount := strings.Count(statements, "?")
	if placeholderCount != len(args) {
		return 0, fmt.Errorf("%w: Expected: %d, received %d", ErrArgumentCountMismatch, placeholderCount, len(args))
	}
	var count uint64
	index := 0
	for statement := range strings.SplitSeq(statements, ";") {
		if len(strings.TrimSpace(statement)) == 0 {
			continue
		}
		statement += ";" // terminate the statement
		indexOffset := strings.Count(statement, "?")
		result, err := db.ExecContext(ctx, statement, args[index:index+indexOffset]...)
		rows, err := RowsAffected(result, err)
		if err != nil {
			return 0, err // already normalized
		}
		count += rows
		index += indexOffset
	}
	return count, nil
}

// RowsAffected returns the rows affected from a sql.Result. This is generally only needed
// by external callers when dealing with the result of a prepared statement.
func RowsAffected(result sql.Result, err error) (uint64, error) {
	if err != nil {
		return 0, normalize(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, normalize(err)
	}
	return uint64(rows), nil
}
func normalize(err error) error {
	if err == nil {
		return nil
	}
	err = fmt.Errorf("%w\nStack Trace:\n%s", err, string(debug.Stack()))
	if strings.Contains(err.Error(), "operation was canceled") {
		return fmt.Errorf("%w: %w", context.Canceled, err)
	}
	return err
}
