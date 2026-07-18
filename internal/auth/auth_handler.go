package auth

import (
	"net/http"

	
	"github.com/Mikhail-Tal63/Orbit/utils/jsonR"
	"github.com/gorilla/mux"
)

type AuthHandler struct{
	service AuthService
}

func NewAuthHandler(service AuthService)*AuthHandler{
	return &AuthHandler{
		service: service,
	}
}

func (h AuthHandler) AuthRouter(mux *mux.Router){
	mux.HandleFunc("/regester",h.CreateUser).Methods("POST")
}
func (h *AuthHandler) ProtectedRouter(router *mux.Router) {
	//router.HandleFunc("/users/{username}", h.GetUserByUsername).Methods("GET")
}
func (h *AuthHandler) CreateUser(w http.ResponseWriter,r *http.Request){
	var payload RegisterRequest
	if err := jsonR.ParseJSON(r,&payload); err != nil{
		jsonR.WriteError(w,http.StatusInternalServerError,err)
		return
	}

    createuser,err:= h.service.CreateUser(r.Context(),&payload)
	if err != nil {
		jsonR.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	if err := jsonR.WriteJSON(w,http.StatusOK,map[string]any{
		"message":"user created seccessfuly",
		"user": createuser,
	});err != nil {
			jsonR.WriteError(w,http.StatusInternalServerError,err)
		return
	}
	

}
func (h *AuthHandler) GetUserByUsername(){
	
}