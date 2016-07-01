package sqldb

import "database/sql"

type Connection struct {
	actual *sql.DB
}

func NewConnection(actual *sql.DB) *Connection {
	return &Connection{actual: actual}
}

func (this *Connection) Ping() error {
	return this.actual.Ping()
}

func (this *Connection) BeginTransaction() (DriverTransaction, error) {
	if tx, err := this.actual.Begin(); err == nil {
		return NewTransaction(tx), nil
	} else {
		return nil, err
	}
}

func (this *Connection) Close() error {
	return this.actual.Close()
}

func (this *Connection) Execute(query string, parameters ...interface{}) (uint64, error) {
	if result, err := this.actual.Exec(query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *Connection) Select(query string, parameters ...interface{}) (SelectResult, error) {
	return this.actual.Query(query, parameters...)
}
