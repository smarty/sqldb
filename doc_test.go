package sqldb

import "strings"

///////////////////////////////////////////////////////////////

type FakeConnectionPool struct {
	pingCalls         int
	pingError         error
	transactionCalls  int
	transaction       *FakeTransaction
	transactionError  error
	closeCalls        int
	closeError        error
	selectCalls       int
	selectStatement   string
	selectParameters  []interface{}
	selectResult      *FakeSelectResult
	selectError       error
	executeCalls      int
	executeStatement  string
	executeParameters []interface{}
	executeResult     uint64
	executeError      error
}

func (this *FakeConnectionPool) Ping() error {
	this.pingCalls++
	return this.pingError
}

func (this *FakeConnectionPool) BeginTransaction() (Transaction, error) {
	this.transactionCalls++
	return this.transaction, this.transactionError
}

func (this *FakeConnectionPool) Close() error {
	this.closeCalls++
	return this.closeError
}

func (this *FakeConnectionPool) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeConnectionPool) Select(statement string, parameters ...interface{}) (SelectResult, error) {
	this.selectCalls++
	this.selectStatement = statement
	this.selectParameters = parameters
	return this.selectResult, this.selectError
}

///////////////////////////////////////////////////////////////

type FakeTransaction struct {
	commitCalls       int
	commitError       error
	rollbackCalls     int
	rollbackError     error
	selectCalls       int
	selectStatement   string
	selectParameters  []interface{}
	selectResult      *FakeSelectResult
	selectError       error
	executeCalls      int
	executeStatement  string
	executeParameters []interface{}
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

func (this *FakeTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeTransaction) Select(statement string, parameters ...interface{}) (SelectResult, error) {
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

func (this *FakeSelectResult) Scan(target ...interface{}) error {
	this.scanCalls++
	return this.scanError
}

///////////////////////////////////////////////////////////////

type FakeExecutor struct {
	affected       uint64
	errorsToReturn []error
	statements     []string
	parameters     [][]interface{}
}

func (this *FakeExecutor) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.statements = append(this.statements, strings.TrimSpace(statement))
	this.parameters = append(this.parameters, parameters)

	if len(this.statements) <= len(this.errorsToReturn) {
		return this.affected, this.errorsToReturn[len(this.statements)-1]
	}

	return this.affected, nil
}

///////////////////////////////////////////////////////////////

type FakeBindingConnectionPool struct {
	pingCalls         int
	pingError         error
	transactionCalls  int
	transaction       *FakeBindingTransaction
	transactionError  error
	closeCalls        int
	closeError        error
	selectCalls       int
	selectBinder      Binder
	selectStatement   string
	selectParameters  []interface{}
	selectResult      *FakeSelectResult
	selectError       error
	executeCalls      int
	executeStatement  string
	executeParameters []interface{}
	executeResult     uint64
	executeError      error
}

func (this *FakeBindingConnectionPool) Ping() error {
	this.pingCalls++
	return this.pingError
}

func (this *FakeBindingConnectionPool) BeginTransaction() (BindingTransaction, error) {
	this.transactionCalls++
	return this.transaction, this.transactionError
}

func (this *FakeBindingConnectionPool) Close() error {
	this.closeCalls++
	return this.closeError
}

func (this *FakeBindingConnectionPool) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeBindingConnectionPool) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	this.selectCalls++
	this.selectBinder = binder
	this.selectStatement = statement
	this.selectParameters = parameters
	return this.selectError
}

///////////////////////////////////////////////////////////////

type FakeBindingTransaction struct {
	commitCalls       int
	commitError       error
	rollbackCalls     int
	rollbackError     error
	selectCalls       int
	selectBinder      Binder
	selectStatement   string
	selectParameters  []interface{}
	selectResult      *FakeSelectResult
	selectError       error
	executeCalls      int
	executeStatement  string
	executeParameters []interface{}
	executeResult     uint64
	executeError      error
}

func (this *FakeBindingTransaction) Commit() error {
	this.commitCalls++
	return this.commitError
}

func (this *FakeBindingTransaction) Rollback() error {
	this.rollbackCalls++
	return this.rollbackError
}

func (this *FakeBindingTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeBindingTransaction) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	this.selectCalls++
	this.selectBinder = binder
	this.selectStatement = statement
	this.selectParameters = parameters
	return this.selectError
}

///////////////////////////////////////////////////////////////
