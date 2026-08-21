package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tynou/avito-assignment/internal/db"
)

type UserRepo struct {
	queries *db.Queries
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		queries: db.New(pool),
	}
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	return r.queries.GetUserByUsername(ctx, username)
}

func (r *UserRepo) CreateUser(ctx context.Context, username, passwordHash string) (int32, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
	})
}
