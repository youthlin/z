# z

- 对 [zap](https://github.com/uber-go/zap) 简单封装了一层
- 支持 github.com/youthlin/logs 日志接口 Adaptor 到 zap

使用方式

```go
// ------- use default config -------
zap.L().Info("info message")
// 2021-09-07T12:58:14.413+0800	INFO	z/example_test.go:79	info message

zap.S().Infow("sugar info", "key", "value")
// 2021-09-07T12:59:58.649+0800	INFO	z/example_test.go:80	sugar info	{"key": "value"}

// ------- use custom config -------
log := z.NewLogger(z.DefaultConfig().Zap)
log.Debug("Debug message")
// 2021-09-07T12:59:58.649+0800	DEBUG	z/example_test.go:85	Debug message

log.With(zap.String("key", "value")).Info("info message")
// 2021-09-07T12:59:58.649+0800	INFO	z/example_test.go:86	info message	{"key": "value"}

sugar := log.Sugar()
sugar.Debugf("Hello, %s, debug message", "tom")
// 2021-09-07T12:59:58.649+0800	DEBUG	z/example_test.go:88	Hello, tom, debug message

sugar.Infow("sugar info with fields", "key", "value", "int", 42)
// 2021-09-07T12:59:58.649+0800	INFO	z/example_test.go:89	sugar info with fields	{"key": "value", "int": 42}

// ------- as logs Adaptor -------
// import "github.com/youthlin/logs"
var log = logs.GetLogger()
log.Ctx(ctx).With("key", 42).Debug("Hello %s", "Tom")

// use custom config
func init() {
	c := DefaultConfig()// or from .yaml / .json
  z.SetConfig(c)
}
var log = logs.GetLogger()
log.Info("Hello %s", "world")
```

`Configs` 也支持 从 yaml/json 读取。 see `testdata/config.yaml` and `testdata/config.json`

```yaml
logs:
  level: # logs config
    root: error # default only print error log
    loggers:
      "github.com": warn # if package name is github.com print warn+ log
      "github.com/youthlin": debug # this package print debug+ log
  zap: # zap config
    - name: console # 控制台输出所有日志
      enable: true
      json: false # 不需要格式化为 json
      level: info
      output:
        type: console
        destination:
          filename: stdout
      encoderConfig:
        levelEncoder: capitalColor # 带颜色大写的日志级别 capital/capitalColor/color/lowcase
        timeEncoder: rfc3339nano # e.g.: 2006-01-02T15:04:05.999999999Z07:00 rfc3339nano/rfc3339/iso8601/millis/nanos/epoch
        durationEncoder: string # 时间段格式化为带单位的: 968.6µs string/nanos/ms/seconds or 带 layout 子字段
        callerEncoder: full # full/short
 ```
