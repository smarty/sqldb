package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type db struct {
	handle    Handle
	logger    Logger
	threshold int
	counts    map[string]int
	prepared  map[string]*sql.Stmt
}

func New(handle Handle, options ...option) DB {
	var config configuration
	Options.apply(options...)(&config)
	return &db{
		handle:    handle,
		logger:    config.logger,
		threshold: config.threshold,
		counts:    make(map[string]int),
		prepared:  make(map[string]*sql.Stmt),
	}
}

func (this *db) prepare(ctx context.Context, rawStatement string) (*sql.Stmt, error) {
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
	statement, err := this.handle.PrepareContext(ctx, rawStatement)
	if err != nil {
		return nil, err
	}
	this.prepared[rawStatement] = statement
	return statement, nil
}

func (this *db) Execute(ctx context.Context, script Script) (err error) {
	defer func() { err = NormalizeErr(err) }()
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
			_, err = this.handle.ExecContext(ctx, statement, params...)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
func (this *db) QueryRow(ctx context.Context, query Query) (err error) {
	defer func() { err = NormalizeErr(err) }()
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
		row = this.handle.QueryRowContext(ctx, statement, parameters...)
	}
	err = query.Scan(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
func (this *db) Query(ctx context.Context, query Query) (err error) {
	defer func() { err = NormalizeErr(err) }()
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
		rows, err = this.handle.QueryContext(ctx, statement, parameters...)
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
