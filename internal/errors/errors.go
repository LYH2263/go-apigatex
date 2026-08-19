package errors

import "errors"

var (
	ErrClosed       = errors.New("apigatex: closed")
	ErrNotFound     = errors.New("apigatex: route not found")
	ErrUnauthorized = errors.New("apigatex: unauthorized")
	ErrForbidden    = errors.New("apigatex: forbidden")
	ErrBadRoute     = errors.New("apigatex: bad route")
	ErrUpstream     = errors.New("apigatex: upstream error")
	ErrPersist      = errors.New("apigatex: persist failed")
	ErrCanceled     = errors.New("apigatex: canceled")
	ErrNilAuth      = errors.New("apigatex: nil authenticator")
	ErrNilTransport = errors.New("apigatex: nil transport")
	ErrEmptyBody    = errors.New("apigatex: empty body")
	ErrTimeout      = errors.New("apigatex: timeout")
	ErrReload       = errors.New("apigatex: reload failed")
	ErrSync         = errors.New("apigatex: sync failed")
	ErrTooLarge     = errors.New("apigatex: body too large")
)

// Is 转发标准库。
func Is(err, target error) bool { return errors.Is(err, target) }

// As 转发标准库。
func As(err error, target any) bool { return errors.As(err, target) }
