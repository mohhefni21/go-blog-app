package errorpkg

import (
	"errors"
	"net/http"
)

// General error
var (
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbiddenAccess = errors.New("forbidden access")
)

var (
	// Auth
	ErrUsernameInvalid     = errors.New("username is invalid")
	ErrUsernameRequired    = errors.New("username is required")
	ErrEmailInvalid        = errors.New("email is invalid")
	ErrEmailRequired       = errors.New("email is required")
	ErrPasswordInvalid     = errors.New("password is invalid")
	ErrPasswordRequired    = errors.New("password is required")
	ErrFullnameRequired    = errors.New("fullname is required")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrPasswordNotMatch    = errors.New("password not match")
)

type Error struct {
	Message  string
	HttpCode int
}

func (e Error) Error() string {
	return e.Message
}

func NewError(message string, httpCode int) Error {
	return Error{
		Message:  message,
		HttpCode: httpCode,
	}
}

var (
	ErrorGeneral         = NewError("internal server error", http.StatusInternalServerError)
	ErrorBadRequest      = NewError("bad request", http.StatusBadRequest)
	ErrorNotFound        = NewError(ErrNotFound.Error(), http.StatusNotFound)
	ErrorUnauthorized    = NewError(ErrUnauthorized.Error(), http.StatusUnauthorized)
	ErrorForbiddenAccess = NewError(ErrForbiddenAccess.Error(), http.StatusForbidden)
)

var (
	ErrorUsernameInvalid     = NewError(ErrUsernameInvalid.Error(), http.StatusBadRequest)
	ErrorUsernameRequired    = NewError(ErrUsernameRequired.Error(), http.StatusBadRequest)
	ErrorEmailInvalid        = NewError(ErrEmailInvalid.Error(), http.StatusBadRequest)
	ErrorEmailRequired       = NewError(ErrEmailRequired.Error(), http.StatusBadRequest)
	ErrorPasswordInvalid     = NewError(ErrPasswordInvalid.Error(), http.StatusBadRequest)
	ErrorPasswordRequired    = NewError(ErrPasswordRequired.Error(), http.StatusBadRequest)
	ErrorFullnameRequired    = NewError(ErrFullnameRequired.Error(), http.StatusBadRequest)
	ErrorEmailAlreadyUsed    = NewError(ErrEmailAlreadyUsed.Error(), http.StatusConflict)
	ErrorUsernameAlreadyUsed = NewError(ErrUsernameAlreadyUsed.Error(), http.StatusConflict)
	ErrorPasswordNotMatch    = NewError(ErrPasswordNotMatch.Error(), http.StatusUnauthorized)
)

var ErrorMapping = map[string]Error{
	ErrUsernameInvalid.Error():     ErrorUsernameInvalid,
	ErrUsernameRequired.Error():    ErrorUsernameRequired,
	ErrEmailInvalid.Error():        ErrorEmailInvalid,
	ErrEmailRequired.Error():       ErrorEmailRequired,
	ErrPasswordInvalid.Error():     ErrorPasswordInvalid,
	ErrPasswordRequired.Error():    ErrorPasswordRequired,
	ErrFullnameRequired.Error():    ErrorFullnameRequired,
	ErrEmailAlreadyUsed.Error():    ErrorEmailAlreadyUsed,
	ErrUsernameAlreadyUsed.Error(): ErrorUsernameAlreadyUsed,
	ErrPasswordNotMatch.Error():    ErrorPasswordNotMatch,
}
