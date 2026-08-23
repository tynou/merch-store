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

type IInfoService interface {
	GetInfo(ctx context.Context, userId int32) (*InfoResponse, error)
}

type ITransferService interface {
	SendCoin(ctx context.Context, fromUserId int32, toUserName string, amount int32) error
}

type ApiHandler struct {
	validate *validator.Validate
	auth     IAuthService
	purchase IPurchaseService
	info     IInfoService
	transfer ITransferService
}

func NewApiHandler(
	validate *validator.Validate,
	auth IAuthService,
	purchase IPurchaseService,
	info IInfoService,
	transfer ITransferService,
) *ApiHandler {
	return &ApiHandler{
		validate: validate,
		auth:     auth,
		purchase: purchase,
		info:     info,
		transfer: transfer,
	}
}
