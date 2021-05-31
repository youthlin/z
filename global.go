package z

import "go.uber.org/zap/zapcore"

var globalLogger Logger = nopLogger

func SetGlobalLogger(logger Logger) {
	// 使用 globalLogger 的场景都是本文件，对 Logger 又包装了一层函数，所以还需要往上跳过一层堆栈
	globalLogger = logger.WithSkip(1)
}
func With(kvs ...interface{}) Logger {
	// 这里返回的是 Logger 后续是调用 Logger 的方法，而不是本文件的包装函数，所以不需要跳过 SetGlobalLogger 时多跳过的一层
	return globalLogger.WithSkip(-1, kvs...)
}
func WithSkip(skip int, kvs ...interface{}) Logger {
	return globalLogger.WithSkip(skip, kvs...)
}

func Enable(level zapcore.Level) bool {
	return globalLogger.Enable(level)
}

func Debug(fmt string, args ...interface{}) {
	globalLogger.Debug(fmt, args...)
}
func Info(fmt string, args ...interface{}) {
	globalLogger.Info(fmt, args...)
}
func Warn(fmt string, args ...interface{}) {
	globalLogger.Warn(fmt, args...)
}
func Error(fmt string, args ...interface{}) {
	globalLogger.Error(fmt, args...)
}

func DebugJSON(fmt string, args ...interface{}) {
	globalLogger.DebugJSON(fmt, args...)
}
func InfoJSON(fmt string, args ...interface{}) {
	globalLogger.InfoJSON(fmt, args...)
}
func WarnJSON(fmt string, args ...interface{}) {
	globalLogger.WarnJSON(fmt, args...)
}
func ErrorJSON(fmt string, args ...interface{}) {
	globalLogger.ErrorJSON(fmt, args...)
}
