package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tynou/avito-assignment/internal/apperr"
	"github.com/tynou/avito-assignment/internal/config"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/service/auth/mocks"
	"golang.org/x/crypto/bcrypt"
)

var cfg = &config.Config{
	JWTKey: []byte("my_super_secret_key"),
}

func Test_Authenticate_UserExists_Success(t *testing.T) {
	ctx := context.Background()
	username, password := "user", "pass"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUserRepo := mocks.NewMockUserRepo(t)

	mockUserRepo.On("GetUserByUsername", ctx, username).Return(db.User{
		ID:           1,
		Username:     username,
		PasswordHash: string(passwordHash),
		Balance:      1000,
	}, nil)

	service := NewAuthService(cfg, mockUserRepo)

	token, err := service.Authenticate(ctx, username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func Test_Authenticate_UserDoesNotExist_Success(t *testing.T) {
	ctx := context.Background()
	username, password := "user", "pass"

	mockUserRepo := mocks.NewMockUserRepo(t)

	mockUserRepo.On("GetUserByUsername", ctx, username).Return(db.User{}, apperr.ErrNotFound)
	mockUserRepo.On("CreateUser", ctx, username, mock.MatchedBy(func(passwordHash string) bool {
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		return err == nil
	})).Return(int32(1), nil)

	service := NewAuthService(cfg, mockUserRepo)

	token, err := service.Authenticate(ctx, username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func Test_Authenticate_GetByUsernameReturnedError_Failure(t *testing.T) {
	ctx := context.Background()
	username, password := "user", "pass"

	mockUserRepo := mocks.NewMockUserRepo(t)

	mockUserRepo.On("GetUserByUsername", ctx, username).Return(db.User{}, errors.New("error"))

	service := NewAuthService(cfg, mockUserRepo)

	token, err := service.Authenticate(ctx, username, password)

	assert.Error(t, err)
	assert.Empty(t, token)
}

func Test_Authenticate_CreateUserReturnedError_Failure(t *testing.T) {
	ctx := context.Background()
	username, password := "user", "pass"

	mockUserRepo := mocks.NewMockUserRepo(t)

	mockUserRepo.On("GetUserByUsername", ctx, username).Return(db.User{}, apperr.ErrNotFound)
	mockUserRepo.On("CreateUser", ctx, username, mock.MatchedBy(func(passwordHash string) bool {
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		return err == nil
	})).Return(int32(0), errors.New("error"))

	service := NewAuthService(cfg, mockUserRepo)

	token, err := service.Authenticate(ctx, username, password)

	assert.Error(t, err)
	assert.Empty(t, token)
}

func Test_Authenticate_WrongPassword_Failure(t *testing.T) {
	ctx := context.Background()
	username, password, wrongPassword := "user", "pass", "wrong_pass"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUserRepo := mocks.NewMockUserRepo(t)

	mockUserRepo.On("GetUserByUsername", ctx, username).Return(db.User{
		ID:           1,
		Username:     username,
		PasswordHash: string(passwordHash),
		Balance:      1000,
	}, nil)

	service := NewAuthService(cfg, mockUserRepo)

	token, err := service.Authenticate(ctx, username, wrongPassword)

	assert.Error(t, err)
	assert.Empty(t, token)
}
