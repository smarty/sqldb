package sqldb

import (
	"database/sql"
	"time"
)

type Wireup struct {
	inner            *sql.DB
	splitStatement   bool
	panicOnBindError bool
	parameterPrefix  string
	retrySleep       time.Duration
}

func Configure(pool *sql.DB) *Wireup {
	return &Wireup{inner: pool}
}

func (this *Wireup) WithPanicOnBindError() *Wireup {
	this.panicOnBindError = true
	return this
}

func (this *Wireup) WithMySQL() *Wireup {
	this.splitStatement = true
	this.parameterPrefix = "?"
	return this
}

func (this *Wireup) WithRetry(retrySleep time.Duration) *Wireup {
	this.retrySleep = retrySleep
	return this
}

func (this *Wireup) Build() ConnectionPool {
	var pool ConnectionPool = NewLibraryConnectionPoolAdapter(this.inner)

	if this.splitStatement {
		pool = NewSplitStatementConnectionPool(pool, this.parameterPrefix)
	}

	return pool
}

func (this *Wireup) BuildWithBinding() BindingConnectionPool {
	inner := this.Build()
	var pool BindingConnectionPool = NewBindingConnectionPoolAdapter(inner, this.panicOnBindError)

	if this.retrySleep > 0 {
		pool = NewRetryBindingConnectionPool(pool, this.retrySleep)
	}

	return pool
}
