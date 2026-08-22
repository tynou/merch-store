package merch

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/repo"
)

type MerchRepo struct {
	queries *db.Queries
}

func NewMerchRepo(pool *pgxpool.Pool) *MerchRepo {
	return &MerchRepo{
		queries: db.New(pool),
	}
}

func (r *MerchRepo) GetMerchByName(ctx context.Context, name string) (db.Merch, error) {
	merch, err := r.queries.GetMerchByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Merch{}, repo.ErrNotFound
		}

		return db.Merch{}, err
	}
	return merch, nil
}
