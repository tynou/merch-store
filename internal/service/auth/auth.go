package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tynou/avito-assignment/internal/apperr"
	"github.com/tynou/avito-assignment/internal/config"
	"github.com/tynou/avito-assignment/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	CreateUser(ctx context.Context, username, passwordHash string) (int32, error)
}

type AuthService struct {
	cfg      *config.Config
	userRepo UserRepo
}

func NewAuthService(cfg *config.Config, userRepo UserRepo) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	var userId int32

	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if !errors.Is(err, apperr.ErrNotFound) {
			return "", err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		userId, err = s.userRepo.CreateUser(ctx, username, string(hash))
		if err != nil {
			return "", err
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			return "", apperr.ErrUnauthorized
		}

		userId = user.ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	signedToken, err := token.SignedString(s.cfg.JWTKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
