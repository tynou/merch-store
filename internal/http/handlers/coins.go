package handlers

import "net/http"

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

func (h *ApiHandler) SendCoin(w http.ResponseWriter, r *http.Request) {

}
