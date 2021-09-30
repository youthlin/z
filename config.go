package z

import (
	"encoding/json"
	"io"
	"os"

	"github.com/youthlin/logs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// debugConsole zap 配置：控制台打印 debug 日志
var debugConsole = &ZapConfig{
	Name:   "default-debug-console",
	Enable: true,
	Level:  zap.DebugLevel,
	Encoder: &Encoder{
		AsJSON:        false,
		EncoderConfig: zap.NewDevelopmentEncoderConfig(),
		EncodeTime:    TimeEncoder{Name: "ISO8601"},
	},
	Output: &Output{
		Type: Console,
		File: lumberjack.Logger{Filename: Stdout},
	},
}

// DefaultConfig 返回默认的日志配置，debug 级别、输出到控制台
func DefaultConfig() *LogsConfig {
	return &LogsConfig{
		Level: logs.LevelConfig(logs.LevelDebug),
		Zap:   []*ZapConfig{debugConsole},
	}
}

// LogsConfig 日志配置
type LogsConfig struct {
	Level *logs.LoggerLevel `json:"level" yaml:"level"` // 各 Logger 的日志级别,配合 logs 库使用
	Zap   []*ZapConfig      `json:"zap" yaml:"zap"`     // zap 的输出配置
}

// ZapConfig 日志输出配置
type ZapConfig struct {
	Name    string        `json:"name" yaml:"name"`       // 配置名称
	Enable  bool          `json:"enable" yaml:"enable"`   // 是否启用
	Level   zapcore.Level `json:"level" yaml:"level"`     // 大于等于该级别的日志才会输出
	Encoder *Encoder      `json:"encoder" yaml:"encoder"` // 输出的各个字段格式
	Output  *Output       `json:"output" yaml:"output"`   // 日志输出去向
}

// Encoder 包装 zapcore 的 EncoderConfig
type Encoder struct {
	zapcore.EncoderConfig
	EncodeTime TimeEncoder `json:"timeEncoder" yaml:"timeEncoder"`
	AsJSON     bool        `json:"json" yaml:"json"` // 整条日志使用 JSON 格式输出
	encoder    zapcore.Encoder
}

// Zap 转为 zap 的 Encoder
func (e *Encoder) Zap() zapcore.Encoder {
	if e.encoder == nil {
		ze := zap.NewProductionEncoderConfig()          // 默认配置
		copyNoneZeroField(&e.EncoderConfig, &ze)        // 自定义配置覆盖默认配置
		ze.EncodeTime = e.EncodeTime.ToZapTimeEncoder() // 时间格式特殊处理
		e.EncoderConfig = ze                            // 保存最终配置
		var encoder zapcore.Encoder
		if e.AsJSON { // 输出为 json
			encoder = zapcore.NewJSONEncoder(ze)
		} else {
			encoder = zapcore.NewConsoleEncoder(ze)
		}
		e.encoder = encoder
	}
	return e.encoder
}

// AsMap 转为 map 用于序列化为 yaml/json.
// zap 配置本身只支持反序列化，不支持序列化，所以需要先转 map
func (e *Encoder) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"json":          e.AsJSON,
		"messageKey":    e.MessageKey,
		"levelKey":      e.LevelKey,
		"timeKey":       e.TimeKey,
		"nameKey":       e.NameKey,
		"callerKey":     e.CallerKey,
		"functionKey":   e.FunctionKey,
		"stacktraceKey": e.StacktraceKey,
		"lineEnding":    e.LineEnding,
		"levelEncoder":  marshalName(funName(e.EncodeLevel)),
		"timeEncoder": func() interface{} {
			// 特殊处理两种格式
			if layout := e.EncodeTime.Layout; layout != "" {
				return map[string]string{
					"layout": layout,
				}
			}
			return e.EncodeTime.Name
		}(),
		"durationEncoder":  marshalName(funName(e.EncodeDuration)),
		"callerEncoder":    marshalName(funName(e.EncodeCaller)),
		"nameEncoder":      marshalName(funName(e.EncodeName)),
		"consoleSeparator": e.ConsoleSeparator,
	}
}
func (e Encoder) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.AsMap())
}
func (e Encoder) MarshalYAML() (interface{}, error) {
	return e.AsMap(), nil
}

