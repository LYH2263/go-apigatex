package apigatex

import (
	"net/http"
	"time"

	"github.com/LYH2263/go-apigatex/internal/auth"
	"github.com/LYH2263/go-apigatex/internal/clock"
)

// Option 配置 Gateway。
type Option func(*Gateway)

// WithClock 注入时钟（测试用 Fake）。
func WithClock(c clock.Clock) Option {
	return func(g *Gateway) {
		if c != nil {
			g.clk = c
		}
	}
}

// WithPersistPath 启用路由表落盘。
func WithPersistPath(path string) Option {
	return func(g *Gateway) { g.persistPath = path }
}

// WithAuthenticator 自定义鉴权；传 nil 时 New 仍会装默认 Token 实现。
func WithAuthenticator(a auth.Authenticator) Option {
	return func(g *Gateway) { g.authn = a }
}

// WithHTTPClient 自定义上游 HTTP 客户端。
func WithHTTPClient(c *http.Client) Option {
	return func(g *Gateway) {
		if c != nil {
			g.client = c
		}
	}
}

// WithProxyTimeout 上游请求超时。
func WithProxyTimeout(d time.Duration) Option {
	return func(g *Gateway) {
		if d > 0 {
			g.proxyTimeout = d
		}
	}
}

// WithMaxBody 缓冲请求体上限。
func WithMaxBody(n int) Option {
	return func(g *Gateway) {
		if n > 0 {
			g.maxBody = n
		}
	}
}

// WithDefaultToken 设置默认 Bearer token（仅默认鉴权器）。
func WithDefaultToken(token string) Option {
	return func(g *Gateway) { g.defaultToken = token }
}

// WithAllowInsecureUpstream 允许 http:// 上游（本地试请求）。
func WithAllowInsecureUpstream(v bool) Option {
	return func(g *Gateway) { g.allowHTTP = v }
}
