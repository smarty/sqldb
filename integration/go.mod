module github.com/smarty/sqldb/integration

go 1.25

require (
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/smarty/gunit/v2 v2.0.0-20250910224800-b24d6a1628bf
	github.com/smarty/sqldb/v3 v3.0.0
)

replace github.com/smarty/sqldb/v3 => ../
