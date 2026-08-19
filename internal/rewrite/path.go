package rewrite

import "strings"

// ApplyPath 应用 strip 前缀改写。
func ApplyPath(in, strip, routePath string) string {
	out := in
	if strip != "" && strings.HasPrefix(out, strip) {
		out = strings.TrimPrefix(out, strip)
		if out == "" {
			out = "/"
		}
	}
	if !strings.HasPrefix(out, "/") {
		out = "/" + out
	}
	_ = routePath
	return out
}

// JoinUpstream 拼接上游 base 与 path。
func JoinUpstream(base, path string) string {
	base = strings.TrimRight(base, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}
