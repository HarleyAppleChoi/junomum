package logging

import (
	"github.com/forbole/flowJuno/modules/modules"
	"github.com/forbole/flowJuno/types"
	"github.com/onflow/flow-go-sdk"
)

const (
	LogKeyModule  = "module"
	LogKeyHeight  = "height"
	LogKeyTxHash  = "tx_hash"
	LogKeyMsgType = "msg_type"
)

// Logger defines a function that takes an error and logs it.
type Logger interface {
	SetLogLevel(level string) error
	SetLogFormat(format string) error
	LogGenesisError(module modules.Module, err error)
	LogBLockError(module modules.Module, block *flow.Block, err error)
	LogTxError(module modules.Module, tx types.Txs, err error)
}

// logger represents the currently used logger
var logger Logger = &defaultLogger{}

// SetLogger sets the given logger as the one to be used
func SetLogger(l Logger) {
	logger = l
}

func SetLogLevel(level string) error {
	return logger.SetLogLevel(level)
}
func SetLogFormat(format string) error {
	return logger.SetLogFormat(format)
}

// LogGenesisError logs the error returned while handling the genesis of the given module
func LogGenesisError(module modules.Module, err error) {
	logger.LogGenesisError(module, err)
}

// LogBLockError logs the error returned while handling the given block inside the specified module
func LogBLockError(module modules.Module, block *flow.Block, err error) {
	logger.LogBLockError(module, block, err)
}

// LogTxError logs the error returned while handling the provided transaction inside the given module
func LogTxError(module modules.Module, tx types.Txs, err error) {
	logger.LogTxError(module, tx, err)
}
