package handlers

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type ApiHandler struct {
}

func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}
