package api

import (
	"net/http"
	"time"
)

// WithTimeout 包装超时。
func WithTimeout(d time.Duration, next http.Handler) http.Handler {
	return http.TimeoutHandler(next, d, "timeout")
}

// WithRequestID 简单请求 ID 头。
func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Id") == "" {
			r.Header.Set("X-Request-Id", time.Now().UTC().Format("20060102T150405.000"))
		}
		next.ServeHTTP(w, r)
	})
}
