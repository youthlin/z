package z_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/youthlin/z"
)

type some struct {
	Name string
}

var s = some{"SomeName"}
var ctx = context.Background()
var err = errors.New("error")

func TestExample(t *testing.T) {
	logger := z.NewLogger(z.DefaultConfig())
	z.SetGlobalLogger(logger)
	Convey("Example", t, func() {
		Convey("logger", func() { // 直接使用 NewLogger 实例，关注行号是否正确
			logger.Debug("Logger Debug")
			logger.InfoJSON("config=%v", z.DefaultConfig())
			logger.Info("config=%#v", z.DefaultConfig())
			logger.With("logger", "name").Warn("see line number")
		})
		Convey("Basic", func() { // 使用全局 Logger 关注行号是否正确
			z.Debug("Hello, World")
			z.Info("Hello, World|s=%v", s)
			z.Warn("Hello, World|s=%#v", &s)
			z.Error("Hello, World|s=%T", s)
			z.DebugJSON("Hello s=%v", s)
			z.InfoJSON("Hello s=%v", s)
			z.WarnJSON("Hello s=%v", s)
			z.ErrorJSON("Hello s=%v", s)
			z.ErrorJSON("Hello s=%v|err=%v", s, err)
			z.ErrorJSON("Hello s=%v|err=%v", s, z.Err("%+v", err))
		})
		Convey("With", func() {
			z.With("key", "value").Info("With key=value")
			z.WithSkip(-2, "-2 key", "value").Info("With key=value")
			z.WithSkip(-1, "-1 key", "value").Info("With key=value")
			z.WithSkip(0, "0  key", "value").Info("With key=value")
			z.WithSkip(1, "1  key", "value").Info("With key=value")
		})
		Convey("Ctx", func() {
			z.CtxDebug(ctx, "No KV")
			ctx = z.CtxAddKV(ctx, "p1", "p1", "p0") // 不是偶数个，忽略
			ctx = z.CtxAddKV(ctx, "k1", "v1", "k0", "v0")
			z.CtxDebug(ctx, "with k1-v1")
			z.CtxInfo(ctx, "with k1-v1")
			z.CtxWarn(ctx, "with k1-v1")
			z.CtxError(ctx, "with k1-v1")
			ctx = z.CtxAddKV(ctx, "k2", "v2")
			z.CtxDebug(ctx, "with k1-v1 and k2-v2")
			z.CtxDebugJSON(ctx, "with k1-v1 and k2-v2|s=%v", s)
			z.CtxInfoJSON(ctx, "with k1-v1 and k2-v2|s=%v", s)
			z.CtxWarnJSON(ctx, "with k1-v1 and k2-v2|s=%v", s)
			z.CtxErrorJSON(ctx, "with k1-v1 and k2-v2|s=%v", s)
			type k string
			ctx = context.WithValue(ctx, k("中间key"), "some...")
			ctx = z.CtxAddKV(ctx, "k3", "v3")
			z.CtxDebug(ctx, "with k1-v1 and k2-v2 k3-v3")
		})

	})
}
