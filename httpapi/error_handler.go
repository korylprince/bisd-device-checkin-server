package httpapi

import (
	"errors"
	"net/http"
)

//ErrorResponse represents an HTTP error
type ErrorResponse struct {
	Code        int    `json:"code"`
	Error       string `json:"error"`
	Description string `json:"description"`
}

//handleError returns a handlerResponse response for the given code
func handleError(code int, err error) *handlerResponse {
	return &handlerResponse{Code: code, Body: &ErrorResponse{Code: code, Error: http.StatusText(code), Description: err.Error()}, Err: err}
}

//notFoundHandler returns a 401 handlerResponse
func notFoundHandler(w http.ResponseWriter, r *http.Request) *handlerResponse {
	return handleError(http.StatusNotFound, errors.New("Could not find handler"))
}

//checkAPIError checks an api.Error and returns a handlerResponse for it, or nil if there was no error
func checkAPIError(err error) *handlerResponse {
	if err == nil {
		return nil
	}

	return handleError(http.StatusInternalServerError, err)
}
