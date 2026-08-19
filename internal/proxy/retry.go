package proxy

import (
	"context"
	"time"

	"github.com/LYH2263/go-apigatex/internal/errors"
)

// RetryPolicy 简单重试策略。
type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

// DefaultRetry 默认策略。
func DefaultRetry() RetryPolicy {
	return RetryPolicy{MaxAttempts: 2, Backoff: 20 * time.Millisecond}
}

// ForwardWithRetry 在可重试状态码或瞬时错误时重试；仍尊重 ctx。
func (e *Engine) ForwardWithRetry(ctx context.Context, in Request, pol RetryPolicy) (*Response, error) {
	if pol.MaxAttempts < 1 {
		pol.MaxAttempts = 1
	}
	var last error
	for i := 0; i < pol.MaxAttempts; i++ {
		if err := ctx.Err(); err != nil {
			return nil, errors.WrapErr(errors.ErrCanceled, err)
		}
		res, err := e.Forward(ctx, in)
		if err == nil {
			if IsRetryable(res.StatusCode) && i+1 < pol.MaxAttempts {
				select {
				case <-ctx.Done():
					return nil, errors.WrapErr(errors.ErrCanceled, ctx.Err())
				case <-time.After(pol.Backoff):
				}
				continue
			}
			return res, nil
		}
		last = err
		if errors.Is(err, errors.ErrCanceled) {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, errors.WrapErr(errors.ErrCanceled, ctx.Err())
		case <-time.After(pol.Backoff):
		}
	}
	return nil, last
}
