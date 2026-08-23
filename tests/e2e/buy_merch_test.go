package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BuyMerch_Success_E2E(t *testing.T) {
	t.Cleanup(func() {
		cleanDB(t)
	})

	token := authenticate(t, "test_user", "password")

	buyMerch(t, token, "t-shirt")

	info := getInfo(t, token)

	assert.Equal(t, int32(920), info.Coins)
	assert.Len(t, info.Inventory, 1)
	assert.Equal(t, "t-shirt", info.Inventory[0].Type)
	assert.Equal(t, int32(1), info.Inventory[0].Quantity)
}
