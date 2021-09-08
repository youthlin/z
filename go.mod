module github.com/youthlin/z

go 1.16

require (
	github.com/cockroachdb/errors v1.8.6
	github.com/youthlin/logs v0.0.5
	go.uber.org/zap v1.17.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/youthlin/logs v0.0.5 => gitee.com/youthlin/logs v0.0.5
