package info

import (
	"context"

	"github.com/tynou/avito-assignment/internal/domain"
	"github.com/tynou/avito-assignment/internal/http/handlers"
)

type UserRepo interface {
	GetUserBalance(ctx context.Context, userId int32) (int32, error)
}

type PurchaseRepo interface {
	GetUserInventory(ctx context.Context, userId int32) ([]domain.InventoryItem, error)
}

type TransferRepo interface {
	GetReceivedTransfers(ctx context.Context, userId int32) ([]domain.ReceivedTransfer, error)
	GetSentTransfers(ctx context.Context, userId int32) ([]domain.SentTransfer, error)
}

type InfoService struct {
	userRepo     UserRepo
	purchaseRepo PurchaseRepo
	transferRepo TransferRepo
}

func NewInfoService(userRepo UserRepo, purchaseRepo PurchaseRepo, transferRepo TransferRepo) *InfoService {
	return &InfoService{
		userRepo:     userRepo,
		purchaseRepo: purchaseRepo,
		transferRepo: transferRepo,
	}
}

func (s *InfoService) GetInfo(ctx context.Context, userId int32) (*handlers.InfoResponse, error) {
	balance, err := s.userRepo.GetUserBalance(ctx, userId)
	if err != nil {
		return nil, err
	}

	inventory, err := s.purchaseRepo.GetUserInventory(ctx, userId)
	if err != nil {
		return nil, err
	}

	received, err := s.transferRepo.GetReceivedTransfers(ctx, userId)
	if err != nil {
		return nil, err
	}

	sent, err := s.transferRepo.GetSentTransfers(ctx, userId)
	if err != nil {
		return nil, err
	}

	response := &handlers.InfoResponse{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: domain.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}

	return response, nil
}
