package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"runtime/debug"
	"strings"
)

type defaultHandle struct {
	pool      Pool
	logger    logger
	threshold int
	counts    map[string]int
	prepared  map[string]*sql.Stmt
}

func New(handle Pool, options ...option) Handle {
	var config configuration
	Options.apply(options...)(&config)
	return &defaultHandle{
		pool:      handle,
		logger:    config.logger,
		threshold: config.threshold,
		counts:    make(map[string]int),
		prepared:  make(map[string]*sql.Stmt),
	}
}

func (this *defaultHandle) prepare(ctx context.Context, rawStatement string) (*sql.Stmt, error) {
	if this.threshold < 0 {
		return nil, nil
	}
	if this.counts[rawStatement] < this.threshold {
		this.counts[rawStatement]++
		return nil, nil
	}
	statement, ok := this.prepared[rawStatement]
	if ok {
		return statement, nil
	}
	statement, err := this.pool.PrepareContext(ctx, rawStatement)
	if err != nil {
		return nil, err
	}
	this.prepared[rawStatement] = statement
	return statement, nil
}

func (this *defaultHandle) Execute(ctx context.Context, script Script) (err error) {
	defer func() { err = normalizeErr(err) }()
	statements := script.Statements()
	parameters := script.Parameters()
	placeholderCount := strings.Count(statements, "?")
	if placeholderCount != len(parameters) {
		return fmt.Errorf("%w: Expected: %d, received %d", ErrParameterCountMismatch, placeholderCount, len(parameters))
	}
	for statement, params := range interleaveParameters(statements, parameters...) {
		prepared, err := this.prepare(ctx, statement)
		if err != nil {
			return err
		}
		if prepared != nil {
			_, err = prepared.ExecContext(ctx, params...)
		} else {
			_, err = this.pool.ExecContext(ctx, statement, params...)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
func (this *defaultHandle) Populate(ctx context.Context, query Query) (err error) {
	defer func() { err = normalizeErr(err) }()
	statement := query.Statement()
	prepared, err := this.prepare(ctx, statement)
	if err != nil {
		return err
	}
	parameters := query.Parameters()
	var rows *sql.Rows
	if prepared != nil {
		rows, err = prepared.QueryContext(ctx, parameters...)
	} else {
		rows, err = this.pool.QueryContext(ctx, statement, parameters...)
	}
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		err = query.Scan(rows)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}
func (this *defaultHandle) PopulateRow(ctx context.Context, query Query) (err error) {
	defer func() { err = normalizeErr(err) }()
	statement := query.Statement()
	prepared, err := this.prepare(ctx, statement)
	if err != nil {
		return err
	}
	parameters := query.Parameters()
	var row *sql.Row
	if prepared != nil {
		row = prepared.QueryRowContext(ctx, parameters...)
	} else {
		row = this.pool.QueryRowContext(ctx, statement, parameters...)
	}
	err = query.Scan(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

// interleaveParameters splits the statements (on ';') and pairs each one with its corresponding parameters.
func interleaveParameters(statements string, parameters ...any) iter.Seq2[string, []any] {
	return func(yield func(string, []any) bool) {
		index := 0
		for statement := range strings.SplitSeq(statements, ";") {
			if len(strings.TrimSpace(statement)) == 0 {
				continue
			}
			statement += ";" // terminate the statement
			indexOffset := strings.Count(statement, "?")
			params := parameters[index : index+indexOffset]
			index += indexOffset
			if !yield(statement, params) {
				return
			}
		}
	}
}

// normalizeErr attaches a stack trace to non-nil errors and also normalizes errors that are
// semantically equal to context.Canceled. At present we are unaware whether this is still a
// commonly encountered scenario.
func normalizeErr(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "operation was canceled") {
		return fmt.Errorf("%w: %w", context.Canceled, err)
	}
	return fmt.Errorf("%w\nStack Trace:\n%s", err, string(debug.Stack()))
}
