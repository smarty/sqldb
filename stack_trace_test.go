package sqldb

import (
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

	actual = actual[len(actual)-1000:]       // last 1000 characters
	expected = expected[len(expected)-1000:] // last 1000 characters

	this.So(actual, should.Equal, expected)
}

func (this *StackTraceFixture) TestWhenNonNil_ReturnsPreSetMessage() {
	this.stack = ContrivedStackTrace("HELLO")
	this.So(this.stack.StackTrace(), should.Equal, "HELLO")
}
