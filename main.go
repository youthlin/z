package z

import (
	"github.com/youthlin/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	c := DefaultConfig()
	SetConfig(c)
}

func SetConfig(c *LogsConfig) {
	zapLogger := NewLogger(c.Zap)
	zap.ReplaceGlobals(zapLogger)

	logs.SetLoggerLevel(c.Level)
	logs.SetAdaptor(NewZapAdaptor(zapLogger))
}

// NewLogger new 一个 Logger
func NewLogger(configs []*ZapConfig) *zap.Logger {
	var core []zapcore.Core
	for i := range configs {
		config := configs[i]
		if config != nil && config.Enable {
			// zapcore.Encoder + zapcore.WriteSyncer => zapcore.Core
			core = append(core, zapcore.NewCore(config.Encoder.Zap(), config.Output.WriteSyncer(), config.Level))
		}
	}
	return zap.New(zapcore.NewTee(core...), zap.AddCaller())
}
