package logger

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

// Env define logger environment
type Env string

// IsDev make env to lower and check if string is dev
func (env Env) IsDev() bool {
	envStr := strings.ToLower(string(env))

	return envStr == "dev" || envStr == "development"
}

// Level define logger level
type Level string

// ZapLevel convert string to zapcore.Level
func (l Level) ZapLevel() zapcore.Level {
	var (
		zapLevel zapcore.Level
		level    = strings.ToLower(string(l))
	)
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	}
	return zapLevel
}

func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// ColorfulLevelEncoder make level colorful
func ColorfulLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var (
		color int
		bold  bool
	)

	switch l {
	case zapcore.DebugLevel:
		color = colorBlue
	case zapcore.InfoLevel:
		color = colorGreen
	case zapcore.WarnLevel:
		color = colorYellow
	case zapcore.ErrorLevel:
		color = colorRed
	case zapcore.FatalLevel:
		color = colorRed
		bold = true
	case zapcore.PanicLevel:
		color = colorRed
		bold = true
	}

	s := colorize(l.CapitalString(), color)

	if bold {
		s = colorize(s, colorBold)
	}
	enc.AppendString(s)
}

// ColorizeCallerEncoder make caller colorful
func ColorizeCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(colorize(caller.TrimmedPath(), colorCyan))
}

func newDefaultZap() Logger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapConfig.EncoderConfig.EncodeLevel = ColorfulLevelEncoder
	zapConfig.EncoderConfig.EncodeCaller = ColorizeCallerEncoder
	logger, _ := zapConfig.Build()
	return &ZapAdapter{
		zaplog: logger,
	}
}

func NewZapAdapter(zaplog *zap.Logger) *ZapAdapter {
	return &ZapAdapter{
		zaplog: zaplog,
	}
}

type ZapAdapter struct {
	zaplog *zap.Logger
}

// Err implement logger interface
func (l *ZapAdapter) Err(err error) Logger {
	var (
		fields = make([]zap.Field, 0, 1)
	)

	fields = append(fields, zap.String("error", fmt.Sprintf("%+v", err)))

	zaplog := l.zaplog.With(fields...)
	return l.clone(zaplog)
}

// Field implement logger interface
func (l *ZapAdapter) Field(key string, val interface{}) Logger {
	field := getField(key, val)
	zaplog := l.zaplog.With(field)
	return l.clone(zaplog)
}

// Fields implement logger interface
func (l *ZapAdapter) Fields(fieldsMap map[string]interface{}) Logger {
	fields := make([]zap.Field, 0, len(fieldsMap))
	for k, v := range fieldsMap {
		fields = append(fields, getField(k, v))
	}
	zaplog := l.zaplog.With(fields...)
	return l.clone(zaplog)
}

// Info implement logger interface
func (l *ZapAdapter) Info(msg string) {
	l.zaplog.Info(msg)
	defer l.zaplog.Sync()
}

// Debug implement logger interface
func (l *ZapAdapter) Debug(msg string) {
	l.zaplog.Debug(msg)
	defer l.zaplog.Sync()

}

// Warn implement logger interface
func (l *ZapAdapter) Warn(msg string) {
	l.zaplog.Warn(msg)
	defer l.zaplog.Sync()

}

// Error implement logger interface
func (l *ZapAdapter) Error(msg string) {
	l.zaplog.Error(msg)
	defer l.zaplog.Sync()

}

// Fatal implement logger interface
func (l *ZapAdapter) Fatal(msg string) {
	l.zaplog.Fatal(msg)
	defer l.zaplog.Sync()

}

func (l *ZapAdapter) Debugf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Debugf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Infof(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Infof(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Warnf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Warnf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Errorf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Errorf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Fatalf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Fatalf(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Printf(msg string, args ...interface{}) {
	s := l.zaplog.Sugar()
	s.Infof(msg, args)
	defer s.Sync()
}

func (l *ZapAdapter) Attach(ctx context.Context) context.Context {
	return Attach(ctx, l)
}

func (l *ZapAdapter) Caller(stack int) Logger {
	return l.clone(l.zaplog.WithOptions(zap.AddCallerSkip(stack)))
}

func (l *ZapAdapter) clone(zaplog *zap.Logger) *ZapAdapter {
	return &ZapAdapter{
		zaplog: zaplog,
	}
}

func getField(key string, val interface{}) zap.Field {
	switch v := val.(type) {
	case int:
		return zap.Int(key, v)
	case int64:
		return zap.Int64(key, v)
	case int32:
		return zap.Int32(key, v)
	case string:
		return zap.String(key, v)
	case time.Time:
		return zap.Time(key, v)
	case fmt.Stringer:
		return zap.Stringer(key, v)
	case bool:
		return zap.Bool(key, v)
	case float32:
		return zap.Float32(key, v)
	case float64:
		return zap.Float64(key, v)
	case []byte:
		return zap.Binary(key, v)
	}
	return zap.Skip()
}
