package auth

import (
	"net/http"

	"github.com/Mikhail-Tal63/Orbit/middleware"
	"github.com/Mikhail-Tal63/Orbit/utils/httperror"
	"github.com/Mikhail-Tal63/Orbit/utils/jsonR"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h AuthHandler) AuthRouter(mux *mux.Router) {
	mux.HandleFunc("/regester", h.CreateUser).Methods("POST")
	mux.HandleFunc("/profile/{id}",h.GetUserByID).Methods("GET")
}
func (h *AuthHandler) ProtectedRouter(router *mux.Router) {
	//router.HandleFunc("/users/{username}", h.GetUserByUsername).Methods("GET")
}
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload RegisterRequest

	if err := jsonR.ParseJSON(r, &payload); err != nil {
		httperror.Handle(w, err)
		return
	}

	createUser, err := h.service.CreateUser(r.Context(), &payload)
	if err != nil {
		httperror.Handle(w, err)
		return
	}

	if err := jsonR.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "user created successfully",
		"user":    createUser,
	}); err != nil {
		httperror.Handle(w, err)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest

	if err := jsonR.ParseJSON(r, &payload); err != nil {
		httperror.Handle(w, err)
		return
	}

	authResponse, err := h.service.Login(
		r.Context(),
		payload.Email,
		payload.Password,
	)
	if err != nil {
		httperror.Handle(w, err)
		return
	}

	if err := jsonR.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "login succeeded",
		"data":    authResponse,
	}); err != nil {
		httperror.Handle(w, err)
		return
	}
}

func (h *AuthHandler) GetUserByID(w http.ResponseWriter,r *http.Request) {
userID,err:= middleware.GetUserID(r.Context())
if err != nil {
	httperror.Handle(w,err)
	return
}
user,err:= h.service.GetUserByID(r.Context(),userID)
if err != nil {
	httperror.Handle(w,err)
	return
}
if err:= jsonR.WriteJSON(w,http.StatusOK,map[string]any{
	"user":user,
});err!= nil {
	httperror.Handle(w,err)
	return
}
}
