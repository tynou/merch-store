package handlers

import (
	"net/http"

	"github.com/tynou/avito-assignment/internal/domain"
	"github.com/tynou/avito-assignment/internal/http/common"
)

type InfoResponse struct {
	Coins       int32                  `json:"coins"`
	Inventory   []domain.InventoryItem `json:"inventory"`
	CoinHistory domain.CoinHistory     `json:"coinHistory"`
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