// UnmarshalYAML 先反序列化到 map，然后通过 json 反序列化
// 因为 yaml 不支持反序列化嵌入结构（需要 `yaml:"",inline"` 不能有同名字段）
func (e *Encoder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var m = make(map[string]interface{})
	if err := unmarshal(&m); err != nil {
		return err
	}
	te := m["timeEncoder"]
	if layoutTe, ok := te.(map[interface{}]interface{}); ok {
		// [key 0] -> "layout"
		// [val 1] -> "2006-01-02"
		for _, v := range layoutTe {
			if layout, ok := v.(string); ok {
				if layout != "layout" {
					m["timeEncoder"] = map[string]interface{}{
						"layout": layout,
					}
				}
			}
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, e)
}

var funcNameMap = map[string]string{
	"<nil>": "",
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

// TimeEncoder zapcore 的 TimeEncoder 可能有两种格式
// 1. ISO8601 这种指定的预置格式 { "timeEncoder": "ISO8601" }
// 2. 用户使用 layout 字段指定时间格式 { "timeEncoder: { "layout": "2006-01-02 15:04:05.000"} }
type TimeEncoder struct {
	Name   string
	Layout string
}

func (e *TimeEncoder) UnmarshalText(text []byte) error {
	e.Name = string(text)
	return nil
}

// UnmarshalYAML 反序列化 yaml，先尝试 layout 模式，再尝试预置格式
func (e *TimeEncoder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var o struct {
		Layout string `json:"layout" yaml:"layout"`
	}
	if err := unmarshal(&o); err == nil {
		e.Layout = o.Layout
		return nil
	}
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	return e.UnmarshalText([]byte(s))
}
func (e *TimeEncoder) UnmarshalJSON(data []byte) error {
	return e.UnmarshalYAML(func(v interface{}) error {
		return json.Unmarshal(data, v)
	})
}

func (e *TimeEncoder) ToZapTimeEncoder() zapcore.TimeEncoder {
	if e.Layout != "" {
		return zapcore.TimeEncoderOfLayout(e.Layout)
	}
	var zte zapcore.TimeEncoder
	zte.UnmarshalText([]byte(e.Name))
	return zte
}

// Output 日志输出配置
type Output struct {
	Type    OutputType        `json:"type" yaml:"type"` // 控制台或文件
	File    lumberjack.Logger `json:"file" yaml:"file"` // 如果是控制台，则只需要填写 Filename(stdout/stderr)，如果是文件，则根据需要填写字段
	writerS zapcore.WriteSyncer
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

// WriteSyncer 转为 zap 的 WriteSyncer
func (o *Output) WriteSyncer() zapcore.WriteSyncer {
	if o.writerS == nil {
		var ws zapcore.WriteSyncer
		switch o.Type {
		case Console:
			if o.File.Filename == Stderr {
				ws = zapcore.AddSync(os.Stderr)
			} else {
				o.File.Filename = Stdout
				ws = zapcore.AddSync(os.Stdout)
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
			copyNoneZeroField(&o.File, fileOut) // 覆盖默认配置
			o.File = lumberjack.Logger{
				Filename:   fileOut.Filename,
				MaxSize:    fileOut.MaxSize,
				MaxAge:     fileOut.MaxAge,
				MaxBackups: fileOut.MaxBackups,
				LocalTime:  fileOut.LocalTime,
				Compress:   fileOut.Compress,
			}
			ws = zapcore.AddSync(fileOut)
		default: // 默认输出到控制台
			ws = zapcore.AddSync(os.Stdout)
		}
		o.writerS = ws
	}
	return o.writerS
}

var _ io.Writer = (*Output)(nil)

func (o *Output) Write(p []byte) (n int, err error) {
	return o.WriteSyncer().Write(p)
}
