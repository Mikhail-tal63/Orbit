package auth

import (
	"net/http"
)


type authError struct {
	msg    string
	status int
	code   string
}

func (e *authError) Error() string      { return e.msg }
func (e *authError) HTTPStatus() int    { return e.status }
func (e *authError) Code() string       { return e.code }

var (
	ErrUserNotFound       error = &authError{"user not found", http.StatusNotFound, "USER_NOT_FOUND"}
	ErrEmailAlreadyExists error = &authError{"email already exists", http.StatusConflict, "EMAIL_ALREADY_EXISTS"}
	ErrUsernameTaken      error = &authError{"username already taken", http.StatusConflict, "USERNAME_TAKEN"}
	ErrInvalidCredentials error = &authError{"invalid email or password", http.StatusUnauthorized, "INVALID_CREDENTIALS"}
	ErrUnauthorized       error = &authError{"unauthorized", http.StatusUnauthorized, "UNAUTHORIZED"}
	ErrInvalidToken       error = &authError{"invalid token", http.StatusUnauthorized, "INVALID_TOKEN"}
	ErrTokenExpired       error = &authError{"token expired", http.StatusUnauthorized, "TOKEN_EXPIRED"}
	ErrWrongPassword      error = &authError{"wrong email or password", http.StatusUnauthorized, "WRONG_PASSWORD"}
)
