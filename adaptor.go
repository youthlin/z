package z

import (
	"github.com/youthlin/logs"
	"github.com/youthlin/logs/pkg/kv"
	"go.uber.org/zap"
)

type ZapAdaptor struct {
	*zap.Logger
}

func (s *ZapAdaptor) Log(msg logs.Message) {
	kvs := append(kv.Get(msg.Ctx()), msg.Kvs()...)
	// kvs = append(kvs, "", msg.LoggerName())
	log := s.Logger
	if msg.Skip() > 0 {
		log = log.WithOptions(zap.AddCallerSkip(msg.Skip()))
	}
	sugar := log.Sugar().With(kvs...)
	switch msg.Level() {
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
