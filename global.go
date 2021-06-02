package z

import "go.uber.org/zap/zapcore"

var globalLogger Logger = nopLogger

// SetGlobalLogger 设置全局 logger
func SetGlobalLogger(logger Logger) {
	// 使用 globalLogger 的场景都是本文件，对 Logger 又包装了一层函数，所以还需要往上跳过一层堆栈
	globalLogger = logger.WithSkip(1)
}

// Flush 刷新日志到输出
func Flush() {
	globalLogger.Flush()
}

// With 返回一个 Logger, 打印的日志会带上指定的 k-v
func With(kvs ...interface{}) Logger {
	// 这里返回的是 Logger 后续是调用 Logger 的方法，而不是本文件的包装函数，所以不需要跳过 SetGlobalLogger 时多跳过的一层
	return globalLogger.WithSkip(-1, kvs...)
}

// WithSkip 同 With 但打印代码行号时多跳过 skip 层调用
func WithSkip(skip int, kvs ...interface{}) Logger {
	return globalLogger.WithSkip(skip, kvs...)
}

// Enable 当前 Logger 是否会打印指定的日志级别
func Enable(level zapcore.Level) bool {
	return globalLogger.Enable(level)
}

// Debug 使用全局 Logger 打印 Debug 日志
func Debug(fmt string, args ...interface{}) {
	globalLogger.Debug(fmt, args...)
}

// Info 使用全局 Logger 打印 Info 日志
func Info(fmt string, args ...interface{}) {
	globalLogger.Info(fmt, args...)
}

// Warn 使用全局 Logger 打印 Warn 日志
func Warn(fmt string, args ...interface{}) {
	globalLogger.Warn(fmt, args...)
}

// Error 使用全局 Logger 打印 Error 日志
func Error(fmt string, args ...interface{}) {
	globalLogger.Error(fmt, args...)
}

// DebugJSON 使用全局 Logger 打印 Debug 日志，参数会转为 JSON 格式
func DebugJSON(fmt string, args ...interface{}) {
	globalLogger.DebugJSON(fmt, args...)
}

// InfoJSON 使用全局 Logger 打印 Info 日志，参数会转为 JSON 格式
func InfoJSON(fmt string, args ...interface{}) {
	globalLogger.InfoJSON(fmt, args...)
}

// WarnJSON 使用全局 Logger 打印 Warn 日志，参数会转为 JSON 格式
func WarnJSON(fmt string, args ...interface{}) {
	globalLogger.WarnJSON(fmt, args...)
}

// ErrorJSON 使用全局 Logger 打印 Error 日志，参数会转为 JSON 格式
func ErrorJSON(fmt string, args ...interface{}) {
	globalLogger.ErrorJSON(fmt, args...)
}
