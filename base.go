package sqldb

// BaseScript is a bare-minimum implementation of Script.
type BaseScript struct {
	Text string
	Args []any
}

func (this BaseScript) Statements() string {
	return this.Text
}
func (this BaseScript) Parameters() []any {
	return this.Args
}

// BaseQuery is a bare-minium, partial implementation of Query.
// Users are invited to embed it on types that define a Scan method,
// thus completing the Query implementation.
type BaseQuery struct {
	Text string
	Args []any
}

func (this BaseQuery) Statement() string {
	return this.Text
}
func (this BaseQuery) Parameters() []any {
	return this.Args
}
