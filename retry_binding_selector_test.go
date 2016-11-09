package sqldb

import (
	"errors"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/clock"
	"github.com/smartystreets/gunit"
)

func TestRetryBindingSelectorFixture(t *testing.T) {
        gunit.Run(new(RetryBindingSelectorFixture), t)
}

type RetryBindingSelectorFixture struct {
	*gunit.Fixture

	inner    *FakeRetrySelector
	selector *RetryBindingSelector
}

func (this *RetryBindingSelectorFixture) Setup() {
	this.inner = &FakeRetrySelector{}
	this.selector = NewRetryBindingSelector(this.inner, time.Second)
	this.selector.sleep = clock.StayAwake()
}

///////////////////////////////////////////////////////////////

func (this *RetryBindingSelectorFixture) TestSelectWithoutErrors() {
	err := this.selector.BindSelect(nil, "statement", 1, 2, 3)

	this.So(err, should.Equal, err)
	this.So(this.inner.count, should.Equal, 1)
	this.So(this.inner.statement, should.Equal, "statement")
	this.So(this.inner.parameters, should.Resemble, []interface{}{1, 2, 3})
}

func (this *RetryBindingSelectorFixture) TestRetryUntilSuccess() {
	this.inner.errorCount = 4

	err := this.selector.BindSelect(nil, "statement", 1, 2, 3)

	this.So(err, should.Equal, err)
	this.So(this.inner.count, should.Equal, 5) // last attempt is successful
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

func (this *FakeRetrySelector) Ping() error {
	panic("Should not be called.")
}

func (this *FakeRetrySelector) BeginTransaction() (BindingTransaction, error) {
	panic("Should not be called.")
}

func (this *FakeRetrySelector) Close() error {
	panic("Should not be called.")
}

func (this *FakeRetrySelector) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
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

func (this *FakeRetrySelector) Execute(string, ...interface{}) (uint64, error) {
	panic("Should not be called.")
}
