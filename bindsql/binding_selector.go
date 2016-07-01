package bindsql

import "github.com/smartystreets/sqldb"

type BindingSelector struct {
	selector         sqldb.Selector
	panicOnBindError bool
}

func NewBindingSelector(selector sqldb.Selector, panicOnBindError bool) *BindingSelector {
	return &BindingSelector{selector: selector, panicOnBindError: panicOnBindError}
}

func (this *BindingSelector) Select(binder Binder, statement string, parameters ...interface{}) error {
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
