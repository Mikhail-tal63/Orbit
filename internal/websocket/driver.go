package websocket

import (
	"net/http"

	"github.com/Mikhail-Tal63/Orbit/internal/driver"

	"github.com/gorilla/mux"
)

type DriverHandler struct {
    service    *driver.DriverSevrice
    hub        *Hub
    jwtSecret  []byte
    driverRepo driver.DriverRepository
}

func NewDriverHandler(
    service *driver.DriverSevrice,
    hub *Hub,
    jwtSecret []byte,
    driverRepo driver.DriverRepository,
) *DriverHandler {
    return &DriverHandler{
        service:    service,
        hub:        hub,
        jwtSecret:  jwtSecret,
        driverRepo: driverRepo,
    }
}
func (h *DriverHandler) DriverRouter(r *mux.Router) {

    r.HandleFunc("/drivers/ws", func(w http.ResponseWriter, req *http.Request) {
       ServerWS(
            h.hub,
            h.jwtSecret,
            h.driverRepo,
            w,
            req,
        )
    })

}