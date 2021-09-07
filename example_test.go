package z_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/youthlin/logs"
	"github.com/youthlin/logs/pkg/arg"
	"github.com/youthlin/logs/pkg/callinfo"
	"github.com/youthlin/logs/pkg/kv"
	"github.com/youthlin/z"
	"go.uber.org/zap"
)

type some struct {
	Name string
}

var s = some{"SomeName"}
var ctx = context.Background()
var err = errors.New("error")

func TestExample(t *testing.T) {
	// import _ "github.com/youthlin/z" // init zap Adaptor
	log := logs.GetLogger()
	{
		// logs
		log.Trace("Hello, Trace")
		log.Debug("Hello, World")
		log.Info("Hello, World|s=%v", s)
		log.Warn("Hello, World|s=%#v", &s)
		log.Error("Hello, World|s=%T", s)
		log.Debug("Hello s=%v", arg.JSON(s))
		log.Info("Hello s=%v", arg.JSON(s))
		log.Warn("Hello s=%v", arg.JSON(s))
		log.Error("Hello s=%v", arg.JSON(s))
		log.Error("Hello s=%v|err=%v", s, err)
		log.Error("Hello s=%v|err=%v", s, arg.ErrJSON("%+v", err))
	}
	{
		// with
		log.With("key", "value").Info("With key=value")
	}
	{
		// ctx
		log.Ctx(ctx).Debug("No KV")
		ctx = kv.Add(ctx, "p1", "p1", "p0") // 不是偶数个，忽略
		ctx = kv.Add(ctx, "k1", "v1", "k0", "v0")
		log.Ctx(ctx).Debug("with k1-v1")
		log.Ctx(ctx).Info("with k1-v1")
		log.Ctx(ctx).Warn("with k1-v1")
		log.Ctx(ctx).Error("with k1-v1")
		ctx = kv.Add(ctx, "k2", "v2")
		log.Ctx(ctx).Debug("with k1-v1 and k2-v2")
		log.Ctx(ctx).Debug("with k1-v1 and k2-v2|s=%v", arg.JSON(s))
		log.Ctx(ctx).Info("with k1-v1 and k2-v2|s=%v", arg.JSON(s))
		log.Ctx(ctx).Warn("with k1-v1 and k2-v2|s=%v", arg.JSON(s))
		log.Ctx(ctx).Error("with k1-v1 and k2-v2|s=%v", arg.JSON(s))
		type k string
		ctx = context.WithValue(ctx, k("中间key"), "some...")
		ctx = kv.Add(ctx, "k3", "v3")
		log.Ctx(ctx).Debug("with k1-v1 and k2-v2 k3-v3")
	}
	{
		logs.SetConfig(&logs.Config{
			Root: logs.LevelWarn,
			Loggers: map[string]logs.Level{
				"github.com": logs.LevelDebug,
			},
		})
		log.Debug("my package is %s", callinfo.Get().PkgName)
		abcLog := logs.GetLogger(logs.WithName("abc"))
		abcLog.Debug("my logger name is abc, debug log not print")
		abcLog.Warn("my logger name is abc, warn log would print")
	}
	{
		// zap global method
		zap.L().Info("info message")
		zap.S().Infow("sugar info", "key", "value")
	}
	{
		// zap interface
		log := z.NewLogger(z.DefaultConfig().Zap)
		log.Debug("Debug message")
		log.With(zap.String("key", "value")).Info("info message")
		sugar := log.Sugar()
		sugar.Debugf("Hello, %s, debug message", "tom")
		sugar.Infow("sugar info with fields", "key", "value", "int", 42)
	}
}
