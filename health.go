package apigatex

import "time"

// Health 健康检查视图。
type Health struct {
	OK        bool      `json:"ok"`
	Closed    bool      `json:"closed"`
	Routes    int       `json:"routes"`
	CheckedAt time.Time `json:"checked_at"`
}

// HealthCheck 返回当前健康快照。
func (g *Gateway) HealthCheck() Health {
	st := g.Stats()
	return Health{
		OK:        !st.Closed,
		Closed:    st.Closed,
		Routes:    st.Routes,
		CheckedAt: g.now(),
	}
}

func (g *Gateway) now() time.Time {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.clk == nil {
		return time.Now().UTC()
	}
	return g.clk.Now().UTC()
}
