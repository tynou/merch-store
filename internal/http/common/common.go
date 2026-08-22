package common

import (
	"context"
	"encoding/json"
	"net/http"
)

var JWTSecret = []byte("my_super_secret_key")

var UserIdKey = "user_id"

type ErrorResponse struct {
	Errors string `json:"errors"`
}

func GetUserIdFromContext(ctx context.Context) (int32, bool) {
	userID, ok := ctx.Value(UserIdKey).(int32)
	return userID, ok
}

func RespondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJson(w, status, ErrorResponse{
		Errors: message,
	})
}
