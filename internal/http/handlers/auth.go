package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

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
		respondWithError(w, http.StatusBadRequest, "Неверное тело запроса.")
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Поля username и password обязательны.")
		return
	}

	token, err := h.auth.Authenticate(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера.")
		return
	}

	respondWithJson(w, http.StatusOK, AuthResponse{
		Token: token,
	})
}
