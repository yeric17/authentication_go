package utils

type ErrorResponse struct {
	ErrorCode uint   `json:"error_code"`
	Message   string `json:"message"`
}

type DefaultResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}
