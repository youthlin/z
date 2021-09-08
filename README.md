# z
[![sync-to-gitee](https://github.com/youthlin/z/actions/workflows/gitee.yaml/badge.svg)](https://github.com/youthlin/z/actions/workflows/gitee.yaml)
[![test](https://github.com/youthlin/z/actions/workflows/test.yaml/badge.svg)](https://github.com/youthlin/z/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/youthlin/z/branch/main/graph/badge.svg?token=7Y3JFVCDJQ)](https://codecov.io/gh/youthlin/z)
[![Go Report Card](https://goreportcard.com/badge/github.com/youthlin/z)](https://goreportcard.com/report/github.com/youthlin/z)
[![Go Reference](https://pkg.go.dev/badge/github.com/youthlin/z.svg)](https://pkg.go.dev/github.com/youthlin/z)

- 快速设置 [zap](https://github.com/uber-go/zap)，支持读取与序列化 zap 配置
- 支持 github.com/youthlin/logs 日志接口 Adaptor 到 zap

## import
```shell
go get -u github.com/youthlin/z
# 国内镜像
go mod edit -replace github.com/youthlin/z@latest=gitee.com/youthlin/logz@latest&&go mod tidy
```
> gitee 镜像：[gitee.com/youthlin/logz](https://gitee.com/youthlin/logz) (logz = logs + zap)
>
> 鸣谢 仓库同步工具 https://github.com/Yikun/hub-mirror-action


## Usage 使用方式

```go
// to init zap config
import _ "github.com/youthlin/z"

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

// ------- ------- ------- -------
// ------- as logs Adaptor -------
// ------- ------- ------- -------
// import "github.com/youthlin/logs"
logs.Ctx(ctx).With("key", 42).Debug("Hello %s", "Tom")

// use custom config
func init() {
	c := DefaultConfig()// or from .yaml / .json
  z.SetConfig(c)
}
logs.Info("Hello %s", "world")
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
      level: info
      output:
        type: console
        destination:
          filename: stdout
      encoder:
        json: false # 不需要格式化为 json
        levelEncoder: capitalColor # 带颜色大写的日志级别 capital/capitalColor/color/lowcase
        timeEncoder: rfc3339nano # e.g.: 2006-01-02T15:04:05.999999999Z07:00 rfc3339nano/rfc3339/iso8601/millis/nanos/epoch
        durationEncoder: string # 时间段格式化为带单位的: 968.6µs string/nanos/ms/seconds or 带 layout 子字段
        callerEncoder: full # full/short
 ```
