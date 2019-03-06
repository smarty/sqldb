package sqldb

type (
	ConnectionPool interface {
		Ping() error
		BeginTransaction() (Transaction, error)
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
		Execute(string, ...interface{}) (uint64, error)
		ExecuteIdentity(string, ...interface{}) (uint64, uint64, error)
	}

	Selector interface {
		Select(string, ...interface{}) (SelectResult, error)
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
		Ping() error
		BeginTransaction() (BindingTransaction, error)
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
		BindSelect(Binder, string, ...interface{}) error
	}

	Binder func(Scanner) error

	BindingSelectExecutor interface {
		BindingSelector
		Executor
	}
)
