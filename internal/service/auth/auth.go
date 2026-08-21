package auth

import (
	"context"
	"errors"

	"github.com/tynou/avito-assignment/internal/db"
)

var (
	ErrUnauthorized = errors.New("Вы не авторизованы.")
)

type IUserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	CreateUser(ctx context.Context, username, passwordHash string) (int32, error)
}

type AuthService struct {
	repo IUserRepo
}

func NewAuthService(repo IUserRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	return "jwt token here", nil
}
