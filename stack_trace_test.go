package sqldb

import (
	"errors"
	"runtime/debug"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestStackTraceFixture(t *testing.T) {
	gunit.Run(new(StackTraceFixture), t)
}

type StackTraceFixture struct {
	*gunit.Fixture

	stack *StackTrace
}

func (this *StackTraceFixture) TestWhenNil_ReturnsActualStackTrace() {
	actual := this.stack.StackTrace()
	expected := string(debug.Stack())

	actual = actual[len(actual)-925:]       // last 925 characters
	expected = expected[len(expected)-925:] // last 925 characters

	this.So(actual, should.Equal, expected)
}

func (this *StackTraceFixture) TestWhenNonNil_ReturnsPreSetMessage() {
	this.stack = ContrivedStackTrace("HELLO")
	this.So(this.stack.StackTrace(), should.Equal, "HELLO")
}

func (this *StackTraceFixture) TestWrap_NilErrorReturned() {
	var err error
	err = this.stack.Wrap(err)
	this.So(err, should.BeNil)
}

func (this *StackTraceFixture) TestWrap_NonNilErrorDecorated() {
	this.stack = ContrivedStackTrace("GOPHER STACK")
	err := errors.New("HELLO")
	err = this.stack.Wrap(err)
	this.So(err.Error(), should.Equal, "HELLO\nStack Trace:\nGOPHER STACK")
}
