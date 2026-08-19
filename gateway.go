package apigatex

import (
	"net/http"
	"sync"
	"time"

	"github.com/LYH2263/go-apigatex/internal/auth"
	"github.com/LYH2263/go-apigatex/internal/buffer"
	"github.com/LYH2263/go-apigatex/internal/clock"
	"github.com/LYH2263/go-apigatex/internal/proxy"
	"github.com/LYH2263/go-apigatex/internal/route"
	"github.com/LYH2263/go-apigatex/internal/upstream"
)

const (
	defaultMaxBody      = 1 << 20
	defaultProxyTimeout = 15 * time.Second
)

// Gateway API 网关门面。零值不可用，须经 New 构造。
type Gateway struct {
	mu sync.Mutex

	closed bool
	clk    clock.Clock
	table  *route.Table
	authn  auth.Authenticator
	client *http.Client
	pool   *upstream.Pool
	cache  *buffer.Cache

	persistPath  string
	proxyTimeout time.Duration
	maxBody      int
	defaultToken string
	allowHTTP    bool
	dirty        bool

	proxied     uint64
	authOK      uint64
	authFail    uint64
	upstreamOK  uint64
	upstreamErr uint64
}

// New 构造 Gateway。鉴权器为空时安装默认 Token 鉴权。
func New(opts ...Option) *Gateway {
	g := &Gateway{
		clk:          clock.Real{},
		table:        route.NewTable(),
		pool:         upstream.NewPool(),
		cache:        buffer.NewCache(defaultMaxBody),
		proxyTimeout: defaultProxyTimeout,
		maxBody:      defaultMaxBody,
		client:       &http.Client{Timeout: defaultProxyTimeout},
		allowHTTP:    true,
	}
	for _, o := range opts {
		if o != nil {
			o(g)
		}
	}
	if g.clk == nil {
		g.clk = clock.Real{}
	}
	if g.authn == nil {
		tok := g.defaultToken
		if tok == "" {
			tok = "dev-token"
		}
		g.authn = auth.NewTokenAuth(tok)
	}
	if g.client == nil {
		g.client = &http.Client{Timeout: g.proxyTimeout}
	}
	if g.table == nil {
		g.table = route.NewTable()
	}
	if g.pool == nil {
		g.pool = upstream.NewPool()
	}
	if g.cache == nil {
		g.cache = buffer.NewCache(g.maxBody)
	}
	if g.maxBody < 1 {
		g.maxBody = defaultMaxBody
	}
	g.cache.SetLimit(g.maxBody)
	return g
}

// Authenticator 返回当前鉴权器（只读引用）。
func (g *Gateway) Authenticator() auth.Authenticator {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.authn
}

// Transport 返回上游客户端（Close 后为 nil）。
func (g *Gateway) Transport() *http.Client {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.client
}

// ProxyEngine 构造一次代理引擎快照。
func (g *Gateway) ProxyEngine() *proxy.Engine {
	g.mu.Lock()
	defer g.mu.Unlock()
	return proxy.NewEngine(g.client, g.pool, g.proxyTimeout, g.allowHTTP)
}
