package clock

import (
	"context"
	"time"

	"github.com/LYH2263/go-apigatex/internal/errors"
)

// Sleep 可取消睡眠。
func Sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return errors.WrapErr(errors.ErrCanceled, ctx.Err())
	case <-t.C:
		return nil
	}
}
