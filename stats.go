package apigatex

// Stats 返回运行统计快照。
func (g *Gateway) Stats() Stats {
	g.mu.Lock()
	defer g.mu.Unlock()
	n := 0
	if g.table != nil {
		n = g.table.Len()
	}
	return Stats{
		Routes:      n,
		Proxied:     g.proxied,
		AuthOK:      g.authOK,
		AuthFail:    g.authFail,
		UpstreamOK:  g.upstreamOK,
		UpstreamErr: g.upstreamErr,
		Closed:      g.closed,
	}
}
