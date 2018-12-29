package sqldb

import (
	"database/sql"
	"log"
)

type LibraryTransactionAdapter struct {
	inner *sql.Tx
}

func NewLibraryTransactionAdapter(actual *sql.Tx) *LibraryTransactionAdapter {
	return &LibraryTransactionAdapter{inner: actual}
}

func (this *LibraryTransactionAdapter) Commit() error {
	log.Println("[INFO] Committing transaction")
	return this.inner.Commit()
}
func (this *LibraryTransactionAdapter) Rollback() error {
	log.Println("[INFO] Rolling back transaction")
	return this.inner.Rollback()
}

func (this *LibraryTransactionAdapter) Execute(query string, parameters ...interface{}) (uint64, error) {
	log.Println("[INFO] Executing transactions SQL statement:", query)
	if result, err := this.inner.Exec(query, parameters...); err != nil {
		log.Println("[INFO] Transactional SQL statement failed:", err)
		return 0, err
	} else {
		count, _ := result.RowsAffected()
		log.Println("[INFO] Transactional SQL Statement succeeded:", count)
		return uint64(count), nil
	}
}

func (this *LibraryTransactionAdapter) Select(query string, parameters ...interface{}) (SelectResult, error) {
	log.Println("[INFO] Executing transactional SQL query:", query)
	return this.inner.Query(query, parameters...)
}
