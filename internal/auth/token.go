package auth

import (
	"net/http"
	"strings"

	"github.com/LYH2263/go-apigatex/internal/errors"
)

// ErrUnauthorized 鉴权失败哨兵。
var ErrUnauthorized = errors.ErrUnauthorized

// TokenAuth Bearer/header token 校验。
type TokenAuth struct {
	Token      string
	HeaderName string
}

// NewTokenAuth 构造。
func NewTokenAuth(token string) *TokenAuth {
	return &TokenAuth{Token: token, HeaderName: "Authorization"}
}

// Authenticate 校验 Authorization: Bearer <token> 或裸 token。
func (a *TokenAuth) Authenticate(r *http.Request) error {
	if a == nil {
		return errors.ErrNilAuth
	}
	if r == nil {
		return ErrUnauthorized
	}
	h := a.HeaderName
	if h == "" {
		h = "Authorization"
	}
	raw := strings.TrimSpace(r.Header.Get(h))
	if raw == "" {
		return ErrUnauthorized
	}
	const prefix = "Bearer "
	if strings.HasPrefix(raw, prefix) {
		raw = strings.TrimSpace(raw[len(prefix):])
	}
	if raw != a.Token {
		return ErrUnauthorized
	}
	return nil
}
