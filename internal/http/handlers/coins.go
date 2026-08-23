package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tynou/avito-assignment/internal/apperr"
	"github.com/tynou/avito-assignment/internal/http/common"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int32  `json:"amount" validate:"required"`
}

func (h *ApiHandler) SendCoin(w http.ResponseWriter, r *http.Request) {
	fromUserId, _ := common.GetUserIdFromContext(r.Context())

	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Некорректное тело запроса.")
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Поля toUser и amount обязательны.")
		return
	}

	err := h.transfer.SendCoin(r.Context(), fromUserId, req.ToUser, req.Amount)
	if err != nil {
		if errors.Is(err, apperr.ErrSelfTransfer) {
			common.RespondWithError(w, http.StatusBadRequest, "Нельзя перевести монеты самому себе.")
			return
		}

		if errors.Is(err, apperr.ErrNotFound) {
			common.RespondWithError(w, http.StatusBadRequest, "Пользователь не найден.")
			return
		}

		if errors.Is(err, apperr.ErrInsufficientFunds) {
			common.RespondWithError(w, http.StatusBadRequest, "Недостаточно средств для перевода.")
			return
		}

		common.RespondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера.")
		return
	}

	common.RespondWithText(w, http.StatusOK, "Передача монет успешно выполнена.")
}
