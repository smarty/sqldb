package sqldb

import "strings"

//go:generate go install github.com/smartystreets/gunit/gunit
//go:generate gunit

///////////////////////////////////////////////////////////////

type FakeInnerConnectionPool struct {
	pingCalls         int
	pingError         error
	transactionCalls  int
	transaction       *FakeInnerTransaction
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

func (this *FakeInnerConnectionPool) Ping() error {
	this.pingCalls++
	return this.pingError
}

func (this *FakeInnerConnectionPool) BeginTransaction() (Transaction, error) {
	this.transactionCalls++
	return this.transaction, this.transactionError
}

func (this *FakeInnerConnectionPool) Close() error {
	this.closeCalls++
	return this.closeError
}

func (this *FakeInnerConnectionPool) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeInnerConnectionPool) Select(statement string, parameters ...interface{}) (SelectResult, error) {
	this.selectCalls++
	this.selectStatement = statement
	this.selectParameters = parameters
	return this.selectResult, this.selectError
}

///////////////////////////////////////////////////////////////

type FakeInnerTransaction struct {
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

func (this *FakeInnerTransaction) Commit() error {
	this.commitCalls++
	return this.commitError
}

func (this *FakeInnerTransaction) Rollback() error {
	this.rollbackCalls++
	return this.rollbackError
}

func (this *FakeInnerTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.executeCalls++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}

func (this *FakeInnerTransaction) Select(statement string, parameters ...interface{}) (SelectResult, error) {
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
