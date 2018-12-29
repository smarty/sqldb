package sqldb

import (
	"database/sql"
	"time"
)

type Wireup struct {
	inner             *sql.DB
	splitStatement    bool
	panicOnBindError  bool
	stackTraceOnError bool
	parameterPrefix   string
	retrySleep        time.Duration
}

func ConfigureConnectionPool(pool *sql.DB, options ...Option) ConnectionPool {
	wireup := &Wireup{inner: pool, parameterPrefix: "?"}
	wireup.configure(options...)
	return wireup.build()
}
func ConfigureBindingConnectionPool(pool *sql.DB, options ...Option) BindingConnectionPool {
	wireup := &Wireup{inner: pool, parameterPrefix: "?"}
	wireup.configure(options...)
	return wireup.buildWithBinding()
}
func (this *Wireup) configure(options ...Option) {
	for _, option := range options {
		option(this)
	}
}

type Option func(wireup *Wireup)

func WithPanicOnBindError() Option {
	return func(wireup *Wireup) { wireup.panicOnBindError = true }
}

func WithMySQL() Option {
	return func(wireup *Wireup) {
		wireup.splitStatement = true
		wireup.parameterPrefix = "?"
	}
}

func WithRetry(retrySleep time.Duration) Option {
	return func(wireup *Wireup) { wireup.retrySleep = retrySleep }
}

func WithStackTraceErrDiagnostics() Option {
	return func(wireup *Wireup) { wireup.stackTraceOnError = true }
}

func (this *Wireup) build() ConnectionPool {
	var pool ConnectionPool = NewLibraryConnectionPoolAdapter(this.inner)

	if this.splitStatement {
		pool = NewSplitStatementConnectionPool(pool, this.parameterPrefix)
	}

	if this.stackTraceOnError {
		pool = NewStackTraceConnectionPool(pool)
	}

	return pool
}

func (this *Wireup) buildWithBinding() BindingConnectionPool {
	inner := this.build()
	var pool BindingConnectionPool = NewBindingConnectionPoolAdapter(inner, this.panicOnBindError)

	if this.retrySleep > 0 {
		pool = NewRetryBindingConnectionPool(pool, this.retrySleep)
	}

	return pool
}
