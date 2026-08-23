package info

import (
	"context"

	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/http/handlers"
)

type IUserRepo interface {
	GetUserBalance(ctx context.Context, userId int32) (int32, error)
}

type IPurchaseRepo interface {
	GetUserInventory(ctx context.Context, userId int32) ([]db.GetUserInventoryRow, error)
}

type ITransferRepo interface {
	GetReceivedTransfers(ctx context.Context, userId int32) ([]db.GetReceivedTransfersRow, error)
	GetSentTransfers(ctx context.Context, userId int32) ([]db.GetSentTransfersRow, error)
}

type InfoService struct {
	userRepo     IUserRepo
	purchaseRepo IPurchaseRepo
	transferRepo ITransferRepo
}

func NewInfoService(userRepo IUserRepo, purchaseRepo IPurchaseRepo, transferRepo ITransferRepo) *InfoService {
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

	dbInventory, err := s.purchaseRepo.GetUserInventory(ctx, userId)
	if err != nil {
		return nil, err
	}

	inventory := make([]handlers.InventoryItem, len(dbInventory))
	for i, item := range dbInventory {
		inventory[i] = handlers.InventoryItem{
			Type:     item.Type,
			Quantity: item.Quantity,
		}
	}

	dbReceived, err := s.transferRepo.GetReceivedTransfers(ctx, userId)
	if err != nil {
		return nil, err
	}

	received := make([]handlers.ReceivedTransfer, len(dbReceived))
	for i, item := range dbReceived {
		received[i] = handlers.ReceivedTransfer{
			FromUser: item.FromUser,
			Amount:   item.Amount,
		}
	}

	dbSent, err := s.transferRepo.GetSentTransfers(ctx, userId)
	if err != nil {
		return nil, err
	}

	sent := make([]handlers.SentTransfer, len(dbSent))
	for i, item := range dbSent {
		sent[i] = handlers.SentTransfer{
			ToUser: item.ToUser,
			Amount: item.Amount,
		}
	}

	response := &handlers.InfoResponse{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: handlers.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}

	return response, nil
}
