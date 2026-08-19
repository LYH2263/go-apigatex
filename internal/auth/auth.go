package auth

import "net/http"

// Authenticator 鉴权接口。
type Authenticator interface {
	Authenticate(r *http.Request) error
}
