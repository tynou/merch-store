package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SendCoins_Success_E2E(t *testing.T) {
	t.Cleanup(func() {
		cleanDB(t)
	})

	user1, user2 := "user1", "user2"

	token1 := authenticate(t, user1, "password")
	token2 := authenticate(t, user2, "password")

	sendCoin(t, token1, user2, 100)

	info1 := getInfo(t, token1)
	info2 := getInfo(t, token2)

	assert.Equal(t, int32(900), info1.Coins)
	assert.Equal(t, int32(1100), info2.Coins)

	assert.Equal(t, int32(100), info1.CoinHistory.Sent[0].Amount)
	assert.Equal(t, "user2", info1.CoinHistory.Sent[0].ToUser)

	assert.Equal(t, int32(100), info2.CoinHistory.Received[0].Amount)
	assert.Equal(t, "user1", info2.CoinHistory.Received[0].FromUser)
}
