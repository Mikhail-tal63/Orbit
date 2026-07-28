package httperror

import (
	"errors"
	"net/http"

	"github.com/Mikhail-Tal63/Orbit/internal/auth"
	"github.com/Mikhail-Tal63/Orbit/utils/jsonR"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Handle(w http.ResponseWriter, err error) {

	switch {

	case errors.Is(err, auth.ErrInvalidCredentials):
		write(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())

	case errors.Is(err, auth.ErrWrongPassword):
		write(w, http.StatusUnauthorized, "WRONG_PASSWORD", err.Error())

	case errors.Is(err, auth.ErrUnauthorized):
		write(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())

	case errors.Is(err, auth.ErrInvalidToken):
		write(w, http.StatusUnauthorized, "INVALID_TOKEN", err.Error())

	case errors.Is(err, auth.ErrTokenExpired):
		write(w, http.StatusUnauthorized, "TOKEN_EXPIRED", err.Error())

	case errors.Is(err, auth.ErrUserNotFound):
		write(w, http.StatusNotFound, "USER_NOT_FOUND", err.Error())

	case errors.Is(err, auth.ErrEmailAlreadyExists):
		write(w, http.StatusConflict, "EMAIL_ALREADY_EXISTS", err.Error())

	case errors.Is(err, auth.ErrUsernameTaken):
		write(w, http.StatusConflict, "USERNAME_TAKEN", err.Error())

	default:
		write(
			w,
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
		)
	}
}

func write(w http.ResponseWriter, status int, code, message string) {
	_ = jsonR.WriteJSON(w, status, ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}