package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type IAuthService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type ApiHandler struct {
	validate *validator.Validate
	auth     IAuthService
}

func NewApiHandler(validate *validator.Validate, auth IAuthService) *ApiHandler {
	return &ApiHandler{
		validate: validate,
		auth:     auth,
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
