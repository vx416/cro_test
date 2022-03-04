package logger

import (
	"context"

	sqlxxlog "github.com/vx416/sqlxx/logger"
)

type (
	// LoggerKey define context key for logger
	LoggerKey struct{}
)

var _logger = newDefaultZap()

// Logger interface
type Logger interface {
	Err(err error) Logger
	Fields(fields map[string]interface{}) Logger
	Field(key string, val interface{}) Logger
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
	Caller(stack int) Logger

	Attach(ctx context.Context) context.Context
	Printf(msg string, args ...interface{})
}

// SetGlobal set global logger
func SetGlobal(logger Logger) Logger {
	_logger = logger
	return _logger
}

// Get get global logger
func Get() Logger {
	return _logger
}

// Attach attach logger instance into context
func Attach(ctx context.Context, logger Logger) context.Context {
	ctx = sqlxxlog.AttachLogger(ctx, SqlxxLogger{logger})
	return context.WithValue(ctx, LoggerKey{}, logger)
}

// Ctx get logger instacne from context
func Ctx(ctx context.Context) Logger {
	val := ctx.Value(LoggerKey{})
	logger, ok := val.(Logger)
	if !ok {
		return Get()
	}

	return logger
}

type SqlxxLogger struct {
	Logger
}

func (l SqlxxLogger) Warn(s string, fields map[string]interface{}) {
	l.Fields(fields).Warn(s)
}
func (l SqlxxLogger) Info(s string, fields map[string]interface{}) {
	l.Fields(fields).Info(s)
}
func (l SqlxxLogger) Debug(s string, fields map[string]interface{}) {
	l.Fields(fields).Debug(s)
}
func (l SqlxxLogger) Error(s string, fields map[string]interface{}) {
	l.Fields(fields).Error(s)
}
