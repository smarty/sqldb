package sqldb

import (
	"database/sql"
	"time"
)

type Wireup struct {
	inner             *sql.DB
	txOptions         *sql.TxOptions
	splitStatement    bool
	panicOnBindError  bool
	stackTraceOnError bool
	parameterPrefix   string
	retrySleep        time.Duration
}

func (this *Wireup) configure(options ...Option) {
	for _, option := range options {
		option(this)
	}
}

func ConfigureConnectionPool(pool *sql.DB, options ...Option) ConnectionPool {
	this := &Wireup{inner: pool}
	this.configure(options...)
	return this.build()
}
func (this *Wireup) build() ConnectionPool {
	if this.txOptions == nil {
		this.txOptions = &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false}
	}

	var pool ConnectionPool = NewLibraryConnectionPoolAdapter(this.inner, this.txOptions)

	if this.splitStatement {
		pool = NewSplitStatementConnectionPool(pool, this.parameterPrefix)
	}

	if this.stackTraceOnError {
		pool = NewStackTraceConnectionPool(pool)
	}

	return pool
}

func ConfigureBindingConnectionPool(pool *sql.DB, options ...Option) BindingConnectionPool {
	this := &Wireup{inner: pool}
	this.configure(options...)
	return this.buildWithBinding()
}
func (this *Wireup) buildWithBinding() BindingConnectionPool {
	inner := this.build()
	var pool BindingConnectionPool = NewBindingConnectionPoolAdapter(inner, this.panicOnBindError)

	if this.retrySleep > 0 {
		pool = NewRetryBindingConnectionPool(pool, this.retrySleep)
	}

	return pool
}

type Option func(*Wireup)

func WithTxOptions(value *sql.TxOptions) Option {
	return func(this *Wireup) { this.txOptions = value }
}
func WithPanicOnBindError() Option {
	return func(this *Wireup) { this.panicOnBindError = true }
}
func WithMySQL() Option {
	return func(this *Wireup) { this.splitStatement = true; this.parameterPrefix = "?" }
}
func WithParameterPrefix(value string) Option {
	return func(this *Wireup) { this.parameterPrefix = value }
}
func WithRetry(value time.Duration) Option {
	return func(this *Wireup) { this.retrySleep = value }
}
func WithStackTraceErrDiagnostics() Option {
	return func(this *Wireup) { this.stackTraceOnError = true }
}
