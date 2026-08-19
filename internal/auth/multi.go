package auth

import "net/http"

// Multi 按序尝试多个鉴权器。
type Multi []Authenticator

func (m Multi) Authenticate(r *http.Request) error {
	if len(m) == 0 {
		return ErrUnauthorized
	}
	var last error
	for _, a := range m {
		if a == nil {
			continue
		}
		if err := a.Authenticate(r); err == nil {
			return nil
		} else {
			last = err
		}
	}
	if last == nil {
		return ErrUnauthorized
	}
	return last
}
