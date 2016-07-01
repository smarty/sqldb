package sqldb

type (
	DriverConnection interface {
		Ping() error
		BeginTransaction() (DriverTransaction, error)
		Close() error
		Executor
		Selector
	}

	DriverTransaction interface {
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
