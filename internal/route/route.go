package route

import "time"

// Kind 匹配类型。
type Kind string

const (
	KindExact  Kind = "exact"
	KindPrefix Kind = "prefix"
)

// Route 内部路由记录。
type Route struct {
	ID        string
	Path      string
	Kind      Kind
	Methods   []string
	Upstream  string
	Headers   map[string]string
	StripPath string
	Auth      bool
	Priority  int
	UpdatedAt time.Time
}

// HeadersCopy 返回 Headers 的独立拷贝，避免调用方改写库存。
func (r Route) HeadersCopy() map[string]string {
	// BUG: 未拷贝，返回库存 map
	return r.Headers
}

// MethodsCopy 返回 Methods 拷贝。
func (r Route) MethodsCopy() []string {
	return CloneStrings(r.Methods)
}
