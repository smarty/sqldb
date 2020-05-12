package sqldb

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type SplitStatementExecutor struct {
	Executor
	delimiter string
}

func NewSplitStatementExecutor(actual Executor, delimiter string) *SplitStatementExecutor {
	return &SplitStatementExecutor{Executor: actual, delimiter: delimiter}
}

func (this *SplitStatementExecutor) Execute(ctx context.Context, statement string, parameters ...interface{}) (uint64, error) {
	if argumentCount := strings.Count(statement, this.delimiter); argumentCount != len(parameters) {
		return 0, fmt.Errorf("%w: Expected: %d, received %d", ErrArgumentCountMismatch, argumentCount, len(parameters))
	}

	var count uint64
	index := 0
	for _, statement = range strings.Split(statement, ";") {
		if len(strings.TrimSpace(statement)) == 0 {
			continue
		}

		statement += ";" // terminate the statement
		indexOffset := strings.Count(statement, this.delimiter)
		if affected, err := this.Executor.Execute(ctx, statement, parameters[index:index+indexOffset]...); err != nil {
			return 0, err
		} else {
			count += affected
		}

		index += indexOffset
	}

	return count, nil
}

var ErrArgumentCountMismatch = errors.New("the number of arguments supplied does not match the statement")
