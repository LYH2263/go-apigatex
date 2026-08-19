package auth

import "net/http"

// AllowAll 测试用放行。
type AllowAll struct{}

func (AllowAll) Authenticate(r *http.Request) error { return nil }

// DenyAll 测试用拒绝。
type DenyAll struct{}

func (DenyAll) Authenticate(r *http.Request) error { return ErrUnauthorized }
