package purchase

import (
	"context"
	"errors"

	"github.com/tynou/avito-assignment/internal/db"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type IMerchRepo interface {
	GetMerchByName(ctx context.Context, name string) (db.Merch, error)
}

type IPurchaseRepo interface {
	CreatePurchase(ctx context.Context, userId int32, merch db.Merch) error
}

type PurchaseService struct {
	merchRepo    IMerchRepo
	purchaseRepo IPurchaseRepo
}

func NewPurchaseService(merchRepo IMerchRepo, purchaseRepo IPurchaseRepo) *PurchaseService {
	return &PurchaseService{
		merchRepo:    merchRepo,
		purchaseRepo: purchaseRepo,
	}
}

func (s *PurchaseService) Buy(ctx context.Context, userId int32, item string) error {
	merch, err := s.merchRepo.GetMerchByName(ctx, item)
	if err != nil {
		return err
	}

	err = s.purchaseRepo.CreatePurchase(ctx, userId, merch)
	if err != nil {
		return err
	}

	return nil
}
