package ride

import (
	"net/http"

	"github.com/Mikhail-Tal63/Orbit/middleware"
	"github.com/Mikhail-Tal63/Orbit/utils/httperror"
	"github.com/Mikhail-Tal63/Orbit/utils/jsonR"
	"github.com/gorilla/mux"
)

type RideHandler struct{
	service *RideService
}

func NewRideHandler (service *RideService) *RideHandler{
	return &RideHandler{
		service: service,
	}
}

func(h *RideHandler) RideRouter(mux *mux.Router){
mux.HandleFunc("/ride_request",h.CreateRideReq).Methods("POST")
}

func(h RideHandler) CreateRideReq(w http.ResponseWriter,req *http.Request){
	var payload RideRequestDto

    if err := jsonR.ParseJSON(req,&payload);err!= nil {
		httperror.Handle(w,err)
		return
	}


	userid,err := middleware.GetUserID(req.Context())
	if err != nil {
		httperror.Handle(w,err)
		return
	}
	request,err := h.service.CreateRideReq(req.Context(),payload,userid)
    if err != nil {
		httperror.Handle(w,err)
		return
	}

	if err := jsonR.WriteJSON(w,http.StatusCreated,map[string]any{
		"message":"ride created",
		"ride_request":request,
	});err != nil{
	 httperror.Handle(w,err)
	 return
	}
}
