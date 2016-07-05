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
	}

	Selector interface {
		Select(string, ...interface{}) (SelectResult, error)
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
		BindingSelector
	}

	BindingSelector interface {
		BindSelect(Binder, string, ...interface{}) error
	}

	Binder func(Scanner) error
)
