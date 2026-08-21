package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type IAuthService interface {
	Authenticate(username, password string) (string, error)
}

type ApiHandler struct {
	auth IAuthService
}

func NewApiHandler(auth IAuthService) *ApiHandler {
	return &ApiHandler{
		auth: auth,
	}
}

func respondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJson(w, status, ErrorResponse{
		Errors: message,
	})
}
