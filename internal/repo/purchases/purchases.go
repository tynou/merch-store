package purchases

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/service/purchase"
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
		return purchase.ErrInsufficientFunds
	}

	newBalance := user.Balance - merch.Price
	err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
		Balance: newBalance,
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
