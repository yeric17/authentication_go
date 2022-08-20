package utils

type ErrorResponse struct {
	Message string `json:"message"`
}

type DefaultResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
