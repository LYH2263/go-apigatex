package apigatex

import (
	"net/url"
	"strings"
)

// ValidateUpstreamURL 校验上游 URL。
func ValidateUpstreamURL(raw string, allowHTTP bool) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ErrBadRoute
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ErrBadRoute
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		return nil
	case "http":
		if allowHTTP {
			return nil
		}
		return ErrBadRoute
	default:
		return ErrBadRoute
	}
}

// ValidatePathPattern 路径须以 / 开头。
func ValidatePathPattern(p string) error {
	if !strings.HasPrefix(p, "/") {
		return ErrBadRoute
	}
	if strings.Contains(p, " ") {
		return ErrBadRoute
	}
	return nil
}
