package z

import (
	"github.com/youthlin/logs"
	"github.com/youthlin/logs/pkg/kv"
	"go.uber.org/zap"
)

func NewZapAdaptor(zLog *zap.Logger) logs.Adaptor {
	// [3]logs.Debug -> [2]logs.Log() -> [1]ZapAdaptor.Log -> [0]sugar.Debug
	zLog = zLog.WithOptions(zap.AddCallerSkip(3))
	return &ZapAdaptor{zLog}
}

// ZapAdaptor is an adaptor for github.com/youthlin/logs
type ZapAdaptor struct {
	*zap.Logger
}

// Log log message to zap
func (s *ZapAdaptor) Log(msg logs.Message) {
	log := s.Logger
	if skip := msg.Skip(); skip != 0 {
		log = log.WithOptions(zap.AddCallerSkip(skip))
	}
	sugar := log.Sugar()
	kvs := append(kv.Get(msg.Ctx()), msg.Kvs()...)
	// kvs = append(kvs, "", msg.LoggerName())
	if len(kvs) > 1 {
		sugar = sugar.With(kvs...)
	}
	switch msg.Level() {
	case logs.LevelUnset:
		logs.Assert(false, "message level must set")
	case logs.LevelAll:
		fallthrough
	case logs.LevelTrace:
		fallthrough
	case logs.LevelDebug:
		sugar.Debugf(msg.Msg(), msg.Args()...)
	case logs.LevelInfo:
		sugar.Infof(msg.Msg(), msg.Args()...)
	case logs.LevelWarn:
		sugar.Warnf(msg.Msg(), msg.Args()...)
	case logs.LevelError:
		sugar.Errorf(msg.Msg(), msg.Args()...)
	case logs.LevelNone:
		return
	}
}
