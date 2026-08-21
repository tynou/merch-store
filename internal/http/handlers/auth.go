package handlers

import "net/http"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *ApiHandler) Auth(w http.ResponseWriter, r *http.Request) {

}
