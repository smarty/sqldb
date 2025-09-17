package integration

import (
	"database/sql"
	"testing"

	"github.com/smarty/gunit/v2"
	"github.com/smarty/gunit/v2/better"
	"github.com/smarty/gunit/v2/should"
	"github.com/smarty/sqldb/v3"

	_ "github.com/mattn/go-sqlite3"
)

const CreateInsert = `

DROP TABLE IF EXISTS sqldb_integration_test;

CREATE TABLE sqldb_integration_test (
	id   INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT    NOT NULL
);

INSERT INTO sqldb_integration_test (name) VALUES (?);
INSERT INTO sqldb_integration_test (name) VALUES (?);
INSERT INTO sqldb_integration_test (name) VALUES (?);
INSERT INTO sqldb_integration_test (name) VALUES (?);
`
const QuerySelectAll = `
SELECT id, name
  FROM sqldb_integration_test;
`

func Test(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	gunit.So(t, err, better.BeNil)

	tx, err := db.BeginTx(t.Context(), nil)
	gunit.So(t, err, better.BeNil)

	err = sqldb.ExecuteScript(t.Context(), tx, CreateInsert,
		"a",
		"b",
		"c",
		"d",
	)
	gunit.So(t, err, better.BeNil)

	rows, err := tx.QueryContext(t.Context(), QuerySelectAll)
	gunit.So(t, err, should.BeNil)

	records := map[int]string{}
	err = sqldb.BindAll(rows, err, func(scanner sqldb.Scanner) error {
		var id int
		var name string
		defer func() { records[id] = name }()
		return scanner.Scan(&id, &name)
	})
	gunit.So(t, err, better.BeNil)
	gunit.So(t, records, should.Equal, map[int]string{
		1: "a",
		2: "b",
		3: "c",
		4: "d",
	})

	var id int
	var value string
	err = sqldb.ScanOptionalRow(tx.QueryRowContext(t.Context(), QuerySelectAll), &id, &value)
	gunit.So(t, err, better.BeNil)
	gunit.So(t, id, should.Equal, 1)
	gunit.So(t, value, should.Equal, "a")

	err = tx.Rollback()
	gunit.So(t, err, better.BeNil)
}
