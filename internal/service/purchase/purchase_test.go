package purchase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tynou/avito-assignment/internal/db"
	"github.com/tynou/avito-assignment/internal/service/purchase/mocks"
)

func Test_Buy_Success(t *testing.T) {
	ctx := context.Background()
	userId, item := int32(1), "t-shirt"
	merch := db.Merch{
		ID:    1,
		Name:  item,
		Price: 50,
	}

	mockMerchRepo := mocks.NewMockMerchRepo(t)
	mockPurchaseRepo := mocks.NewMockPurchaseRepo(t)

	mockMerchRepo.On("GetMerchByName", ctx, item).Return(merch, nil)

	mockPurchaseRepo.On("CreatePurchase", ctx, userId, merch).Return(nil)

	service := NewPurchaseService(mockMerchRepo, mockPurchaseRepo)

	err := service.Buy(ctx, userId, item)

	assert.NoError(t, err)
}
