package auth

import (
	"net/http"
	"strings"
)

// HeaderAuth 校验固定 header 值。
type HeaderAuth struct {
	Name  string
	Value string
}

func (h HeaderAuth) Authenticate(r *http.Request) error {
	if r == nil || h.Name == "" {
		return ErrUnauthorized
	}
	if strings.TrimSpace(r.Header.Get(h.Name)) != h.Value {
		return ErrUnauthorized
	}
	return nil
}
