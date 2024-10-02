package sqldb

import (
	"database/sql"
	"time"
)

type configuration struct {
	txOptions                    *sql.TxOptions
	splitStatement               bool
	panicOnBindError             bool
	normalizeContextCancellation bool
	stackTraceOnError            bool
	parameterPrefix              string
	retrySleep                   time.Duration
}

func NewPool(handle *sql.DB, options ...option) ConnectionPool {
	var config configuration
	Options.apply(options...)(&config)
	return newPool(handle, config)
}
func NewBindingPool(handle *sql.DB, options ...option) BindingConnectionPool {
	var config configuration
	Options.apply(options...)(&config)
	return newBindingPool(handle, config)
}
func newPool(handle *sql.DB, config configuration) ConnectionPool {
	var pool ConnectionPool = NewLibraryConnectionPoolAdapter(handle, config.txOptions)

	if config.splitStatement {
		pool = NewSplitStatementConnectionPool(pool, config.parameterPrefix)
	}

	if config.normalizeContextCancellation {
		pool = NewNormalizeContextCancellationConnectionPool(pool)
	}

	if config.stackTraceOnError {
		pool = NewStackTraceConnectionPool(pool)
	}

	return pool
}
func newBindingPool(handle *sql.DB, config configuration) BindingConnectionPool {
	inner := newPool(handle, config)
	var pool BindingConnectionPool = NewBindingConnectionPoolAdapter(inner, config.panicOnBindError)

	if config.retrySleep > 0 {
		pool = NewRetryBindingConnectionPool(pool, config.retrySleep)
	}

	return pool
}

var Options singleton

type singleton struct{}
type option func(*configuration)

func (singleton) TxOptions(value *sql.TxOptions) option {
	return func(this *configuration) { this.txOptions = value }
}
func (singleton) PanicOnBindError(value bool) option {
	return func(this *configuration) { this.panicOnBindError = value }
}
func (singleton) NormalizeContextCancellation(value bool) option {
	return func(this *configuration) { this.normalizeContextCancellation = value }
}
func (singleton) MySQL() option {
	return func(this *configuration) { this.splitStatement = true; this.parameterPrefix = "?" }
}
func (singleton) ParameterPrefix(value string) option {
	return func(this *configuration) { this.parameterPrefix = value }
}
func (singleton) SplitStatement(value bool) option {
	return func(this *configuration) { this.splitStatement = value }
}
func (singleton) RetrySleep(value time.Duration) option {
	return func(this *configuration) { this.retrySleep = value }
}
func (singleton) StackTraceErrDiagnostics(value bool) option {
	return func(this *configuration) { this.stackTraceOnError = value }
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, option := range Options.defaults(options...) {
			option(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	var defaultTxOptions = &sql.TxOptions{Isolation: sql.LevelReadCommitted}
	const defaultStackTraceErrDiagnostics = true
	const defaultPanicOnBindError = true
	const defaultNormalizeContextCancellation = true
	const defaultSplitStatement = true
	const defaultParameterPrefix = "?"
	const defaultRetrySleep = 0

	return append([]option{
		Options.TxOptions(defaultTxOptions),
		Options.PanicOnBindError(defaultPanicOnBindError),
		Options.NormalizeContextCancellation(defaultNormalizeContextCancellation),
		Options.StackTraceErrDiagnostics(defaultStackTraceErrDiagnostics),
		Options.ParameterPrefix(defaultParameterPrefix),
		Options.SplitStatement(defaultSplitStatement),
		Options.RetrySleep(defaultRetrySleep),
	}, options...)
}
