package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kevin-griley/api/internal/types"
)

func GenerateRequestID() string {
	return uuid.New().String()
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate or retrieve a requestID in context.
		reqID, ok := r.Context().Value(types.ContextKeyRequestID).(string)
		if !ok || reqID == "" {
			reqID = GenerateRequestID()
			ctx := context.WithValue(r.Context(), types.ContextKeyRequestID, reqID)
			r = r.WithContext(ctx)
		}

		// Retrieve remote IP (r.RemoteAddr isn't perfect, check X-Forwarded-For if behind a proxy)
		remoteIP := r.RemoteAddr

		start := time.Now()
		slog.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
			"requestID", reqID,
			"remoteIP", remoteIP,
		)

		next(w, r)

		duration := time.Since(start)
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration,
			"requestID", reqID,
		)
	}
}
