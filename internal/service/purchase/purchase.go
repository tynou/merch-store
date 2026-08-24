package purchase

import (
	"context"

	"github.com/tynou/avito-assignment/internal/db"
)

type MerchRepo interface {
	GetMerchByName(ctx context.Context, name string) (db.Merch, error)
}

type PurchaseRepo interface {
	CreatePurchase(ctx context.Context, userId int32, merch db.Merch) error
}

type PurchaseService struct {
	merchRepo    MerchRepo
	purchaseRepo PurchaseRepo
}

func NewPurchaseService(merchRepo MerchRepo, purchaseRepo PurchaseRepo) *PurchaseService {
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
