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
const query = `
SELECT id, name
  FROM sqldb_integration_test;
`

func Test(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	gunit.So(t, err, better.BeNil)

	tx, err := db.BeginTx(t.Context(), nil)
	gunit.So(t, err, better.BeNil)

	count, err := sqldb.ExecuteStatements(t.Context(), tx, CreateInsert,
		"a",
		"b",
		"c",
		"d",
	)
	gunit.So(t, err, better.BeNil)
	gunit.So(t, count, should.Equal, 4)

	rows, err := tx.QueryContext(t.Context(), query)
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

	err = tx.Rollback()
	gunit.So(t, err, better.BeNil)
}
