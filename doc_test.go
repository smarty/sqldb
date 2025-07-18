package sqldb

import (
	"context"
	"strings"
)

///////////////////////////////////////////////////////////////

type FakeConnectionPool struct {
	pingCalls int
	pingError error

	transactionCalls int
	transaction      *FakeTransaction
	transactionError error

	closeCalls int
	closeError error

	selectCalls      int
	selectStatement  string
	selectParameters []any
	selectResult     *FakeSelectResult
	selectError      error

	executeCalls      int
	executeStatement  string
	executeParameters []any
	executeResult     uint64
	executeError      error
}

func (this *FakeConnectionPool) Ping(_ context.Context) error {
	this.pingCalls++
	return this.pingError
}

func (this *FakeConnectionPool) BeginTransaction(_ context.Context) (Transaction, error) {
	this.transactionCalls++
	return this.transaction, this.transactionError
}

func (this *FakeConnectionPool) Close() error {
	this.closeCalls++
	return this.closeError
}

func (this *FakeConnectionPool) Execute(_ context.Context, statement string, parameters ...any) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeConnectionPool) Select(_ context.Context, statement string, parameters ...any) (SelectResult, error) {
	this.selectCalls++
	this.selectStatement = statement
	this.selectParameters = parameters
	return this.selectResult, this.selectError
}

///////////////////////////////////////////////////////////////

type FakeTransaction struct {
	commitCalls int
	commitError error

	rollbackCalls int
	rollbackError error

	selectCalls      int
	selectStatement  string
	selectParameters []any
	selectResult     *FakeSelectResult
	selectError      error

	executeCalls      int
	executeStatement  string
	executeParameters []any
	executeResult     uint64
	executeError      error
}

func (this *FakeTransaction) Commit() error {
	this.commitCalls++
	return this.commitError
}

func (this *FakeTransaction) Rollback() error {
	this.rollbackCalls++
	return this.rollbackError
}

func (this *FakeTransaction) Execute(_ context.Context, statement string, parameters ...any) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeTransaction) Select(_ context.Context, statement string, parameters ...any) (SelectResult, error) {
	this.selectCalls++
	this.selectStatement = statement
	this.selectParameters = parameters
	return this.selectResult, this.selectError
}

///////////////////////////////////////////////////////////////

type FakeSelectResult struct {
	nextCalls  int
	errCalls   int
	closeCalls int
	scanCalls  int
	iterations int

	errError   error
	closeError error
	scanError  error
}

func (this *FakeSelectResult) Next() bool {
	this.nextCalls++
	return this.iterations >= this.nextCalls
}

func (this *FakeSelectResult) Err() error {
	this.errCalls++
	return this.errError
}

func (this *FakeSelectResult) Close() error {
	this.closeCalls++
	return this.closeError
}

func (this *FakeSelectResult) Scan(_ ...any) error {
	this.scanCalls++
	return this.scanError
}

///////////////////////////////////////////////////////////////

type FakeExecutor struct {
	affected       uint64
	errorsToReturn []error
	statements     []string
	parameters     [][]any
}

func (this *FakeExecutor) Execute(_ context.Context, statement string, parameters ...any) (uint64, error) {
	this.statements = append(this.statements, strings.TrimSpace(statement))
	this.parameters = append(this.parameters, parameters)

	if len(this.statements) <= len(this.errorsToReturn) {
		return this.affected, this.errorsToReturn[len(this.statements)-1]
	}

	return this.affected, nil
}

///////////////////////////////////////////////////////////////

type FakeBindingConnectionPool struct {
	pingCalls int
	pingError error

	transactionCalls int
	transaction      *FakeBindingTransaction
	transactionError error

	closeCalls int
	closeError error

	selectCalls      int
	selectBinder     Binder
	selectStatement  string
	selectParameters []any
	selectResult     *FakeSelectResult
	selectError      error

	executeCalls      int
	executeStatement  string
	executeParameters []any
	executeResult     uint64
	executeError      error
}

func (this *FakeBindingConnectionPool) Ping(_ context.Context) error {
	this.pingCalls++
	return this.pingError
}

func (this *FakeBindingConnectionPool) BeginTransaction(_ context.Context) (BindingTransaction, error) {
	this.transactionCalls++
	return this.transaction, this.transactionError
}

func (this *FakeBindingConnectionPool) Close() error {
	this.closeCalls++
	return this.closeError
}

func (this *FakeBindingConnectionPool) Execute(_ context.Context, statement string, parameters ...any) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeBindingConnectionPool) BindSelect(_ context.Context, binder Binder, statement string, parameters ...any) error {
	this.selectCalls++
	this.selectBinder = binder
	this.selectStatement = statement
	this.selectParameters = parameters
	return this.selectError
}

///////////////////////////////////////////////////////////////

type FakeBindingTransaction struct {
}

func (this *FakeBindingTransaction) Commit() error {
	panic("Not called")
}

func (this *FakeBindingTransaction) Rollback() error {
	panic("Not called")
}

func (this *FakeBindingTransaction) Execute(_ context.Context, _ string, _ ...any) (uint64, error) {
	panic("Not called")
}

func (this *FakeBindingTransaction) BindSelect(_ context.Context, _ Binder, _ string, _ ...any) error {
	panic("Not called")
}

///////////////////////////////////////////////////////////////
