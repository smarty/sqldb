package sqldb_test

import (
	"context"
	"time"

	"github.com/smarty/sqldb/v3"
)

// Business Event:
type FooEstablished struct {
	Timestamp time.Time
	FooID     uint64
}

// Storage Operation:
type LoadFooName struct {
	FooID  uint64
	Result struct {
		FooID   uint64
		FooName string
	}
}

type Mapper struct {
	db sqldb.DBTx
}

func (this Mapper) fooEstablished(ctx context.Context, operation FooEstablished) (uint64, error) {
	return sqldb.ExecuteStatements(ctx, this.db, `
		INSERT
		  INTO Foos
		       ( foo_id, created, foo_name )
		VALUES ( ?, ?, '' )
		    ON DUPLICATE KEY 
		UPDATE created = created;`,
		operation.FooID, operation.Timestamp,
	)
}
func (this Mapper) fooEstablished_Prepared(ctx context.Context, operation FooEstablished) (uint64, error) {
	statement, err := this.db.PrepareContext(ctx, `
		INSERT
		  INTO Foos
		       ( foo_id, created, foo_name )
		VALUES ( ?, ?, '' )
		    ON DUPLICATE KEY 
		UPDATE created = created;`,
	)
	if err != nil {
		return 0, err
	}
	return sqldb.RowsAffected(statement.ExecContext(ctx, operation.FooID, operation.Timestamp))
}

func (this Mapper) loadFooName(ctx context.Context, operation *LoadFooName) error {
	rows, err := this.db.QueryContext(ctx, `
		SELECT foo_id, foo_name
		  FROM Foos 
		 WHERE foo_id = ?;`,
		operation.FooID,
	)
	return sqldb.BindAll(rows, err, func(source sqldb.Scanner) error {
		return source.Scan(&operation.Result.FooID, &operation.Result.FooName)
	})
}
func (this Mapper) loadFooName_Prepared(ctx context.Context, operation *LoadFooName) error {
	statement, err := this.db.PrepareContext(ctx, `
		SELECT foo_id, foo_name
		  FROM Foos 
		 WHERE foo_id = ?;`,
	)
	if err != nil {
		return err
	}
	rows, err := statement.QueryContext(ctx, operation.FooID)
	return sqldb.BindAll(rows, err, func(source sqldb.Scanner) error {
		return source.Scan(&operation.Result.FooID, &operation.Result.FooName)
	})
}
