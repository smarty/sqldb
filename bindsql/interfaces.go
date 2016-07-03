package bindsql

type (
	BindingConnection interface {
		Ping() error
		BeginTransaction() (BindingTransaction, error)
		Close() error
		Executor
		Selector
	}

	BindingTransaction interface {
		Commit() error
		Rollback() error
		Executor
		Selector
	}

	Executor interface {
		Execute(string, ...interface{}) (uint64, error)
	}

	Selector interface {
		Select(Binder, string, ...interface{}) error
	}

	Binder func(Scanner) error

	Scanner interface {
		Scan(...interface{}) error
	}
)
