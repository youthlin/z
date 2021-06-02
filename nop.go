package z

import "go.uber.org/zap/zapcore"

// nopLogger 一个没有任何功能的 Logger 实例
const nopLogger = nop(0)

// nop 没有任何功能的 Logger
type nop int

func (n nop) Flush()                              {}
func (n nop) With(...interface{}) Logger          { return nopLogger }
func (n nop) WithSkip(int, ...interface{}) Logger { return nopLogger }
func (n nop) Enable(zapcore.Level) bool           { return false }
func (n nop) Debug(string, ...interface{})        {}
func (n nop) Info(string, ...interface{})         {}
func (n nop) Warn(string, ...interface{})         {}
func (n nop) Error(string, ...interface{})        {}
func (n nop) DebugJSON(string, ...interface{})    {}
func (n nop) InfoJSON(string, ...interface{})     {}
func (n nop) WarnJSON(string, ...interface{})     {}
func (n nop) ErrorJSON(string, ...interface{})    {}
