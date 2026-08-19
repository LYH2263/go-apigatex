package apigatex

// Metrics 导出计数器（只读拷贝）。
type Metrics struct {
	Proxied     uint64
	AuthOK      uint64
	AuthFail    uint64
	UpstreamOK  uint64
	UpstreamErr uint64
}

// Metrics 返回计数快照。
func (g *Gateway) Metrics() Metrics {
	g.mu.Lock()
	defer g.mu.Unlock()
	return Metrics{
		Proxied:     g.proxied,
		AuthOK:      g.authOK,
		AuthFail:    g.authFail,
		UpstreamOK:  g.upstreamOK,
		UpstreamErr: g.upstreamErr,
	}
}

// ResetMetrics 清零计数（测试用）。
func (g *Gateway) ResetMetrics() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.proxied = 0
	g.authOK = 0
	g.authFail = 0
	g.upstreamOK = 0
	g.upstreamErr = 0
}
