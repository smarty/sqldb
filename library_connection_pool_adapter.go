package sqldb

import (
	"database/sql"
	"log"
)

type LibraryConnectionPoolAdapter struct {
	inner *sql.DB
}

func NewLibraryConnectionPoolAdapter(actual *sql.DB) *LibraryConnectionPoolAdapter {
	return &LibraryConnectionPoolAdapter{inner: actual}
}

func (this *LibraryConnectionPoolAdapter) Ping() error {
	err := this.inner.Ping()
	log.Println("[INFO] Pinging database. Error result:", err)
	return err
}
func (this *LibraryConnectionPoolAdapter) BeginTransaction() (Transaction, error) {
	if tx, err := this.inner.Begin(); err == nil {
		log.Println("[INFO] Beginning transaction")
		return NewLibraryTransactionAdapter(tx), nil
	} else {
		log.Println("[INFO] Unable to begin transaction")
		return nil, err
	}
}
func (this *LibraryConnectionPoolAdapter) Close() error {
	return this.inner.Close()
}

func (this *LibraryConnectionPoolAdapter) Execute(query string, parameters ...interface{}) (uint64, error) {
	log.Println("[INFO] Executing SQL statement:", query)
	if result, err := this.inner.Exec(query, parameters...); err != nil {
		log.Println("[INFO] SQL statement failed:", err)
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		log.Println("[INFO] SQL Statement succeeded:", count)
		return uint64(count), nil
	}
}

func (this *LibraryConnectionPoolAdapter) Select(query string, parameters ...interface{}) (SelectResult, error) {
	log.Println("[INFO] Executing SQL query:", query)
	return this.inner.Query(query, parameters...)
}
