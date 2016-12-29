package sqldb

import (
	"fmt"
	"runtime/debug"
)

// StackTrace, like github.com/smartystreets/clock.Clock performs in production mode
// when used as a nil pointer struct field. When non-nil, it returns a preset value.
// This is useful during testing when asserting on simple, deterministic values is helpful.
type StackTrace struct {
	trace string
}

func ContrivedStackTrace(trace string) *StackTrace {
	return &StackTrace{trace: trace}
}

func (this *StackTrace) Wrap(err error) error {
	if err != nil {
		err = fmt.Errorf("%s\nStack Trace:\n%s", err.Error(), this.StackTrace())
	}
	return err
}

func (this *StackTrace) StackTrace() string {
	if this == nil {
		return string(debug.Stack())
	}
	return this.trace
}
