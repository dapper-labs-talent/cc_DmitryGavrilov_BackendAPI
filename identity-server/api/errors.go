package api

import (
	"fmt"
	"net/http"
)

type HttpErrorResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"errorMessage"`
}

func (er *HttpErrorResponse) Error() string {
	return fmt.Sprintf("Code = %d; ErrorMessage = %s", er.Code, er.ErrorMessage)
}

func badRequestError(message string) *HttpErrorResponse {
	return &HttpErrorResponse{Code: http.StatusBadRequest, ErrorMessage: message}
}

func internalServerError(message string) *HttpErrorResponse {
	return &HttpErrorResponse{Code: http.StatusInternalServerError, ErrorMessage: message}
}
