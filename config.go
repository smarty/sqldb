package sqldb

type configuration struct {
	logger    logger
	threshold int
}
type option func(*configuration)

var Options singleton

type singleton struct{}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, item := range Options.defaults(options...) {
			item(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	return append([]option{
		Options.Logger(&nop{}),
		Options.PreparationThreshold(10),
	}, options...)
}

func (singleton) Logger(logger logger) option {
	return func(this *configuration) { this.logger = logger }
}

// PreparationThreshold specifies the number of times a give sql statement can be
// executed before it will be transitioned to a prepared statement. Passing a negative
// value will disable any use of prepared statements.
func (singleton) PreparationThreshold(n int) option {
	return func(this *configuration) { this.threshold = n }
}

type nop struct{}

func (nop) Printf(string, ...any) {}
