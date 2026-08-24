package handlers

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
}

type PurchaseService interface {
	Buy(ctx context.Context, userId int32, item string) error
}

type InfoService interface {
	GetInfo(ctx context.Context, userId int32) (*InfoResponse, error)
}

type TransferService interface {
	SendCoin(ctx context.Context, fromUserId int32, toUserName string, amount int32) error
}

type ApiHandler struct {
	validate *validator.Validate
	auth     AuthService
	purchase PurchaseService
	info     InfoService
	transfer TransferService
}

func NewApiHandler(
	validate *validator.Validate,
	auth AuthService,
	purchase PurchaseService,
	info InfoService,
	transfer TransferService,
) *ApiHandler {
	return &ApiHandler{
		validate: validate,
		auth:     auth,
		purchase: purchase,
		info:     info,
		transfer: transfer,
	}
}
