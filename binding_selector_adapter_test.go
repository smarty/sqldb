package sqldb

import (
	"errors"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BindingSelectorAdapterFixture struct {
	*gunit.Fixture

	fakeResult        *FakeResult
	fakeInnerSelector *FakeSelector
	selector          *BindingSelectorAdapter
}

func (this *BindingSelectorAdapterFixture) Setup() {
	this.fakeResult = &FakeResult{}
	this.fakeInnerSelector = &FakeSelector{selectResult: this.fakeResult}
	this.selector = NewBindingSelectorAdapter(this.fakeInnerSelector, false)
}

///////////////////////////////////////////////////////////////

func (this *BindingSelectorAdapterFixture) TestFailedSelectReturnsError() {
	this.fakeInnerSelector.selectError = errors.New("")

	err := this.selector.Select(nil, "query", 1, 2, 3)

	this.So(err, should.Equal, this.fakeInnerSelector.selectError)
	this.So(this.fakeInnerSelector.selects, should.Equal, 1)
	this.So(this.fakeInnerSelector.statement, should.Equal, "query")
	this.So(this.fakeInnerSelector.parameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *BindingSelectorAdapterFixture) TestEmptyResult() {
	err := this.selector.Select(nil, "query", 1, 2, 3)
	this.So(err, should.BeNil)
	this.So(this.fakeInnerSelector.selects, should.Equal, 1)
	this.So(this.fakeResult.nextCalls, should.Equal, 1)
	this.So(this.fakeResult.closeCalls, should.Equal, 1)
}

func (this *BindingSelectorAdapterFixture) TestResultErrorClosesAndReturnsError() {
	this.fakeResult.iterations = 1
	this.fakeResult.errError = errors.New("")

	err := this.selector.Select(nil, "query", 1, 2, 3)
	this.So(err, should.Equal, this.fakeResult.errError)
	this.So(this.fakeInnerSelector.selects, should.Equal, 1)
	this.So(this.fakeResult.nextCalls, should.Equal, 1)
	this.So(this.fakeResult.errCalls, should.Equal, 1)
	this.So(this.fakeResult.closeCalls, should.Equal, 1)
}

func (this *BindingSelectorAdapterFixture) TestScanErrorClosesAndReturnsError() {
	this.fakeResult.iterations = 1
	this.fakeResult.scanError = errors.New("")

	err := this.selector.Select(func(source Scanner) error {
		return source.Scan()
	}, "query", 1, 2, 3)

	this.So(err, should.Equal, this.fakeResult.scanError)
	this.So(this.fakeInnerSelector.selects, should.Equal, 1)
	this.So(this.fakeResult.nextCalls, should.Equal, 1)
	this.So(this.fakeResult.errCalls, should.Equal, 1)
	this.So(this.fakeResult.scanCalls, should.Equal, 1)
	this.So(this.fakeResult.closeCalls, should.Equal, 1)
}

func (this *BindingSelectorAdapterFixture) TestScanErrorClosesAndPanicsWhenConfigured() {
	this.selector.panicOnBindError = true
	this.fakeResult.iterations = 1
	this.fakeResult.scanError = errors.New("")

	this.So(func() {
		this.selector.Select(func(source Scanner) error {
			return source.Scan()
		}, "query", 1, 2, 3)
	}, should.Panic)
}

///////////////////////////////////////////////////////////////

type FakeSelector struct {
	selects      int
	statement    string
	parameters   []interface{}
	selectResult *FakeResult
	selectError  error
}

func (this *FakeSelector) Select(statement string, parameters ...interface{}) (SelectResult, error) {
	this.selects++
	this.statement = statement
	this.parameters = parameters
	return this.selectResult, this.selectError
}

///////////////////////////////////////////////////////////////

type FakeResult struct {
	nextCalls  int
	errCalls   int
	closeCalls int
	scanCalls  int
	iterations int

	errError   error
	closeError error
	scanError  error
}

func (this *FakeResult) Next() bool {
	this.nextCalls++
	return this.iterations >= this.nextCalls
}

func (this *FakeResult) Err() error {
	this.errCalls++
	return this.errError
}

func (this *FakeResult) Close() error {
	this.closeCalls++
	return this.closeError
}

func (this *FakeResult) Scan(target ...interface{}) error {
	this.scanCalls++
	return this.scanError
}
