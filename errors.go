package apigatex

import ierr "github.com/LYH2263/go-apigatex/internal/errors"

var (
	ErrClosed       = ierr.ErrClosed
	ErrNotFound     = ierr.ErrNotFound
	ErrUnauthorized = ierr.ErrUnauthorized
	ErrForbidden    = ierr.ErrForbidden
	ErrBadRoute     = ierr.ErrBadRoute
	ErrUpstream     = ierr.ErrUpstream
	ErrPersist      = ierr.ErrPersist
	ErrCanceled     = ierr.ErrCanceled
	ErrNilAuth      = ierr.ErrNilAuth
	ErrNilTransport = ierr.ErrNilTransport
	ErrEmptyBody    = ierr.ErrEmptyBody
	ErrTimeout      = ierr.ErrTimeout
	ErrReload       = ierr.ErrReload
	ErrSync         = ierr.ErrSync
)
