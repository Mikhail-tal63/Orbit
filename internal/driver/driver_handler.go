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

func (h *DriverHandler) DriverRouter(mux *mux.Router) {
	mux.HandleFunc("driver/online/{id}", h.DriverOnline).Methods("POST")
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

func (h *DriverHandler) GetDriverByUserId(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		httperror.Handle(w, err)
		return
	}
	driver, err := h.service.GetDriverByUserId(r.Context(), userID)
	if err != nil {
		httperror.Handle(w, err)
		return
	}

	if err := jsonR.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "here you go",
		"dirver":  driver,
	}); err != nil {
		httperror.Handle(w, err)
		return
	}
}
func (h *DriverHandler) DriverOnline(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		jsonR.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.service.DriverOnline(r.Context(), userID); err != nil {
		jsonR.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := jsonR.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "you are online",
	}); err != nil {
		jsonR.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
