package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type IAuthService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type IPurchaseService interface {
	Buy(ctx context.Context, userId int32, item string) error
}

type ApiHandler struct {
	validate *validator.Validate
	auth     IAuthService
	purchase IPurchaseService
}

func NewApiHandler(validate *validator.Validate, auth IAuthService, purchase IPurchaseService) *ApiHandler {
	return &ApiHandler{
		validate: validate,
		auth:     auth,
		purchase: purchase,
	}
}
