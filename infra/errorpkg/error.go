package errorpkg

import "errors"

// General error

var (
	ErrUsernameInvalid  = errors.New("username is invalid")
	ErrUsernameRequired = errors.New("username is required")
	ErrEmailInvalid     = errors.New("email is invalid")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordInvalid  = errors.New("password is invalid")
	ErrPasswordRequired = errors.New("password is required")
	ErrFullnameRequired = errors.New("fullname is required")
)
