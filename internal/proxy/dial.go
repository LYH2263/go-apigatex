package proxy

import (
	"context"
	"net"
	"time"

	"github.com/LYH2263/go-apigatex/internal/errors"
)

// DialContext 带超时与取消的拨号辅助。
func DialContext(ctx context.Context, network, address string, timeout time.Duration) (net.Conn, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	d := net.Dialer{Timeout: timeout}
	conn, err := d.DialContext(ctx, network, address)
	if err != nil {
		if ctx.Err() != nil {
			return nil, errors.WrapErr(errors.ErrCanceled, ctx.Err())
		}
		return nil, errors.WrapErr(errors.ErrUpstream, err)
	}
	return conn, nil
}
