package transfers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tynou/avito-assignment/internal/apperr"
	"github.com/tynou/avito-assignment/internal/db"
)

type TransferRepo struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewTransferRepo(pool *pgxpool.Pool) *TransferRepo {
	return &TransferRepo{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *TransferRepo) CreateTransfer(ctx context.Context, fromUserId, toUserId int32, amount int32) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// блокируем пользователей в порядке возрастания Id, чтобы не было дедлоков
	var fromUser, toUser db.GetUserForUpdateRow
	if fromUserId < toUserId {
		fromUser, err = qtx.GetUserForUpdate(ctx, fromUserId)
		if err != nil {
			return err
		}
		toUser, err = qtx.GetUserForUpdate(ctx, toUserId)
		if err != nil {
			return err
		}
	} else {
		toUser, err = qtx.GetUserForUpdate(ctx, toUserId)
		if err != nil {
			return err
		}
		fromUser, err = qtx.GetUserForUpdate(ctx, fromUserId)
		if err != nil {
			return err
		}
	}

	if fromUser.Balance < amount {
		return apperr.ErrInsufficientFunds
	}

	err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
		Balance: fromUser.Balance - amount,
		ID:      fromUserId,
	})
	if err != nil {
		return err
	}

	err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
		Balance: toUser.Balance + amount,
		ID:      toUserId,
	})
	if err != nil {
		return err
	}

	err = qtx.CreateTransfer(ctx, db.CreateTransferParams{
		FromID: fromUserId,
		ToID:   toUserId,
		Amount: amount,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *TransferRepo) GetReceivedTransfers(ctx context.Context, userId int32) ([]db.GetReceivedTransfersRow, error) {
	transfers, err := r.queries.GetReceivedTransfers(ctx, userId)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (r *TransferRepo) GetSentTransfers(ctx context.Context, userId int32) ([]db.GetSentTransfersRow, error) {
	transfers, err := r.queries.GetSentTransfers(ctx, userId)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}
