package auth

import "context"

type ctxKey int

const principalKey ctxKey = 1

// WithPrincipal 注入主体。
func WithPrincipal(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, principalKey, name)
}

// Principal 读取主体。
func Principal(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(principalKey).(string)
	return v, ok
}
