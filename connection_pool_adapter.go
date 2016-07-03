package sqldb

import "database/sql"

type ConnectionPoolAdapter struct {
	actual *sql.DB
}

func NewConnectionPoolAdapter(actual *sql.DB) *ConnectionPoolAdapter {
	return &ConnectionPoolAdapter{actual: actual}
}

func (this *ConnectionPoolAdapter) Ping() error {
	return this.actual.Ping()
}

func (this *ConnectionPoolAdapter) BeginTransaction() (Transaction, error) {
	if tx, err := this.actual.Begin(); err == nil {
		return NewTransactionAdapter(tx), nil
	} else {
		return nil, err
	}
}

func (this *ConnectionPoolAdapter) Close() error {
	return this.actual.Close()
}

func (this *ConnectionPoolAdapter) Execute(query string, parameters ...interface{}) (uint64, error) {
	if result, err := this.actual.Exec(query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *ConnectionPoolAdapter) Select(query string, parameters ...interface{}) (SelectResult, error) {
	return this.actual.Query(query, parameters...)
}
