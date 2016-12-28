package sqldb

import (
	"fmt"
	"runtime/debug"
)

type StackTraceConnectionPool struct {
	inner ConnectionPool
	stack *StackTrace
}

func NewStackTraceConnectionPool(inner ConnectionPool) *StackTraceConnectionPool {
	return &StackTraceConnectionPool{inner: inner}
}

func (this *StackTraceConnectionPool) Ping() error {
	return this.wrap(this.inner.Ping())
}

func (this *StackTraceConnectionPool) BeginTransaction() (Transaction, error) {
	tx, err := this.inner.BeginTransaction()
	return tx, this.wrap(err)
}

func (this *StackTraceConnectionPool) Close() error {
	return this.wrap(this.inner.Close())
}

func (this *StackTraceConnectionPool) Execute(query string, parameters ...interface{}) (uint64, error) {
	rows, err := this.inner.Execute(query, parameters...)
	return rows, this.wrap(err)
}

func (this *StackTraceConnectionPool) Select(query string, parameters ...interface{}) (SelectResult, error) {
	result, err := this.inner.Select(query, parameters...)
	return result, this.wrap(err)
}

func (this *StackTraceConnectionPool) wrap(err error) error {
	if err != nil {
		err = fmt.Errorf("%s\nStack Trace:\n%s", err.Error(), this.stack.StackTrace())
	}
	return err
}

/**************************************************************************/

// StackTrace, like github.com/smartystreets/clock.Clock performs in production mode
// when used as a nil pointer struct field. When non-nil, it returns a preset value.
// This is useful during testing when asserting on simple, deterministic values is helpful.
type StackTrace struct {
	trace string
}

func ContrivedStackTrace(trace string) *StackTrace {
	return &StackTrace{trace: trace}
}

func (this *StackTrace) StackTrace() string {
	if this == nil {
		return string(debug.Stack())
	}
	return this.trace
}
