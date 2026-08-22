package handlers

import (
	"errors"
	"net/http"

	"github.com/tynou/avito-assignment/internal/http/common"
	"github.com/tynou/avito-assignment/internal/repo"
	"github.com/tynou/avito-assignment/internal/service/purchase"
)

func (h *ApiHandler) Buy(w http.ResponseWriter, r *http.Request) {
	userId, _ := common.GetUserIdFromContext(r.Context())

	item := r.PathValue("item")

	err := h.purchase.Buy(r.Context(), userId, item)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			common.RespondWithError(w, http.StatusBadRequest, "Мерч не найден.")
			return
		}

		if errors.Is(err, purchase.ErrInsufficientFunds) {
			common.RespondWithError(w, http.StatusBadRequest, "Недостаточно средств для покупки.")
			return
		}

		common.RespondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера.")
		return
	}

	common.RespondWithText(w, http.StatusOK, "Покупка прошла успешно.")
}
