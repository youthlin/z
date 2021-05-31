package z

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// region interface

// Logger 分级别日志打印
type Logger interface {
	Flush()
	With(kvs ...interface{}) Logger
	WithSkip(skip int, kvs ...interface{}) Logger
	Enable(level zapcore.Level) bool

	Debug(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})

	DebugJSON(fmt string, args ...interface{})
	InfoJSON(fmt string, args ...interface{})
	WarnJSON(fmt string, args ...interface{})
	ErrorJSON(fmt string, args ...interface{})
}

// endregion interface

type logger struct {
	base   *zap.Logger
	logger *zap.SugaredLogger
}

// newLogger create a logger, skip: 相对于本结构体跳过几层调用
func newLogger(zapL *zap.Logger) *logger {
	zapL.WithOptions()
	// 对 SugaredLogger.Debugf 包装了一层，所以需要多往上跳一层调用才是用户的行号
	zapL = zapL.WithOptions(zap.AddCallerSkip(1))
	return &logger{
		base:   zapL,
		logger: zapL.Sugar(),
	}
}

func (l *logger) Flush() {
	_ = l.base.Sync()
}

func (l *logger) With(kvs ...interface{}) Logger {
	return l.WithSkip(0, kvs...)
}
func (l *logger) WithSkip(skip int, kvs ...interface{}) Logger {
	base := l.base.WithOptions(zap.AddCallerSkip(skip))
	sug := base.Sugar().With(kvs...)
	return &logger{
		base:   base,
		logger: sug,
	}
}

// regionLog

func (l *logger) Debug(fmt string, args ...interface{}) {
	l.logger.Debugf(fmt, args...)
}

func (l *logger) Info(fmt string, args ...interface{}) {
	l.logger.Infof(fmt, args...)
}

func (l *logger) Warn(fmt string, args ...interface{}) {
	l.logger.Warnf(fmt, args...)
}

func (l *logger) Error(fmt string, args ...interface{}) {
	l.logger.Errorf(fmt, args...)
}

func (l *logger) Enable(level zapcore.Level) bool {
	return l.base.Core().Enabled(level)
}

// endregion Log

// region  xxxJSON

func (l *logger) DebugJSON(fmt string, args ...interface{}) {
	if l.Enable(zapcore.DebugLevel) {
		l.logger.Debugf(fmt, toJSON(args...)...)
	}
}

func (l *logger) InfoJSON(fmt string, args ...interface{}) {
	if l.Enable(zapcore.InfoLevel) {
		l.logger.Infof(fmt, toJSON(args...)...)
	}
}

func (l *logger) WarnJSON(fmt string, args ...interface{}) {
	if l.Enable(zapcore.WarnLevel) {
		l.logger.Warnf(fmt, toJSON(args...)...)
	}
}

func (l *logger) ErrorJSON(fmt string, args ...interface{}) {
	if l.Enable(zapcore.ErrorLevel) {
		l.logger.Errorf(fmt, toJSON(args...)...)
	}
}

// endregion  xxxJSON
