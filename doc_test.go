package sqldb

//go:generate go install github.com/smartystreets/gunit/gunit
//go:generate gunit

type FakeInnerTransaction struct {
	commit            int
	commitError       error
	rollback          int
	rollbackError     error
	selects           int
	selectStatement   string
	selectParameters  []interface{}
	selectResult      *FakeSelectResult
	selectError       error
	execute           int
	executeStatement  string
	executeParameters []interface{}
	executeResult     uint64
	executeError      error
}

func (this *FakeInnerTransaction) Commit() error {
	this.commit++
	return this.commitError
}
func (this *FakeInnerTransaction) Rollback() error {
	this.rollback++
	return this.rollbackError
}

func (this *FakeInnerTransaction) Execute(statement string, parameters ...interface{}) (uint64, error) {
	this.execute++
	this.executeStatement = statement
	this.executeParameters = parameters
	return this.executeResult, this.executeError
}
func (this *FakeInnerTransaction) Select(statement string, parameters ...interface{}) (SelectResult, error) {
	this.selects++
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
