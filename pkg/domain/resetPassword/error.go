package resetPassword

import "errors"

// ErrResetPasswordCodeExpired is the error returned when a reset password code has expired
var ErrResetPasswordCodeExpired = errors.New("reset password code has expired")

// ErrResetPasswordCodeInvalid is the error returned when a reset password code is invalid
var ErrResetPasswordCodeInvalid = errors.New("reset password code is invalid")

// ErrResetPasswordNotFound is the error returned when a reset password is not found
var ErrResetPasswordNotFound = errors.New("reset password not found")
