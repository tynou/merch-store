package info

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tynou/avito-assignment/internal/apperr"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/service/info/mocks"
)

func TestGetUserInfo_Success(t *testing.T) {
	ctx := context.Background()
	userID := int32(1)

	mockUserRepo := mocks.NewMockUserRepo(t)
	mockPurchaseRepo := mocks.NewMockPurchaseRepo(t)
	mockTransferRepo := mocks.NewMockTransferRepo(t)

	mockUserRepo.On("GetUserBalance", ctx, userID).Return(int32(1000), nil)

	mockPurchaseRepo.On("GetUserInventory", ctx, userID).Return([]db.GetUserInventoryRow{
		{Type: "t-shirt", Quantity: 2},
		{Type: "book", Quantity: 1},
	}, nil)

	mockTransferRepo.On("GetReceivedTransfers", ctx, userID).Return([]db.GetReceivedTransfersRow{
		{FromUser: "user1", Amount: 100},
	}, nil)

	mockTransferRepo.On("GetSentTransfers", ctx, userID).Return([]db.GetSentTransfersRow{
		{ToUser: "user2", Amount: 50},
		{ToUser: "user3", Amount: 67},
	}, nil)

	service := NewInfoService(mockUserRepo, mockPurchaseRepo, mockTransferRepo)

	resp, err := service.GetInfo(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1000), resp.Coins)
	assert.Len(t, resp.Inventory, 2)
	assert.Equal(t, "t-shirt", resp.Inventory[0].Type)
	assert.Equal(t, int32(2), resp.Inventory[0].Quantity)
	assert.Len(t, resp.CoinHistory.Received, 1)
	assert.Equal(t, "user1", resp.CoinHistory.Received[0].FromUser)
	assert.Equal(t, int32(100), resp.CoinHistory.Received[0].Amount)
	assert.Len(t, resp.CoinHistory.Sent, 2)
}

func TestGetUserInfo_GetUserBalance_Failure(t *testing.T) {
	ctx := context.Background()
	userID := int32(1)

	mockUserRepo := mocks.NewMockUserRepo(t)

	mockUserRepo.On("GetUserBalance", ctx, userID).Return(int32(0), apperr.ErrNotFound)

	service := NewInfoService(mockUserRepo, mocks.NewMockPurchaseRepo(t), mocks.NewMockTransferRepo(t))

	resp, err := service.GetInfo(ctx, userID)

	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestGetUserInfo_GetUserInventory_Failure(t *testing.T) {
	ctx := context.Background()
	userID := int32(1)

	mockUserRepo := mocks.NewMockUserRepo(t)
	mockPurchaseRepo := mocks.NewMockPurchaseRepo(t)

	mockUserRepo.On("GetUserBalance", ctx, userID).Return(int32(1000), nil)

	mockPurchaseRepo.On("GetUserInventory", ctx, userID).Return([]db.GetUserInventoryRow{}, errors.New("error"))

	service := NewInfoService(mockUserRepo, mockPurchaseRepo, mocks.NewMockTransferRepo(t))

	resp, err := service.GetInfo(ctx, userID)

	assert.Nil(t, resp)
	assert.Error(t, err)
}
