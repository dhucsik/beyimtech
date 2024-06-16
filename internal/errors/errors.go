package errors

import "net/http"

var (
	ErrNotSupportedImgFormat = NewAPIError("not supported img format", http.StatusBadRequest, "400001")
	ErrFileSizeExceeds       = NewAPIError("file size exceeds 5MB", http.StatusBadRequest, "400002")
)

type APIError struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	InnerCode string `json:"inner_code"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(message string, code int, innerCode string) *APIError {
	return &APIError{
		Message:   message,
		Code:      code,
		InnerCode: innerCode,
	}
}
