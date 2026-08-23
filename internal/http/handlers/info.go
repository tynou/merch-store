package handlers

import (
	"net/http"

	"github.com/tynou/avito-assignment/internal/http/common"
)

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int32  `json:"quantity"`
}

type ReceivedTransfer struct {
	FromUser string `json:"fromUser"`
	Amount   int32  `json:"amount"`
}

type SentTransfer struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedTransfer `json:"received"`
	Sent     []SentTransfer     `json:"sent"`
}

type InfoResponse struct {
	Coins       int32           `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

func (h *ApiHandler) Info(w http.ResponseWriter, r *http.Request) {
	userId, _ := common.GetUserIdFromContext(r.Context())

	response, err := h.info.GetInfo(r.Context(), userId)
	if err != nil {
		common.RespondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера.")
		return
	}

	common.RespondWithJson(w, http.StatusOK, response)
}
