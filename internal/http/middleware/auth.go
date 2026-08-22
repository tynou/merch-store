package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tynou/avito-assignment/internal/http/common"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.RespondWithError(w, http.StatusUnauthorized, "Отсутствует заголовок 'Authorization'.")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.RespondWithError(w, http.StatusUnauthorized, "Некорректный формат заголовка авторизации.")
			return
		}

		jwtToken := parts[1]

		token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (any, error) {
			return common.JWTSecret, nil
		})
		if err != nil || !token.Valid {
			common.RespondWithError(w, http.StatusUnauthorized, "Некорректный токен.")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			common.RespondWithError(w, http.StatusUnauthorized, "Некорректный токен.")
			return
		}

		ctx := context.WithValue(r.Context(), common.UserIdKey, int32(claims[common.UserIdKey].(float64)))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
