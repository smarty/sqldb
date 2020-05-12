package sqldb

import "database/sql"

type LibraryConnectionPoolAdapter struct {
	inner *sql.DB
}

func NewLibraryConnectionPoolAdapter(actual *sql.DB) *LibraryConnectionPoolAdapter {
	return &LibraryConnectionPoolAdapter{inner: actual}
}

func (this *LibraryConnectionPoolAdapter) Ping() error {
	return this.inner.Ping()
}
func (this *LibraryConnectionPoolAdapter) BeginTransaction() (Transaction, error) {
	if tx, err := this.inner.Begin(); err == nil {
		return NewLibraryTransactionAdapter(tx), nil
	} else {
		return nil, err
	}
}
func (this *LibraryConnectionPoolAdapter) Close() error {
	return this.inner.Close()
}

func (this *LibraryConnectionPoolAdapter) Execute(query string, parameters ...interface{}) (uint64, error) {
	if result, err := this.inner.Exec(query, parameters...); err != nil {
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		return uint64(count), nil
	}
}

func (this *LibraryConnectionPoolAdapter) Select(query string, parameters ...interface{}) (SelectResult, error) {
	return this.inner.Query(query, parameters...)
}
