package apigatex

import "time"

// MatchKind 路由匹配方式。
type MatchKind string

const (
	MatchExact  MatchKind = "exact"
	MatchPrefix MatchKind = "prefix"
)

// RouteSpec 对外路由配置视图。
type RouteSpec struct {
	ID        string
	Path      string
	Kind      MatchKind
	Methods   []string
	Upstream  string
	Headers   map[string]string
	StripPath string
	Auth      bool
	Priority  int
	UpdatedAt time.Time
}

// ProxyResult 试请求/代理结果摘要。
type ProxyResult struct {
	Status     int
	Upstream   string
	BytesOut   int
	DurationMs int64
	Headers    map[string]string
}

// Stats 网关运行统计。
type Stats struct {
	Routes      int
	Proxied     uint64
	AuthOK      uint64
	AuthFail    uint64
	UpstreamOK  uint64
	UpstreamErr uint64
	Closed      bool
}
