package constants

import "errors"

var (
	ErrorUnauthorized        = errors.New("unauthorized")
	ErrorPermissionDenied    = errors.New("permission denied")
	ErrorBadRequest          = errors.New("bad request")
	ErrorNotFound            = errors.New("not found")
	ErrorInternalServerError = errors.New("internal server error")
	ErrorConflict            = errors.New("conflict")
	ErrorOAuthWithPassWord   = errors.New("cannot use password with an account that has OAuth")
)
