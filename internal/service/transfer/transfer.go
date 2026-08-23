package transfer

import (
	"context"

	"github.com/tynou/avito-assignment/internal/apperr"
	"github.com/tynou/avito-assignment/internal/db"
)

type IUserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
}

type ITransferRepo interface {
	CreateTransfer(ctx context.Context, fromUserId, toUserId int32, amount int32) error
}

type TransferService struct {
	userRepo     IUserRepo
	transferRepo ITransferRepo
}

func NewTransferService(userRepo IUserRepo, transferRepo ITransferRepo) *TransferService {
	return &TransferService{
		userRepo:     userRepo,
		transferRepo: transferRepo,
	}
}

func (s *TransferService) SendCoin(ctx context.Context, fromUserId int32, toUserName string, amount int32) error {
	user, err := s.userRepo.GetUserByUsername(ctx, toUserName)
	if err != nil {
		return err
	}

	toUserId := user.ID

	if fromUserId == toUserId {
		return apperr.ErrSelfTransfer
	}

	err = s.transferRepo.CreateTransfer(ctx, fromUserId, toUserId, amount)
	if err != nil {
		return err
	}

	return nil
}
