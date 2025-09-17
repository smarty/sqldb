package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
)

// ScanOptionalRow receipts the *sql.Row, feeds the supplied arguments to the
// row.Scan() method and masks the result of sql.ErrNoRows.
func ScanOptionalRow(row *sql.Row, args ...any) error {
	err := row.Scan(args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return NormalizeErr(err)
}

// BindAll receives the *sql.Rows + error from the QueryContext method of either
// a *sql.DB, a *sql.Tx, or a *sql.Stmt, as well as a binder callback, to be called
// for each record, which gives the caller the opportunity to scan and aggregate values.
func BindAll(rows *sql.Rows, err error, scanner func(Scanner) error) error {
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		err = scanner(rows)
		if err != nil {
			return NormalizeErr(err)
		}
	}
	return nil
}

// ExecuteScript receives a Handle as well as one or more SQL statements (each ending in ';')
// with corresponding args. IOt executes each statement with the arguments corresponding to it.
func ExecuteScript(ctx context.Context, db Handle, statements string, args ...any) error {
	placeholderCount := strings.Count(statements, "?")
	if placeholderCount != len(args) {
		return fmt.Errorf("%w: Expected: %d, received %d", ErrParameterCountMismatch, placeholderCount, len(args))
	}
	for statement, params := range interleaveParameters(statements, args...) {
		_, err := db.ExecContext(ctx, statement, params...)
		if err != nil {
			return NormalizeErr(err)
		}
	}
	return nil
}

// NormalizeErr attaches a stack trace to non-nil errors and also normalizes errors that are
// semantically equal to context.Canceled. At present we are unaware whether this is still a
// commonly encountered scenario.
func NormalizeErr(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "operation was canceled") {
		return fmt.Errorf("%w: %w", context.Canceled, err)
	}
	return fmt.Errorf("%w\nStack Trace:\n%s", err, string(debug.Stack()))
}
