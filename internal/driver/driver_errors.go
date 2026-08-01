package driver

import "net/http"

type DriverErrors struct {
	msg    string
	status int
	code   string
}

func (e *DriverErrors) Error() string {
	return e.msg
}

func (e *DriverErrors) HTTPStatus() int {
	return e.status
}

func (e *DriverErrors) Code() string {
	return e.code
}

var (
	ErrDriverNotFound = &DriverErrors{
		msg:    "driver not found",
		status: http.StatusNotFound,
		code:   "DRIVER_NOT_FOUND",
	}

	ErrDriverAlreadyExists = &DriverErrors{
		msg:    "driver already exists",
		status: http.StatusConflict,
		code:   "DRIVER_ALREADY_EXISTS",
	}

	ErrDriverNotAvailable = &DriverErrors{
		msg:    "driver is not available",
		status: http.StatusBadRequest,
		code:   "DRIVER_NOT_AVAILABLE",
	}

	ErrDriverOffline = &DriverErrors{
		msg:    "driver is offline",
		status: http.StatusBadRequest,
		code:   "DRIVER_OFFLINE",
	}

	ErrDriverNotFoundForUser = &DriverErrors{
		msg:    "user is not a driver",
		status: http.StatusForbidden,
		code:   "USER_NOT_DRIVER",
	}

	ErrDriverAlreadyOnline = &DriverErrors{
		msg:    "driver already online",
		status: http.StatusConflict,
		code:   "DRIVER_ALREADY_ONLINE",
	}

	ErrDriverAlreadyOffline = &DriverErrors{
		msg:    "driver already offline",
		status: http.StatusConflict,
		code:   "DRIVER_ALREADY_OFFLINE",
	}

	ErrInvalidDriverLocation = &DriverErrors{
		msg:    "invalid driver location",
		status: http.StatusBadRequest,
		code:   "INVALID_DRIVER_LOCATION",
	}
)