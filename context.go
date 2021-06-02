package z

import (
	"context"

	"go.uber.org/zap/zapcore"
)

// ctxKey 新定义一个类型用在 context.WithValue
type ctxKey string

// ctxKeyKV 用于 context.WithValue 的 key 的实例
var ctxKeyKV = ctxKey("ctx_kvs")

// kv 用于日志打印的 key-value 字段。current 是本次添加的，pre 是之前添加的。
type kv struct {
	current []interface{}
	pre     *kv
}

// getKv 获取当前上下文的 kv 结构
func getKv(ctx context.Context) *kv {
	if ctx == nil {
		return nil
	}
	if result := ctx.Value(ctxKeyKV); result != nil {
		if result, ok := result.(*kv); ok {
			return result
		}
	}
	return nil
}

// getAllKvs 获取当前上下文及其所有祖先的 kv 结构
func getAllKvs(ctx context.Context) []interface{} {
	if ctx == nil {
		return nil
	}
	currentKV := getKv(ctx)
	if currentKV == nil {
		return nil
	}
	var result []interface{}
	p := currentKV
	for p != nil {
		result = append(p.current, result...)
		p = p.pre
	}
	return result
}

// CtxAddKV 为上下文注入 kv 字段
func CtxAddKV(ctx context.Context, kvs ...interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(kvs) == 0 || (len(kvs)&1 == 1) { // 忽律不是偶数的情况
		return ctx
	}
	current := make([]interface{}, 0, len(kvs))
	current = append(current, kvs...)
	return context.WithValue(ctx, ctxKeyKV, &kv{
		current: current,
		pre:     getKv(ctx),
	})
}

// CtxDebug 同 Debug, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxDebug(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.DebugLevel) {
		// WithSkip 后返回的是 Logger, 但是 Debug 是在 CtxDebug 中调用的，所以需要跳过的层数和 SetGlobalLogger 是一致的
		WithSkip(0, getAllKvs(ctx)...).Debug(fmt, args...)
	}
}

// CtxInfo 同 Info, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxInfo(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.InfoLevel) {
		WithSkip(0, getAllKvs(ctx)...).Info(fmt, args...)
	}
}

// CtxWarn 同 Warn, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxWarn(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.WarnLevel) {
		WithSkip(0, getAllKvs(ctx)...).Warn(fmt, args...)
	}
}

// CtxError 同 Error, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxError(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.ErrorLevel) {
		WithSkip(0, getAllKvs(ctx)...).Error(fmt, args...)
	}
}

// CtxDebugJSON 同 DebugJSON, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxDebugJSON(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.DebugLevel) {
		WithSkip(0, getAllKvs(ctx)...).DebugJSON(fmt, args...)
	}
}

// CtxInfoJSON 同 InfoJSON, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxInfoJSON(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.InfoLevel) {
		WithSkip(0, getAllKvs(ctx)...).InfoJSON(fmt, args...)
	}
}

// CtxWarnJSON 同 WarnJSON, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxWarnJSON(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.WarnLevel) {
		WithSkip(0, getAllKvs(ctx)...).WarnJSON(fmt, args...)
	}
}

// CtxErrorJSON 同 ErrorJSON, 带 ctx 参数，可将 ctx 上带的 k-v 打印，见 CtxAddKV
func CtxErrorJSON(ctx context.Context, fmt string, args ...interface{}) {
	if Enable(zapcore.ErrorLevel) {
		WithSkip(0, getAllKvs(ctx)...).ErrorJSON(fmt, args...)
	}
}
