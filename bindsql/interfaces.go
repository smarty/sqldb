package bindsql

type (
	Connection interface {
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
		Select(Binder, string, ...interface{}) error
	}

	Binder func(Scanner) error

	Scanner interface {
		Scan(...interface{}) error
	}
)
