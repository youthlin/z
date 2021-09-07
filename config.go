package z

import (
	"encoding/json"

	"github.com/youthlin/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var debugConsole = &Config{
	Name:          "default-debug-console",
	Enable:        true,
	Level:         zap.DebugLevel,
	AsJSON:        false,
	Output:        Output{Console, lumberjack.Logger{Filename: Stdout}},
	EncoderConfig: zap.NewDevelopmentEncoderConfig(),
}

// DefaultConfig 返回默认的日志配置，debug 级别、输出到控制台
func DefaultConfig() *LogsConfig {
	return &LogsConfig{
		Level: logs.LevelConfig(logs.LevelDebug),
		Zap:   []*Config{debugConsole},
	}
}

// LogsConfig 日志配置
type LogsConfig struct {
	// 各 Logger 的日志级别
	Level *logs.Config `json:"level" yaml:"level"`
	// zap 的输出配置
	Zap []*Config `json:"zap" yaml:"zap"`
}

// Config 日志输出配置
type Config struct {
	Name          string                `json:"name" yaml:"name"`                   // 配置名称
	Enable        bool                  `json:"enable" yaml:"enable"`               // 是否启用
	Level         zapcore.Level         `json:"level" yaml:"level"`                 // 大于等于该级别的日志才会输出
	AsJSON        bool                  `json:"json" yaml:"json"`                   // 整条日志使用 JSON 格式输出
	Output        Output                `json:"output" yaml:"output"`               // 日志输出去向
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"` // 输出的各个字段格式
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

func (c *Config) AsMap() map[string]interface{} {
	if c == nil {
		return nil
	}
	return map[string]interface{}{
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
			"levelEncoder":     marshalName(funName(c.EncoderConfig.EncodeLevel)),
			"timeEncoder":      marshalName(funName(c.EncoderConfig.EncodeTime)),
			"durationEncoder":  marshalName(funName(c.EncoderConfig.EncodeDuration)),
			"callerEncoder":    marshalName(funName(c.EncoderConfig.EncodeCaller)),
			"nameEncoder":      marshalName(funName(c.EncoderConfig.EncodeName)),
			"consoleSeparator": c.EncoderConfig.ConsoleSeparator,
		},
	}
}

var funcNameMap = map[string]string{
	"go.uber.org/zap/zapcore.CapitalLevelEncoder":        "capital",
	"go.uber.org/zap/zapcore.CapitalColorLevelEncoder":   "capitalColor",
	"go.uber.org/zap/zapcore.LowercaseColorLevelEncoder": "color",
	"go.uber.org/zap/zapcore.LowercaseLevelEncoder":      "",

	"go.uber.org/zap/zapcore.RFC3339NanoTimeEncoder": "rfc3339nano",
	"go.uber.org/zap/zapcore.RFC3339TimeEncoder":     "rfc3339",
	"go.uber.org/zap/zapcore.ISO8601TimeEncoder":     "iso8601",
	"go.uber.org/zap/zapcore.EpochMillisTimeEncoder": "millis",
	"go.uber.org/zap/zapcore.EpochNanosTimeEncoder":  "nanos",
	"go.uber.org/zap/zapcore.EpochTimeEncoder":       "",

	"go.uber.org/zap/zapcore.StringDurationEncoder":  "string",
	"go.uber.org/zap/zapcore.NanosDurationEncoder":   "nanos",
	"go.uber.org/zap/zapcore.MillisDurationEncoder":  "ms",
	"go.uber.org/zap/zapcore.SecondsDurationEncoder": "",

	"go.uber.org/zap/zapcore.FullCallerEncoder":  "full",
	"go.uber.org/zap/zapcore.ShortCallerEncoder": "",

	"go.uber.org/zap/zapcore.FullNameEncoder": "full",
}

func marshalName(funcName string) string {
	if name, ok := funcNameMap[funcName]; ok {
		return name
	}
	return funcName
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.AsMap())
}
func (c *Config) MarshalYAML() (interface{}, error) {
	return c.AsMap(), nil
}
