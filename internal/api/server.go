package api

import (
	"net/http"

	"github.com/LYH2263/go-apigatex"
)

// Options 管理服务选项。
type Options struct {
	WebDir    string
	AllowCORS bool
}

// Server 管理 API + 静态页。
type Server struct {
	gw   *apigatex.Gateway
	opts Options
	mux  *http.ServeMux
}

// New 构造。
func New(gw *apigatex.Gateway, opts Options) *Server {
	s := &Server{gw: gw, opts: opts, mux: http.NewServeMux()}
	s.routes()
	return s
}

// ServeHTTP 入口。
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.opts.AllowCORS {
		applyCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.mux.ServeHTTP(w, r)
}
