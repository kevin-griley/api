package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/kevin-griley/api/internal/middleware"
)

func GetPathID(r *http.Request) (uuid.UUID, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("id is required")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id: %s", idStr)
	}

	return id, nil
}

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

func WriteJSON(w http.ResponseWriter, status int, v any) *ApiError {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return &ApiError{http.StatusInternalServerError, err.Error()}
	}
	return nil
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) *ApiError

func HandleApiError(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {

			reqID, ok := middleware.GetRequestID(r.Context())
			if !ok {
				reqID = "unknown"
			}

			slog.Error("API Error",
				"status", err.Status,
				"error", err.Message,
				"requestID", reqID,
			)
			WriteJSON(w, err.Status, err)
		}
	}
}

func DecodeJSONRequest(r *http.Request, dest any, maxSize int64) error {
	if maxSize <= 0 {
		maxSize = 1 << 20 // Default to 1MB if not specified
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return fmt.Errorf("invalid content type: expected application/json")
	}

	body := http.MaxBytesReader(nil, r.Body, maxSize)
	defer r.Body.Close()

	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		if err == io.EOF {
			return fmt.Errorf("empty request body")
		}
		return fmt.Errorf("invalid JSON: %v", err)
	}

	if decoder.More() {
		return fmt.Errorf("unexpected extra JSON content")
	}

	return nil
}
