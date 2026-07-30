package httperror

import (
	"errors"
	"net/http"

	"github.com/Mikhail-Tal63/Orbit/utils/jsonR"
)

type AppError interface {
	error
	HTTPStatus() int
	Code() string
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Handle(w http.ResponseWriter, err error) {
	var appErr AppError
	if errors.As(err, &appErr) {
		write(w, appErr.HTTPStatus(), appErr.Code(), appErr.Error())
		return
	}

	write(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
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
