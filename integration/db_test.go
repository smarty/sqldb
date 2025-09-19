package integration

import (
	"database/sql"
	"log"
	"testing"

	"github.com/smarty/gunit/v2"
	"github.com/smarty/gunit/v2/better"
	"github.com/smarty/gunit/v2/should"
	"github.com/smarty/sqldb/v3"

	_ "github.com/mattn/go-sqlite3"
)

func TestFixture(t *testing.T) {
	gunit.Run(new(Fixture), t)
}

type Fixture struct {
	*gunit.Fixture
	db *sql.DB
	tx *sql.Tx
	DB sqldb.Handle
}

func (this *Fixture) Setup() {
	var err error
	this.db, err = sql.Open("sqlite3", ":memory:")
	this.So(err, better.BeNil)

	this.tx, err = this.db.BeginTx(this.Context(), nil)
	this.So(err, better.BeNil)

	this.DB = sqldb.New(this.db,
		sqldb.Options.Logger(log.New(this.Output(), this.Name()+": ", 0)),
		sqldb.Options.PreparationThreshold(5),
	)
	ddl := &DDL{}
	err = this.DB.Execute(this.Context(), ddl)
	this.So(err, better.BeNil)
	this.So(ddl.totalRows, should.Equal, 4)
}
func (this *Fixture) Teardown() {
	this.So(this.tx.Rollback(), better.BeNil)
	this.So(this.db.Close(), better.BeNil)
}

func (this *Fixture) TestQuery() {
	for range 10 { // should transition to prepared statements
		query := &SelectAll{Result: make(map[int]string)}
		err := this.DB.Populate(this.Context(), query)
		this.So(err, better.BeNil)
		this.So(query.Result, should.Equal, map[int]string{
			1: "a",
			2: "b",
			3: "c",
			4: "d",
		})
	}
}
func (this *Fixture) TestQueryRow() {
	query := &SelectRow{id: 1}
	err := this.DB.PopulateRow(this.Context(), query)
	this.So(err, better.BeNil)
	this.So(query.value, should.Equal, "a")
}
func (this *Fixture) TestQueryQueryRow_NoResult() {
	query := &SelectRow{id: 5}
	err := this.DB.PopulateRow(this.Context(), query)
	this.So(err, better.BeNil)
	this.So(query.value, should.BeEmpty)
}

///////////////////////////////////////////////

type DDL struct {
	totalRows uint64
}

func (this *DDL) Statements() string {
	return `
		DROP TABLE IF EXISTS sqldb_integration_test;
		
		CREATE TABLE sqldb_integration_test (
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT    NOT NULL
		);
		
		INSERT INTO sqldb_integration_test (name) VALUES (?);
		INSERT INTO sqldb_integration_test (name) VALUES (?);
		INSERT INTO sqldb_integration_test (name) VALUES (?);
		INSERT INTO sqldb_integration_test (name) VALUES (?);`
}
func (this *DDL) Parameters() []any {
	return []any{
		"a",
		"b",
		"c",
		"d",
	}
}
func (this *DDL) RowsAffected(rows uint64) {
	this.totalRows += rows
}

///////////////////////////////////////////////

type SelectAll struct {
	Result map[int]string
}

func (this *SelectAll) Statement() string {
	return `
		SELECT id, name
		  FROM sqldb_integration_test;`
}

func (this *SelectAll) Parameters() []any {
	return nil
}

func (this *SelectAll) Scan(scanner sqldb.Scanner) error {
	var id int
	var name string
	defer func() { this.Result[id] = name }()
	return scanner.Scan(&id, &name)
}

///////////////////////////////////////////////

type SelectRow struct {
	id    int
	value string
}

func (this *SelectRow) Statement() string {
	return `SELECT name FROM sqldb_integration_test WHERE id = ?;`
}
func (this *SelectRow) Parameters() []any {
	return []any{this.id}
}
func (this *SelectRow) Scan(scanner sqldb.Scanner) error {
	return scanner.Scan(&this.value)
}

///////////////////////////////////////////////
