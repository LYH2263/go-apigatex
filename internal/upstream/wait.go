package upstream

import (
	"context"
	"time"

	"github.com/LYH2263/go-apigatex/internal/errors"
)

// WaitReady 等待上游就绪信号；必须响应 ctx 取消。
func WaitReady(ctx context.Context, ready <-chan struct{}, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return errors.WrapErr(errors.ErrCanceled, ctx.Err())
	case <-ready:
		return nil
	case <-timer.C:
		return errors.ErrTimeout
	}
}

// WaitHealthy 轮询健康标记直到 true 或 ctx 取消。
func WaitHealthy(ctx context.Context, check func() bool, interval time.Duration) error {
	if interval <= 0 {
		interval = 20 * time.Millisecond
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		if check != nil && check() {
			return nil
		}
		select {
		case <-ctx.Done():
			return errors.WrapErr(errors.ErrCanceled, ctx.Err())
		case <-t.C:
		}
	}
}
