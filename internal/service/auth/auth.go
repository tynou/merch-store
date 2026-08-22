package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/repo/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized = errors.New("user is not authorized")
)

// TODO (TASK-001): Заменить хардкод на взятие ключа из конфига
var JWTSecret = []byte("my_super_secret_key")

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

func (s *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	var userId int32

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if !errors.Is(err, users.ErrNotFound) {
			return "", err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}

		userId, err = s.repo.CreateUser(ctx, username, string(hash))
		if err != nil {
			return "", err
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			return "", ErrUnauthorized
		}

		userId = user.ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	signedToken, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
