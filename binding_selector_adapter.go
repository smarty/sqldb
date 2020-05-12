package sqldb

type BindingSelectorAdapter struct {
	selector         Selector
	panicOnBindError bool
}

func NewBindingSelectorAdapter(selector Selector, panicOnBindError bool) *BindingSelectorAdapter {
	return &BindingSelectorAdapter{selector: selector, panicOnBindError: panicOnBindError}
}

func (this *BindingSelectorAdapter) BindSelect(binder Binder, statement string, parameters ...interface{}) error {
	result, err := this.selector.Select(statement, parameters...)
	if err != nil {
		return err
	}

	for result.Next() {
		if err := result.Err(); err != nil {
			_ = result.Close()
			return err
		}

		if err := binder(result); err != nil {
			_ = result.Close()
			if this.panicOnBindError {
				panic(err)
			} else {
				return err
			}
		}
	}

	return result.Close()
}
