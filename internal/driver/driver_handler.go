package driver

import (
	"net/http"

	"github.com/Mikhail-Tal63/Orbit/middleware"
	"github.com/Mikhail-Tal63/Orbit/utils/httperror"
	"github.com/Mikhail-Tal63/Orbit/utils/jsonR"
	"github.com/gorilla/mux"
)

type DriverHandler struct {
	service *DriverSevrice
}

func NewDriverHandler(service *DriverSevrice) *DriverHandler {
	return &DriverHandler{
		service: service,
	}
}

func DriverRouter(mux *mux.Router) {

}

func (h *DriverHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	var payload *CreateDriverRequest

	if err := jsonR.ParseJSON(r, &payload); err != nil {
		httperror.Handle(w, err)
		return
	}

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		httperror.Handle(w, err)
		return
	}

	driver, err := h.service.CreateDriver(r.Context(), userID, *payload)
	if err != nil {
		httperror.Handle(w, err)
		return
	}

	jsonR.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "Driver created",
		"driver":  driver,
	})
}
