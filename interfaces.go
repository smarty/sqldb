package sqldb

import "context"

type (
	ConnectionPool interface {
		Ping(context.Context) error
		BeginTransaction(context.Context) (Transaction, error)
		Close() error
		Executor
		Selector
	}

	Transaction interface {
		Commit() error
		Rollback() error
		Executor
		Selector
	}

	Executor interface {
		Execute(context.Context, string, ...interface{}) (uint64, error)
	}

	Selector interface {
		Select(context.Context, string, ...interface{}) (SelectResult, error)
	}

	SelectExecutor interface {
		Selector
		Executor
	}

	SelectResult interface {
		Next() bool
		Err() error
		Close() error
		Scanner
	}

	Scanner interface {
		Scan(...interface{}) error
	}
)

type (
	BindingConnectionPool interface {
		Ping(context.Context) error
		BeginTransaction(context.Context) (BindingTransaction, error)
		Close() error
		Executor
		BindingSelector
	}

	BindingTransaction interface {
		Commit() error
		Rollback() error
		Executor
		BindingSelector
	}

	BindingSelector interface {
		BindSelect(context.Context, Binder, string, ...interface{}) error
	}

	Binder func(Scanner) error

	BindingSelectExecutor interface {
		BindingSelector
		Executor
	}
)
