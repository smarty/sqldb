#### SMARTY DISCLAIMER: Subject to the terms of the associated license agreement, this software is freely available for your use. This software is FREE, AS IN PUPPIES, and is a gift. Enjoy your new responsibility. This means that while we may consider enhancement requests, we may or may not choose to entertain requests at our sole and absolute discretion.

[![Build Status](https://travis-ci.org/smarty/sqldb.svg?branch=master)](https://travis-ci.org/smarty/sqldb)
[![Code Coverage](https://codecov.io/gh/smarty/sqldb/branch/master/graph/badge.svg)](https://codecov.io/gh/smarty/sqldb)
[![Go Report Card](https://goreportcard.com/badge/github.com/smarty/sqldb)](https://goreportcard.com/report/github.com/smarty/sqldb)
[![GoDoc](https://godoc.org/github.com/smarty/sqldb?status.svg)](http://godoc.org/github.com/smarty/sqldb)


## Upgrading from v2 to v3

Those upgrading from v2 to v3 may be interested in the following code, which can be used to adapt write operations ('Scripts') to v3 while still using a v2-style interface:

```go
package mysql

import (
	"context"

	"github.com/smarty/sqldb/v3"
)

// Deprecated
type LegacyExecutor interface {
	Execute(context.Context, string, ...any) (uint64, error)
}

// Deprecated
func newLegacyExecutor(handle sqldb.Handle) LegacyExecutor {
	return &legacyExecutor{handle: handle}
}

// Deprecated
type legacyExecutor struct {
	handle sqldb.Handle
}

// Deprecated
func (this *legacyExecutor) Execute(ctx context.Context, statement string, args ...any) (uint64, error) {
	script := &rowCountScript{
		BaseScript: sqldb.BaseScript{
			Text: statement,
			Args: args,
		},
	}
	err := this.handle.Execute(ctx, script)
	return script.rowsAffectedCount, err
}

// Deprecated
type rowCountScript struct {
	sqldb.BaseScript
	rowsAffectedCount uint64
}

// Deprecated
func (this *rowCountScript) RowsAffected(rowCount uint64) {
	this.rowsAffectedCount += rowCount
}
```