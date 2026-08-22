package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
)

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
