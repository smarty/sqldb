package sqldb

import (
	"database/sql"
)

type txConfig struct {
	splitStatement    bool
	panicOnBindError  bool
	stackTraceOnError bool
	parameterPrefix   string
}

func NewTransaction(handle *sql.Tx, options ...txOption) Transaction {
	var config txConfig
	TxOptions.apply(options...)(&config)
	return newTx(handle, config)
}
func NewBindingTransaction(handle *sql.Tx, options ...txOption) BindingTransaction {
	var config txConfig
	TxOptions.apply(options...)(&config)
	return newBindingTx(handle, config)
}
func newTx(handle *sql.Tx, config txConfig) Transaction {
	var tx Transaction = NewLibraryTransactionAdapter(handle)

	if config.splitStatement {
		tx = NewSplitStatementTransaction(tx, config.parameterPrefix)
	}

	if config.stackTraceOnError {
		tx = NewStackTraceTransaction(tx)
	}

	return tx
}
func newBindingTx(handle *sql.Tx, config txConfig) BindingTransaction {
	inner := newTx(handle, config)
	return NewBindingTransactionAdapter(inner, config.panicOnBindError)
}

var TxOptions txSingleton

type txSingleton struct{}
type txOption func(*txConfig)

func (txSingleton) PanicOnBindError(value bool) txOption {
	return func(this *txConfig) { this.panicOnBindError = value }
}
func (txSingleton) MySQL() txOption {
	return func(this *txConfig) { this.splitStatement = true; this.parameterPrefix = "?" }
}
func (txSingleton) ParameterPrefix(value string) txOption {
	return func(this *txConfig) { this.parameterPrefix = value }
}
func (txSingleton) SplitStatement(value bool) txOption {
	return func(this *txConfig) { this.splitStatement = value }
}
func (txSingleton) StackTraceErrDiagnostics(value bool) txOption {
	return func(this *txConfig) { this.stackTraceOnError = value }
}

func (txSingleton) apply(txOptions ...txOption) txOption {
	return func(this *txConfig) {
		for _, txOption := range TxOptions.defaults(txOptions...) {
			txOption(this)
		}
	}
}
func (txSingleton) defaults(txOptions ...txOption) []txOption {
	const defaultStackTraceErrDiagnostics = true
	const defaultPanicOnBindError = true
	const defaultSplitStatement = true
	const defaultParameterPrefix = "?"

	return append([]txOption{
		TxOptions.PanicOnBindError(defaultPanicOnBindError),
		TxOptions.StackTraceErrDiagnostics(defaultStackTraceErrDiagnostics),
		TxOptions.ParameterPrefix(defaultParameterPrefix),
		TxOptions.SplitStatement(defaultSplitStatement),
	}, txOptions...)
}
