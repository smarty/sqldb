package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"
	"iter"
	"runtime/debug"
	"strings"
	"sync"
)

type defaultHandle struct {
	pool      Pool
	logger    logger
	threshold int
	lock      *sync.Mutex
	counts    map[uint64]int       // map[sql-statement-checksum]count
	prepared  map[uint64]*sql.Stmt // map[sql-statement-checksum]stmt
}

func New(handle Pool, options ...option) Handle {
	var config configuration
	Options.apply(options...)(&config)
	return &defaultHandle{
		pool:      handle,
		logger:    config.logger,
		threshold: config.threshold,
		lock:      new(sync.Mutex),
		counts:    make(map[uint64]int),
		prepared:  make(map[uint64]*sql.Stmt),
	}
}

func (this *defaultHandle) prepare(ctx context.Context, rawStatement string) (*sql.Stmt, error) {
	if this.threshold < 0 {
		return nil, nil
	}
	if len(this.counts) > 1024*64 { // Put some kind of cap on how many statements we will track.
		return nil, nil
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	checksum := checksum([]byte(rawStatement))
	if this.counts[checksum] < this.threshold {
		this.counts[checksum]++
		return nil, nil
	}
	statement, ok := this.prepared[checksum]
	if ok {
		return statement, nil
	}
	statement, err := this.pool.PrepareContext(ctx, rawStatement)
	if err != nil {
		return nil, err
	}
	this.prepared[checksum] = statement
	return statement, nil
}
func checksum(x []byte) (hash uint64) {
	h := fnv.New64a()
	_, _ = h.Write(x)
	return h.Sum64()
}

func (this *defaultHandle) Execute(ctx context.Context, scripts ...Script) (err error) {
	defer func() { err = normalizeErr(err) }()
	for _, script := range scripts {
		statements := script.Statements()
		parameters := script.Parameters()
		placeholderCount := strings.Count(statements, "?")
		if placeholderCount != len(parameters) {
			return fmt.Errorf("%w: Expected: %d, received %d", ErrParameterCountMismatch, placeholderCount, len(parameters))
		}
		var actualRowsAffectedCount uint64
		for statement, params := range interleaveParameters(statements, parameters...) {
			prepared, err := this.prepare(ctx, statement)
			if err != nil {
				return err
			}
			var result sql.Result
			if prepared != nil {
				result, err = prepared.ExecContext(ctx, params...)
			} else {
				result, err = this.pool.ExecContext(ctx, statement, params...)
			}
			if err != nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			rowCount := uint64(max(0, affected))
			actualRowsAffectedCount += rowCount
			if rows, ok := script.(RowsAffected); ok {
				rows.RowsAffected(rowCount)
			}
		}
		if check, ok := script.(OptimisticConcurrencyCheck); ok {
			expectedRowsAffectedCount := check.ExpectedRowsAffected()
			if actualRowsAffectedCount != expectedRowsAffectedCount {
				return fmt.Errorf("%w: expected rows affected: %d (actual: %d)",
					ErrOptimisticConcurrencyCheckFailed, expectedRowsAffectedCount, actualRowsAffectedCount,
				)
			}
		}
	}
	return nil
}
func (this *defaultHandle) Populate(ctx context.Context, queries ...Query) (err error) {
	defer func() { err = normalizeErr(err) }()
	for _, query := range queries {
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
		for rows.Next() {
			err = query.Scan(rows)
			if err != nil {
				_ = rows.Close()
				return err
			}
		}
		_ = rows.Close()
	}
	return nil
}
func (this *defaultHandle) PopulateRow(ctx context.Context, queries ...Query) (err error) {
	defer func() { err = normalizeErr(err) }()
	for _, query := range queries {
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
		if err == nil {
			continue
		}
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		return err
	}
	return nil
}

// interleaveParameters splits the statements (on ';') and yields each with its corresponding parameters.
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
