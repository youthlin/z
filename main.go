package z

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger new 一个 Logger
func NewLogger(configs []*Config) Logger {
	var core []zapcore.Core
	for i := range configs {
		config := configs[i]
		if config != nil && config.Enable {
			// zapcore.Encoder + zapcore.WriteSyncer => zapcore.Core
			core = append(core, zapcore.NewCore(buildEncoder(config), buildOut(&config.Output), config.Level))
		}
	}
	zapLogger := zap.New(zapcore.NewTee(core...), zap.AddCaller())
	zap.ReplaceGlobals(zapLogger)
	return newLogger(zapLogger)
}

// buildEncoder 根据配置设置输出格式
func buildEncoder(config *Config) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()        // 默认配置
	copyNoneZeroField(&config.EncoderConfig, &encoderConfig) // 覆盖默认配置
	var encoder zapcore.Encoder
	if config.AsJSON { // 输出为 json
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	return encoder
}

// buildOut 配置日志目的地
func buildOut(output *Output) zapcore.WriteSyncer {
	switch output.Type {
	case Console:
		if output.Destination.Filename == Stderr {
			return zapcore.AddSync(os.Stderr)
		} else {
			output.Destination.Filename = Stdout
			return zapcore.AddSync(os.Stdout)
		}
	case File:
		var fileOut = &lumberjack.Logger{ // 日志切割: 默认配置
			Filename:   "app.log", // 文件名
			MaxSize:    100,       // MB 超过这个大小会切割日志
			MaxAge:     30,        // day 切割的日志最多保存几天
			MaxBackups: 30,        // 切割的日志最多最多保存几个
			LocalTime:  false,     // 默认 false=UTC 时间
			Compress:   true,      // 压缩
		}
		copyNoneZeroField(&output.Destination, fileOut) // 覆盖默认配置
		return zapcore.AddSync(fileOut)
	default: // 默认输出到控制台
		return zapcore.AddSync(os.Stdout)
	}
}
