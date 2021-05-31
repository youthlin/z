package z

import "go.uber.org/zap/zapcore"

const nopLogger = nop(0)

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
