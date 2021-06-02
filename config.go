package z

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var infoConsole = &Config{
	Name:          "default-info-console",
	Enable:        true,
	Level:         zap.DebugLevel,
	AsJSON:        false,
	Output:        Output{Console, lumberjack.Logger{Filename: Stdout}},
	EncoderConfig: zap.NewDevelopmentEncoderConfig(),
}

// DefaultConfig 返回默认的日志配置，debug 级别、输出到控制台
func DefaultConfig() []*Config {
	return []*Config{infoConsole}
}

// Output 日志输出配置
type Output struct {
	Type        OutputType        `json:"type" yaml:"type"`               // 控制台或文件
	Destination lumberjack.Logger `json:"destination" yaml:"destination"` // 如果是控制台，则只需要填写 Filename(stdout/stderr)，如果是文件，则根据需要填写字段
}

// OutputType 日志输出类型
type OutputType string

const (
	Console OutputType = "console" // 控制台
	File    OutputType = "file"    // 文件
)
const (
	Stdout = "stdout" // 标准输出
	Stderr = "stderr" // 标准错误
)

// Config 日志输出配置
type Config struct {
	Name          string                `json:"name" yaml:"name"`                   // 配置名称
	Enable        bool                  `json:"enable" yaml:"enable"`               // 是否启用
	Level         zapcore.Level         `json:"level" yaml:"level"`                 // 大于等于该级别的日志才会输出
	AsJSON        bool                  `json:"json" yaml:"json"`                   // 整条日志使用 JSON 格式输出
	Output        Output                `json:"output" yaml:"output"`               // 输出配置
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"` // 输出的各个字段配置
}

// Configs 日志配置数组
type Configs []*Config

func (c *Config) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}
	m := map[string]interface{}{
		"name":   c.Name,
		"enable": c.Enable,
		"level":  c.Level,
		"json":   c.AsJSON,
		"output": map[string]interface{}{
			"type": c.Output.Type,
			"destination": map[string]interface{}{
				"filename":   c.Output.Destination.Filename,
				"maxsize":    c.Output.Destination.MaxSize,
				"maxage":     c.Output.Destination.MaxAge,
				"maxbackups": c.Output.Destination.MaxBackups,
				"localtime":  c.Output.Destination.LocalTime,
				"compress":   c.Output.Destination.Compress,
			},
		},
		"encoderConfig": map[string]interface{}{
			"messageKey":       c.EncoderConfig.MessageKey,
			"levelKey":         c.EncoderConfig.LevelKey,
			"timeKey":          c.EncoderConfig.TimeKey,
			"nameKey":          c.EncoderConfig.NameKey,
			"callerKey":        c.EncoderConfig.CallerKey,
			"functionKey":      c.EncoderConfig.FunctionKey,
			"stacktraceKey":    c.EncoderConfig.StacktraceKey,
			"lineEnding":       c.EncoderConfig.LineEnding,
			"levelEncoder":     funName(c.EncoderConfig.EncodeLevel),
			"timeEncoder":      funName(c.EncoderConfig.EncodeTime),
			"durationEncoder":  funName(c.EncoderConfig.EncodeDuration),
			"callerEncoder":    funName(c.EncoderConfig.EncodeCaller),
			"nameEncoder":      funName(c.EncoderConfig.EncodeName),
			"consoleSeparator": c.EncoderConfig.ConsoleSeparator,
		},
	}
	return json.Marshal(m)
}
