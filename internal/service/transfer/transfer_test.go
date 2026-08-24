package transfer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/service/transfer/mocks"
)

func Test_SendCoin_Success(t *testing.T) {
	ctx := context.Background()
	fromUserId, toUserId, toUserName, amount := int32(1), int32(2), "user2", int32(100)

	mockUserRepo := mocks.NewMockUserRepo(t)
	mockTransferRepo := mocks.NewMockTransferRepo(t)

	mockUserRepo.On("GetUserByUsername", ctx, toUserName).Return(db.User{
		ID: toUserId,
	}, nil)

	mockTransferRepo.On("CreateTransfer", ctx, fromUserId, toUserId, amount).Return(nil)

	service := NewTransferService(mockUserRepo, mockTransferRepo)

	err := service.SendCoin(ctx, fromUserId, toUserName, amount)

	assert.NoError(t, err)
}
