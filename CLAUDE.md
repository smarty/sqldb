# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make test      # Format, then run unit tests (race detection, 1s timeout) + integration tests
make fmt       # Format code and tidy go.mod
make compile   # Build all packages
make build     # Run test, then compile
```

To run a single test:
```bash
go test -run TestName ./...
```

Integration tests require CGO (uses go-sqlite3):
```bash
go test github.com/smarty/sqldb/integration
```

## Architecture

This is a lightweight Go database abstraction library (v3) providing high-level interfaces over `database/sql`. Zero external dependencies in core module.

### Core Interfaces

**Script** - For SQL statements without result rows (INSERT, UPDATE, DELETE):
- `Statements() string` - Returns semicolon-separated SQL statements
- `Parameters() []any` - Returns parameters interleaved across all statements

**Query** - For SELECT statements returning rows:
- `Statement() string` - Single SQL query
- `Parameters() []any` - Query parameters
- `Scan(Scanner) error` - Called for each row

**Handle** - Main operations:
- `Execute(ctx, ...Script)` - Run scripts
- `Populate(ctx, ...Query)` - Run queries (multiple rows)
- `PopulateRow(ctx, ...Query)` - Run queries (single row, errors if none)

### Key Implementation Details

- **Prepared statement caching**: Uses FNV-64a hash of SQL text. After configurable threshold (default 1), statements become prepared. Capped at 64k statements.
- **Thread safety**: RWMutex for concurrent access to prepared statement cache
- **Parameter interleaving**: Scripts can contain multiple semicolon-separated statements; parameters are distributed across them using Go iterators
- **BaseScript/BaseQuery**: Helper types in `base.go` implementing the interfaces

### Optional Hooks

Scripts can implement:
- `RowsAffected(uint64)` - Receive row count per statement
- `ExpectedRowsAffected() uint64` - Optimistic concurrency check; returns `ErrOptimisticConcurrencyCheckFailed` on mismatch

### File Layout

- `db.go` - Handle implementation with prepared statement caching
- `contracts.go` - Interface definitions
- `base.go` - BaseScript/BaseQuery helpers
- `config.go` - Options pattern (Logger, PreparationThreshold)
- `integration/` - Integration tests using SQLite
