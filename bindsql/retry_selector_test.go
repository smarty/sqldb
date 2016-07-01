package bindsql

import (
	"errors"
	"time"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
)

type RetrySelectorFixture struct {
	*gunit.Fixture

	fakeSelector *FakeRetrySelector
	selector     *RetrySelector
}

func (this *RetrySelectorFixture) Setup() {
	this.fakeSelector = &FakeRetrySelector{}
	this.selector = NewRetrySelector(this.fakeSelector, time.Second)
	this.selector.sleep = clock.StayAwake()
}

///////////////////////////////////////////////////////////////

func (this *RetrySelectorFixture) TestSelectWithoutErrors() {
	err := this.selector.Select(nil, "statement", 1, 2, 3)

	this.So(err, should.Equal, err)
	this.So(this.fakeSelector.count, should.Equal, 1)
	this.So(this.fakeSelector.statement, should.Equal, "statement")
	this.So(this.fakeSelector.parameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *RetrySelectorFixture) TestRetryUntilSuccess() {
	this.fakeSelector.errorCount = 4

	err := this.selector.Select(nil, "statement", 1, 2, 3)

	this.So(err, should.Equal, err)
	this.So(this.fakeSelector.count, should.Equal, 5) // last attempt is successful
	this.So(this.selector.sleep.Naps, should.Resemble, []time.Duration{
		time.Second,
		time.Second,
		time.Second,
		time.Second,
	})
}

///////////////////////////////////////////////////////////////

type FakeRetrySelector struct {
	count      int
	errorCount int
	binder     Binder
	statement  string
	parameters []interface{}
}

func (this *FakeRetrySelector) Select(binder Binder, statement string, parameters ...interface{}) error {
	if this.binder == nil {
		this.binder = binder
	} else {
		assertions.So(this.binder, should.Equal, binder)
	}

	if this.statement == "" {
		this.statement = statement
	} else {
		assertions.So(this.statement, should.Equal, statement)
	}

	if len(this.parameters) == 0 {
		this.parameters = parameters
	} else {
		assertions.So(this.parameters, should.Resemble, parameters)
	}

	this.count++
	if this.errorCount > 0 && this.count <= this.errorCount {
		return errors.New("")
	} else {
		return nil
	}
}
