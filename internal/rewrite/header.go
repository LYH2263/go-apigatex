package rewrite

import (
	"net/http"

	"github.com/LYH2263/go-apigatex/internal/route"
)

// MergeHeaders 克隆入站 header 并覆盖路由配置。
func MergeHeaders(in http.Header, extra map[string]string) http.Header {
	out := in.Clone()
	if out == nil {
		out = make(http.Header)
	}
	for k, v := range route.CloneHeaders(extra) {
		out.Set(k, v)
	}
	return out
}

// DropHopByHop 删除逐跳头。
func DropHopByHop(h http.Header) {
	for _, k := range []string{
		"Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization",
		"Te", "Trailers", "Transfer-Encoding", "Upgrade",
	} {
		h.Del(k)
	}
}
