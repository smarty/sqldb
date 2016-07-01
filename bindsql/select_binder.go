package bindsql

import (
	"github.com/smartystreets/sqldb"
)

type SelectBinder struct {
	selector         sqldb.Selector
	panicOnBindError bool
}

func NewSelectBinder(selector sqldb.Selector, panicOnBindError bool) *SelectBinder {
	return &SelectBinder{selector: selector, panicOnBindError: panicOnBindError}
}

func (this *SelectBinder) Select(binder Binder, statement string, parameters ...interface{}) error {
	result, err := this.selector.Select(statement, parameters...)
	if err != nil {
		return err
	}

	for result.Next() {
		if err := result.Err(); err != nil {
			result.Close()
			return err
		}

		if err := binder(result); err != nil {
			result.Close()
			if this.panicOnBindError {
				panic(err)
			} else {
				return err
			}
		}
	}

	return result.Close()
}
