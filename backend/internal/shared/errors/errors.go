package errors

import "net/http"

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func BadRequest(msg string) ErrorResponse {
	return ErrorResponse{Message: msg, Code: http.StatusBadRequest}
}

func Internal(msg string) ErrorResponse {
	return ErrorResponse{Message: msg, Code: http.StatusInternalServerError}
}

// etc...
