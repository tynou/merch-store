package handlers

import "net/http"

type ItemDto struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceivedCoinsDto struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoinsDto struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistoryDto struct {
	Received []ReceivedCoinsDto `json:"received"`
	Sent     []SentCoinsDto     `json:"sent"`
}

type InfoResponse struct {
	Coins       int            `json:"coins"`
	Inventory   []ItemDto      `json:"inventory"`
	CoinHistory CoinHistoryDto `json:"coinHistory"`
}

func (h *ApiHandler) Info(w http.ResponseWriter, r *http.Request) {

}
