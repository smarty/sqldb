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
