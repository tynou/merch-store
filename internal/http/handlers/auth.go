package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tynou/avito-assignment/internal/http/common"
	"github.com/tynou/avito-assignment/internal/service/auth"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *ApiHandler) Auth(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Некорректное тело запроса.")
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Поля username и password обязательны.")
		return
	}

	token, err := h.auth.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			common.RespondWithError(w, http.StatusUnauthorized, "Вы не авторизованы.")
			return
		}

		common.RespondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера.")
		return
	}

	common.RespondWithJson(w, http.StatusOK, AuthResponse{
		Token: token,
	})
}
