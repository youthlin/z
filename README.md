# z

对 [zap](https://github.com/uber-go/zap) 简单封装了一层

使用方式

```go
// use default config
logger := z.NewLogger(z.DefaultConfig())
logger.Info("format template|%v", "some var")

z.SetGlobalLogger(logger)
z.WarnJSON("something wrong|req=%v|err=%v", req, z.Err("%+v", err))

z.With(k, v, k2, v2).Info("msg")

ctx = z.CtxAddKV(ctx, "p1", "p1")
z.CtxInfo(ctx, "with context based key value")
```

`Configs` 也支持 从 yaml 读取

```yaml
logs:
  - name: console # 控制台输出所有日志
    enable: true
    json: false # 不需要格式化为 json
    level: info
    output:
      type: console
      destination:
        filename: stdout
    encoderConfig:
      levelEncoder: capitalColor # 带颜色大写的日志级别
      timeEncoder: rfc3339nano # e.g.: 2006-01-02T15:04:05.999999999Z07:00
      durationEncoder: string # 时间段格式化为带单位的: 968.6µs
 ```
