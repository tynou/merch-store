package purchases

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tynou/avito-assignment/internal/apperr"
	"github.com/tynou/avito-assignment/internal/db"
)

type PurchaseRepo struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewPurchaseRepo(pool *pgxpool.Pool) *PurchaseRepo {
	return &PurchaseRepo{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *PurchaseRepo) CreatePurchase(ctx context.Context, userId int32, merch db.Merch) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	user, err := qtx.GetUserForUpdate(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to lock user: %w", err)
	}

	if user.Balance < merch.Price {
		return apperr.ErrInsufficientFunds
	}

	err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
		Balance: user.Balance - merch.Price,
		ID:      userId,
	})
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	err = qtx.CreatePurchase(ctx, db.CreatePurchaseParams{
		UserID:  userId,
		MerchID: merch.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create purchase: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *PurchaseRepo) GetUserInventory(ctx context.Context, userId int32) ([]db.GetUserInventoryRow, error) {
	inventory, err := r.queries.GetUserInventory(ctx, userId)
	if err != nil {
		return nil, err
	}
	return inventory, nil
}
