package log

import (
	logging "github.com/op/go-logging"
)

type twlogger struct {
}

// MustGetLogger 返回 logger
func MustGetLogger(module string) *logging.Logger {
	return logging.MustGetLogger(module)
}
