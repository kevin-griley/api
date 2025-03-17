package handlers

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func Error(status int, message string) *ApiError {
	return &ApiError{Status: status, Message: message}
}

func HTTPErrorHandler(err error, w http.ResponseWriter) {
	apiErr, ok := err.(*ApiError)
	if !ok {
		apiErr = Error(http.StatusInternalServerError, err.Error())
	}
	WriteJSON(w, apiErr.Status, apiErr)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) *ApiError {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}
	return nil
}
