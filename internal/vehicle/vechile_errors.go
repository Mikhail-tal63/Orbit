package vehicle

import "net/http"

type VehicleErrors struct {
	msg    string
	status int
	code   string
}

func (e *VehicleErrors) Error() string {
	return e.msg
}

func (e *VehicleErrors) HTTPStatus() int {
	return e.status
}

func (e *VehicleErrors) Code() string {
	return e.code
}

var (
	ErrVehicleNotFound = &VehicleErrors{
		msg:    "vehicle not found",
		status: http.StatusNotFound,
		code:   "VEHICLE_NOT_FOUND",
	}

	ErrVehicleAlreadyExists = &VehicleErrors{
		msg:    "vehicle already exists",
		status: http.StatusConflict,
		code:   "VEHICLE_ALREADY_EXISTS",
	}

	ErrInvalidVehicleData = &VehicleErrors{
		msg:    "invalid vehicle data",
		status: http.StatusBadRequest,
		code:   "INVALID_VEHICLE_DATA",
	}

	ErrVehicleNotOwned = &VehicleErrors{
		msg:    "vehicle does not belong to user",
		status: http.StatusForbidden,
		code:   "VEHICLE_NOT_OWNED",
	}

	ErrVehicleUnavailable = &VehicleErrors{
		msg:    "vehicle is not available",
		status: http.StatusBadRequest,
		code:   "VEHICLE_UNAVAILABLE",
	}

	ErrVehicleAlreadyAssigned = &VehicleErrors{
		msg:    "vehicle already assigned",
		status: http.StatusConflict,
		code:   "VEHICLE_ALREADY_ASSIGNED",
	}

	ErrInvalidVehicleType = &VehicleErrors{
		msg:    "invalid vehicle type",
		status: http.StatusBadRequest,
		code:   "INVALID_VEHICLE_TYPE",
	}
)
